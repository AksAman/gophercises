package main

import (
	"flag"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AksAman/gophercises/sitemap"
)

type WebsiteURL string
type WebsiteURLs []WebsiteURL

func (urls *WebsiteURLs) String() string {
	var urlsAsString []string
	for _, url := range *urls {
		urlsAsString = append(urlsAsString, string(url))
	}
	return strings.Join(urlsAsString, ",")
}

func (urls *WebsiteURLs) Set(value string) error {
	*urls = append(*urls, WebsiteURL(value))
	return nil
}

var websiteURLs WebsiteURLs
var useThreads bool

func init() {
	flag.Var(&websiteURLs, "url", "website addresses to be parsed to generate sitemap")
	flag.BoolVar(&useThreads, "threaded", false, "use threading for generation")
	flag.Parse()

	if len(websiteURLs) == 0 {
		panic("No urls passed")
	}
}

func main() {
	if useThreads {
		RunThreaded()
	} else {
		RunNonThreaded()
	}
}

func RunNonThreaded() {
	start := time.Now()
	for _, websiteURL := range websiteURLs {
		CreateSitemap(websiteURL)
	}
	fmt.Println("RunNonThreaded Took ", time.Since(start))
}

func RunThreaded() {
	start := time.Now()
	var wg sync.WaitGroup
	for _, websiteURL := range websiteURLs {
		wg.Add(1)
		websiteURL := websiteURL
		go func() {
			defer wg.Done()
			CreateSitemap(websiteURL)
		}()
	}
	wg.Wait()
	fmt.Println("RunThreaded Took ", time.Since(start))
}

func CreateSitemap(websiteURL WebsiteURL) {
	parsedURL, err := url.Parse(string(websiteURL))
	if err != nil {
		fmt.Printf("Unable to parse website url:%q due to error: \n\t%v\n", websiteURL, err)
		return
	}
	sitemap, err := sitemap.Generate(parsedURL)
	if err != nil {
		fmt.Printf("Unable to generate sitemap for site:%q due to error: \n\t%v\n", websiteURL, err)
		return
	}

	jsonFileName := parsedURL.Host + ".ignore.json"
	jsonFilePath := filepath.Join(".", "results", jsonFileName)
	err = sitemap.SerializeToFile(jsonFilePath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	} else {
		fmt.Println("Saved sitemap to", jsonFilePath)
	}
}
