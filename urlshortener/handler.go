package urlshortener

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

//go:embed data
var dataFS embed.FS

func handleRedirect(w http.ResponseWriter, r *http.Request, src, dest string) {
	fmt.Println("handle redirect")
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `<h3>Redirecting from <a href="%v">%v</a> to <a href="%v">%v</a></h3><hr>`, src, src, dest, dest)
	// http.Redirect(w, r, dest, http.StatusFound)
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		src := r.URL.Path
		if dest, ok := pathsToUrls[src]; ok {
			handleRedirect(w, r, src, dest)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.

type pathToUrl struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

func buildMapFromPathToUrls(pathToUrls []pathToUrl) map[string]string {
	pathToUrlsMap := map[string]string{}
	for _, pathToUrl := range pathToUrls {
		pathToUrlsMap[pathToUrl.Path] = pathToUrl.Url
	}
	return pathToUrlsMap
}

func YAMLHandler(yml []byte, fallback http.Handler) http.HandlerFunc {
	pathToUrls := []pathToUrl{}

	err := yaml.Unmarshal(yml, &pathToUrls)
	if err != nil {
		return fallback.ServeHTTP
	}

	pathToUrlsMap := buildMapFromPathToUrls(pathToUrls)
	return MapHandler(pathToUrlsMap, fallback)
}

func YAMLFileHandler(filename string, fallback http.Handler) http.HandlerFunc {

	// read json file
	ymlBytes, err := dataFS.ReadFile(filename)
	if err != nil {
		log.Println("error while reading file", err)
		return fallback.ServeHTTP
	}

	return YAMLHandler(ymlBytes, fallback)
}

func JSONFileHandler(filename string, fallback http.Handler) http.HandlerFunc {

	// read json file
	jsonBytes, err := dataFS.ReadFile(filename)
	if err != nil {
		log.Println("error while reading file", err)
		return fallback.ServeHTTP
	}

	// convert raw Bytes to struct data
	pathToUrls := []pathToUrl{}
	err = json.Unmarshal(jsonBytes, &pathToUrls)
	if err != nil {
		log.Println("error while converting to json", err)
		return fallback.ServeHTTP
	}

	pathToUrlsMap := buildMapFromPathToUrls(pathToUrls)
	return MapHandler(pathToUrlsMap, fallback)
}
