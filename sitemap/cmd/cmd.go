package main

import (
	"flag"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AksAman/gophercises/sitemap"
	"github.com/AksAman/gophercises/sitemap/utils"
	"go.uber.org/zap"
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
var maxDepth int
var useThreads bool

func init() {
	flag.Var(&websiteURLs, "url", "website addresses to be parsed to generate sitemap")
	flag.BoolVar(&useThreads, "threaded", false, "use threading for generation")
	flag.IntVar(&maxDepth, "depth", 3, "maximum depth of links to be parsed")
	flag.Parse()

	// if len(websiteURLs) == 0 {
	// 	logger.Errorf("No urls passed")
	// }
}

var logger *zap.SugaredLogger

func init() {
	utils.InitializeLogger("cmd.log")
	logger = utils.Logger
}

func main() {
	if !useThreads || len(websiteURLs) == 1 {
		RunNonThreaded()
	} else {
		RunThreaded()
	}
}

func RunNonThreaded() {
	start := time.Now()
	for _, websiteURL := range websiteURLs {
		CMDCreateSitemap(websiteURL)
	}
	logger.Infoln("RunNonThreaded Took ", time.Since(start))
}

func RunThreaded() {
	start := time.Now()
	var wg sync.WaitGroup
	for _, websiteURL := range websiteURLs {
		wg.Add(1)
		websiteURL := websiteURL
		go func() {
			defer wg.Done()
			CMDCreateSitemap(websiteURL)
		}()
	}
	wg.Wait()
	logger.Infoln("RunThreaded Took ", time.Since(start))
}

func CMDCreateSitemap(websiteURL WebsiteURL) {
	parsedURL, err := url.Parse(string(websiteURL))
	if err != nil {
		logger.Errorf("Unable to parse website url:%q due to error: \n\t%v\n", websiteURL, err)
		return
	}
	sitemap, err := sitemap.Generate(parsedURL, maxDepth)
	if err != nil {
		logger.Errorf("Unable to generate sitemap for site:%q due to error: \n\t%v\n", websiteURL, err)
		return
	}

	jsonFileName := parsedURL.Host + ".ignore.json"
	jsonFilePath := filepath.Join(".", "results", jsonFileName)
	err = sitemap.SerializeToFile(jsonFilePath)
	if err != nil {
		logger.Errorf("err: %v\n", err)
		return
	} else {
		logger.Infoln("Saved sitemap to", jsonFilePath)
	}
}
