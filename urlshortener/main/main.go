package main

import (
	"fmt"
	"net/http"

	"github.com/AksAman/gophercises/urlshortener"
)

func createDefaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	return mux
}

func main() {
	mux := createDefaultMux()

	mapHandler := getMapHandler(mux)
	yamlHandler := getYAMLHandler(mapHandler)
	yamlFileHandler := getYAMLFileHandler(yamlHandler)
	jsonFileHandler := getJSONFileHandler(yamlFileHandler)

	address := ":8080"
	fmt.Println("starting server on", address)
	err := http.ListenAndServe(address, jsonFileHandler)
	if err != nil {
		panic(err)
	}
}

func getJSONFileHandler(fallbackHandler http.HandlerFunc) http.HandlerFunc {
	return urlshortener.JSONFileHandler("./data/urls.json", fallbackHandler)
}

func getYAMLFileHandler(fallbackHandler http.HandlerFunc) http.HandlerFunc {
	return urlshortener.YAMLFileHandler("./data/urls.yaml", fallbackHandler)
}

func getYAMLHandler(fallbackHandler http.HandlerFunc) http.HandlerFunc {
	yamlPathsToUrls := `
- path: /my-urlshortener
  url:  https://github.com/AksAman/gophercises.git
- path: /john-urlshortener
  url:  https://github.com/gophercises/urlshort
`
	yamlHandler := urlshortener.YAMLHandler([]byte(yamlPathsToUrls), fallbackHandler)
	return yamlHandler
}

func getMapHandler(fallbackHandler *http.ServeMux) http.HandlerFunc {
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
		"/dybo":           "https://dybo.io",
		"/mygit":          "https://github.com/AksAman",
	}

	mapHandler := urlshortener.MapHandler(pathsToUrls, fallbackHandler)
	return mapHandler
}
