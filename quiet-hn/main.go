package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AksAman/gophercises/quietHN/routing"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/gorilla/mux"
)

func RunServer() {
	router := mux.NewRouter()
	router.NotFoundHandler = views.NotFoundHandler{}

	for _, route := range routing.Routes {
		router.HandleFunc(route.Pattern, route.Handler)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Settings.Port),
		Handler: router,
	}

	fmt.Printf("Starting server on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func main() {
	RunServer()
}
