package sitemap

import (
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/AksAman/gophercises/linkparser"
	"github.com/AksAman/gophercises/sitemap/models"
	"github.com/AksAman/gophercises/sitemap/utils"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	utils.InitializeLogger("sitemap.log")
	logger = utils.Logger
}

func Generate(websiteURL *url.URL, maxDepth int) (*models.Sitemap, error) {
	logger.Infof("generating sitemap for site: %v\n", websiteURL)

	finalURL, err := getFinalURL(websiteURL)
	if err != nil {
		return nil, err
	}

	baseHost := finalURL.Hostname()
	return generateSitemap(finalURL, baseHost, maxDepth)
}

// Gets the base url either from final redirect or from a long url
func getFinalURL(originalURL *url.URL) (*url.URL, error) {
	r, err := http.Head(originalURL.String())
	if err != nil {
		return nil, err
	}
	finalURL := r.Request.URL
	return &url.URL{
		Scheme: finalURL.Scheme,
		Host:   finalURL.Host,
	}, nil
}

func getHTMLBody(websiteURLStr string) ([]byte, error) {
	response, err := http.Get(websiteURLStr)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	return data, err
}

func getURLsFromResourceURL(baseResourceURL *url.URL, basehost string, depth int) ([]models.URL, error) {
	urls := []models.URL{}

	visitedURLs := map[string]struct{}{}
	crawledURLs := map[models.URL]struct{}{}
	currentQ := map[string]struct{}{
		baseResourceURL.String(): {},
	}
	nextQ := map[string]struct{}{}

	filterFunc := func(uStr string) bool {
		return !isInvalidURLString(uStr, basehost) && !isAlreadyVisitedURL(uStr, &visitedURLs)
	}

	for d := 0; d <= depth; d++ {
		// logger.Infof("At Level: %v, %d urls to be processed : %v", d, len(currentQ), currentQ)
		for currentURL := range currentQ {
			currentURL := cleanURL(currentURL)
			if _, seen := visitedURLs[currentURL]; seen {
				// fmt.Printf("\t\t ------- %q already seen\n", currentURL)
				continue
			}
			visitedURLs[currentURL] = struct{}{}
			crawledURLs[models.URL{Loc: currentURL, Depth: d}] = struct{}{}
			bodyBytes, err := getHTMLBody(currentURL)
			if err != nil {
				logger.Errorf("error while getting htmlbody from url: %v, : %q ", err, baseResourceURL)
				continue
			}
			links, err := linkparser.ParseBytes(bodyBytes)
			if err != nil {
				logger.Errorf("error while getting links from url: %v, : %q ", err, baseResourceURL)
				continue
			}
			// flinks := filterLinks(links, filterFunc)
			// logger.Infof("\tAt Level: %v, %d of %d links to be processed : %v", d, len(flinks), len(links), flinks)
			for _, link := range links {
				if isInvalidURLString(link.Href, basehost) || isAlreadyVisitedURL(link.Href, &visitedURLs) {
					continue
				}
				completePath := baseResourceURL.String() + cleanURL(link.Href)
				nextQ[completePath] = struct{}{}
			}
		}
		currentQ, nextQ = nextQ, map[string]struct{}{}
		currentQ = filterURLQ(currentQ, filterFunc)
	}

	for u := range crawledURLs {
		urls = append(urls, u)
	}
	return urls, nil
}

func cleanURL(u string) string {
	u = strings.TrimSuffix(u, "/")
	u = strings.TrimSuffix(u, "#")
	return u
}

func isAlreadyVisitedURL(uStr string, visitedURLs *map[string]struct{}) bool {
	_, seen := (*visitedURLs)[uStr]
	return seen
}

func generateSitemap(baseResourceURL *url.URL, basehost string, maxDepth int) (*models.Sitemap, error) {
	visitedURLs := map[string]struct{}{}

	urls, err := getURLsFromResourceURL(baseResourceURL, basehost, maxDepth)
	if err != nil {
		return nil, err
	}

	sitemapURLs := []models.URL{
		*models.NewURLWithPriority(baseResourceURL.String(), 1.0, 0),
	}

	visitedURLs[baseResourceURL.Path] = struct{}{}

	for _, u := range urls {
		parsedURL, err := url.Parse(u.Loc)
		if err != nil {
			continue
		}
		path := parsedURL.Path
		_, seen := visitedURLs[path]
		if seen {
			continue
		}
		if isInvalidURL(parsedURL, basehost) {
			logger.Infof("%q is invalid url with host %q", parsedURL.String(), basehost)
			continue
		}
		u.Priority = calculatePriority(u.Depth, maxDepth)
		sitemapURLs = append(sitemapURLs, u)
		visitedURLs[path] = struct{}{}
	}

	sort.Sort(models.URLsByDepthAndPriority(sitemapURLs))

	urlset := models.URLSet{
		URLs: sitemapURLs,
	}

	return models.NewSitemap(
		urlset,
	), nil
}

func calculatePriority(depth, maxDepth int) float64 {
	p := 1.0 - (float64(depth) / float64(maxDepth))
	return math.Round(p*100) / 100
}

func isInvalidURL(u *url.URL, basehost string) bool {
	// remove fragment
	path := u.Path
	currentHost := u.Hostname()

	return (currentHost != "" && currentHost != basehost) || path == "/" || strings.HasPrefix(u.String(), "mailto:") || strings.HasPrefix(u.String(), "tel:") || strings.HasPrefix(u.String(), "javascript:") || strings.HasPrefix(u.String(), "#")
}

func isInvalidURLString(uStr string, basehost string) bool {
	if u, err := url.Parse(uStr); err != nil {
		// logger.Errorf("error while parsing url: %v, : %q ", err, uStr)
		return true
	} else {
		return isInvalidURL(u, basehost)
	}
}

func filterURLQ(urlQ map[string]struct{}, filter func(string) bool) map[string]struct{} {
	filteredURLs := make(map[string]struct{})
	for u := range urlQ {
		if filter(u) {
			filteredURLs[u] = struct{}{}
		}
	}
	return filteredURLs
}

func filterLinks(links []linkparser.Link, filter func(string) bool) []linkparser.Link {
	filteredLinks := []linkparser.Link{}
	for _, link := range links {
		if filter(link.Href) {
			filteredLinks = append(filteredLinks, link)
		}
	}
	return filteredLinks
}
