package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

var logger *zap.SugaredLogger

func main() {
	utils.InitializeLogger("server.log")
	logger = utils.Logger
	RunServer()
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Infof("%v %v", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func RunServer() {
	router := http.NewServeMux()
	routerWithLogger := LoggingMiddleware(router)

	router.HandleFunc("/sitemap/", CreateSitemapView)

	address := ":8080"
	logger.Info("Starting server at", address)
	server := http.Server{
		Addr:    address,
		Handler: routerWithLogger,
	}

	server.ListenAndServe()
}

func CreateSitemapView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	websiteURL := r.URL.Query().Get("url")
	maxDepth := 3
	if r.URL.Query().Has("depth") {
		maxDepth, _ = strconv.Atoi(r.URL.Query().Get("depth"))
	}
	if websiteURL == "" {
		logger.Error("No url passed")
		JSONError(w, HTTPError{Error: "No url passed"}, http.StatusBadRequest)
		return
	}
	logger.Info("Creating sitemap for url: ", websiteURL, " with max depth: ", maxDepth)

	parsedURL, err := url.Parse(string(websiteURL))
	if err != nil {
		logger.Error("Unable to parse website url", zap.String("url", websiteURL), zap.Error(err))
		JSONError(w, HTTPError{Error: "Invalid url passed"}, http.StatusBadRequest)
		return
	}

	sitemap, err := sitemap.Generate(parsedURL, maxDepth)
	if err != nil {
		logger.Error("Unable to generate sitemap for site", zap.String("url", websiteURL), zap.Error(err))
		JSONError(w, HTTPError{Error: err.Error()}, http.StatusBadRequest)
		return
	}

	serializers := map[string]struct {
		serializer  func() ([]byte, error)
		contentType string
	}{
		"json": {
			serializer:  sitemap.SerializeToJSON,
			contentType: "application/json",
		},
		"xml": {
			serializer:  sitemap.SerializeToXML,
			contentType: "application/xml",
		},
	}
	s, found := serializers[format]

	if !found {
		logger.Error("Invalid format passed", zap.String("format", format))
		JSONError(w, HTTPError{Error: "Invalid format passed"}, http.StatusBadRequest)
		return
	}

	serialized, err := s.serializer()
	w.Header().Set("Content-Type", s.contentType)
	if err != nil {
		logger.Error("Unable to serialize sitemap", zap.String("url", websiteURL), zap.Error(err))
		JSONError(w, HTTPError{Error: err.Error()}, http.StatusInternalServerError)
		return
	}
	w.Write(serialized)

}
