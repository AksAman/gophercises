package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/AksAman/gophercises/quietHN/routing"
	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/gorilla/mux"
)

var port int

func init() {
	flag.IntVar(&port, "port", 8080, "Port to start server on")
	flag.Parse()
}

func RunServer() {
	router := mux.NewRouter()
	router.NotFoundHandler = views.NotFoundHandler{}

	for _, route := range routing.Routes {
		router.HandleFunc(route.Pattern, route.Handler)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	fmt.Printf("Starting server on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func main() {
	RunServer()
}
