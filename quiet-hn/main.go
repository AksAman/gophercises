package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AksAman/gophercises/quietHN/routing"
	"github.com/AksAman/gophercises/quietHN/settings"
)

func RunServer() {
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Settings.Port),
		Handler: routing.NewRouter(),
	}

	fmt.Printf("Starting server on %s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func main() {
	RunServer()
}
