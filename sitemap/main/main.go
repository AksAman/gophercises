package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/AksAman/gophercises/sitemap"
	"github.com/AksAman/gophercises/sitemap/utils"
	"go.uber.org/zap"
)

type WebsiteURL string
type WebsiteURLs []WebsiteURL

type HTTPError struct {
	Error string `json:"error"`
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}

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
	// 	fmt.Printf("No urls passed")
	// }
}

var logger *zap.SugaredLogger

func main() {
	utils.InitializeLogger()
	logger = utils.Logger
	// if !useThreads || len(websiteURLs) == 1 {
	// 	RunNonThreaded()
	// } else {
	// 	RunThreaded()
	// }
	RunServer()
}

func RunNonThreaded() {
	start := time.Now()
	for _, websiteURL := range websiteURLs {
		CMDCreateSitemap(websiteURL)
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
			CMDCreateSitemap(websiteURL)
		}()
	}
	wg.Wait()
	fmt.Println("RunThreaded Took ", time.Since(start))
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%v %v", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func RunServer() {
	router := http.NewServeMux()
	configuredRouter := LoggingMiddleware(router)

	router.HandleFunc("/", CreateSitemapView)

	address := ":8080"
	logger.Info("Starting server at", address)
	server := http.Server{
		Addr:    address,
		Handler: configuredRouter,
	}

	server.ListenAndServe()
}

func CreateSitemapView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	websiteURL := r.URL.Query().Get("url")
	maxDepth := 3
	if r.URL.Query().Has("depth") {
		maxDepth, _ = strconv.Atoi(r.URL.Query().Get("depth"))
	}
	logger.Info("Creating sitemap for url: ", websiteURL, " with max depth: ", maxDepth)
	if websiteURL == "" {
		logger.Warn("No url passed")
		JSONError(w, HTTPError{Error: "No url passed"}, http.StatusBadRequest)
		return
	}

	parsedURL, err := url.Parse(string(websiteURL))
	if err != nil {
		logger.Warn("Unable to parse website url", zap.String("url", websiteURL), zap.Error(err))
		JSONError(w, HTTPError{Error: "Invalid url passed"}, http.StatusBadRequest)
		return
	}

	sitemap, err := sitemap.Generate(parsedURL, maxDepth)
	if err != nil {
		logger.Warn("Unable to generate sitemap for site", zap.String("url", websiteURL), zap.Error(err))
		JSONError(w, HTTPError{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	serialized, err := sitemap.Serialize()
	if err != nil {
		logger.Warn("Unable to serialize sitemap", zap.String("url", websiteURL), zap.Error(err))
		JSONError(w, HTTPError{Error: err.Error()}, http.StatusBadRequest)
		return
	}
	w.Write(serialized)

}

func CMDCreateSitemap(websiteURL WebsiteURL) {
	parsedURL, err := url.Parse(string(websiteURL))
	if err != nil {
		fmt.Printf("Unable to parse website url:%q due to error: \n\t%v\n", websiteURL, err)
		return
	}
	sitemap, err := sitemap.Generate(parsedURL, maxDepth)
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
