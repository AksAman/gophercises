package sitemap

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/AksAman/gophercises/linkparser"
	"github.com/AksAman/gophercises/sitemap/models"
)

func Generate(websiteURL *url.URL) (*models.Sitemap, error) {
	fmt.Printf("generating sitemap for site: %v\n", websiteURL)

	body, err := getHTMLBody(websiteURL.String())
	if err != nil {
		return nil, err
	}
	return generateSitemap(body, websiteURL.Hostname())
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

func generateSitemap(siteBytes []byte, hostname string) (*models.Sitemap, error) {

	visitedURLs := map[string]bool{}
	links, _ := linkparser.Parse(bytes.NewReader(siteBytes))

	urlset := models.URLSet{}
	for _, link := range links {
		u, err := url.Parse(link.Href)
		if err != nil {
			// fmt.Printf("hostname: %q, err: %v\n", hostname, err)
			continue
		}
		path := u.Path
		if u.Hostname() != "" && u.Hostname() != hostname {
			continue
		}
		if u.Path == "" {
			continue
		}
		if visitedURLs[path] {
			continue
		}
		// hrefs := strings.Split(link.Href, hostname)
		urlset.URLs = append(urlset.URLs, models.URL{
			Loc:      path,
			Priority: 0,
		})
		visitedURLs[path] = true
	}

	return &models.Sitemap{
		Urlset: urlset,
	}, nil
}
