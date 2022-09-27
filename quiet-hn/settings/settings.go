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
		Timeout:    5,
	}

	flag.IntVar(&Settings.Port, "port", 8080, "Port to start server on")
	flag.IntVar(&Settings.MaxStories, "n", 30, "Number of stories to fetch")
	flag.DurationVar(&Settings.Timeout, "timeout", 5*time.Second, "Timeout for cache stories")

	flag.Parse()
}
