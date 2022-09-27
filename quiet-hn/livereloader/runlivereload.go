package main

import (
	"flag"
	"log"

	"github.com/AksAman/gophercises/quietHN/livereloader/livereload"
)

var port int

func init() {
	flag.IntVar(&port, "port", 9090, "Port to start livereload server on")
	flag.Parse()
}

func main() {
	log.Fatal(livereload.StartLiveReloadSocketOnPort(port))
}
