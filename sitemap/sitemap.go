package sitemap

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	"github.com/AksAman/gophercises/linkparser"
	"github.com/AksAman/gophercises/sitemap/models"
	"github.com/AksAman/gophercises/sitemap/utils"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	utils.InitializeLogger()
	logger = utils.Logger
}

func Generate(websiteURL *url.URL, maxDepth int) (*models.Sitemap, error) {
	logger.Infof("generating sitemap for site: %v\n", websiteURL)

	finalURL, err := getFinalURL(websiteURL)
	if err != nil {
		return nil, err
	}

	body, err := getHTMLBody(finalURL.String())
	if err != nil {
		return nil, err
	}
	return generateSitemap(body, finalURL)
}

func getHTMLBody(websiteURLStr string) ([]byte, error) {
	response, err := http.Get(websiteURLStr)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
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

func generateSitemap(siteBytes []byte, websiteURL *url.URL) (*models.Sitemap, error) {

	visitedURLs := map[string]bool{}
	links, _ := linkparser.Parse(bytes.NewReader(siteBytes))

	urlset := models.URLSet{
		URLs: []models.URL{
			{
				Loc:      websiteURL.String(),
				Priority: 1.00,
			},
		},
	}
	visitedURLs[websiteURL.Path] = true

	for _, link := range links {
		u, err := url.Parse(link.Href)
		if err != nil {
			// fmt.Printf("hostname: %q, err: %v\n", hostname, err)
			continue
		}
		path := u.Path
		if u.Hostname() != "" && u.Hostname() != websiteURL.Hostname() {
			continue
		}
		if u.Path == "" || u.Path == "/" {
			continue
		}
		if visitedURLs[path] {
			continue
		}
		// hrefs := strings.Split(link.Href, hostname)
		urlset.URLs = append(urlset.URLs, models.URL{
			Loc:      websiteURL.String() + path,
			Priority: 0.50,
		})
		visitedURLs[path] = true
	}

	return models.NewSitemap(
		urlset,
	), nil
}
