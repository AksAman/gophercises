package settings

import (
	"flag"
	"time"
)

type GlobalSettings struct {
	Port       int
	MaxStories int
	Timeout    time.Duration
}

var Settings *GlobalSettings

func init() {
	Settings = &GlobalSettings{
		Port:       8080,
		MaxStories: 30,
		Timeout:    time.Second * 10,
	}

	flag.IntVar(&Settings.Port, "port", Settings.Port, "Port to start server on")
	flag.IntVar(&Settings.MaxStories, "n", Settings.MaxStories, "Number of stories to fetch")
	flag.DurationVar(&Settings.Timeout, "timeout", Settings.Timeout, "Timeout for cache stories")

	flag.Parse()
}
