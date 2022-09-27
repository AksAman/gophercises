package views

import (
	"fmt"
	"html/template"
	"time"

	"github.com/AksAman/gophercises/quietHN/caching"
	"github.com/AksAman/gophercises/quietHN/models"
	"github.com/AksAman/gophercises/quietHN/settings"
)

var (
	storyTemplate *template.Template
	cache         caching.InMemoryCache[models.HNItem]
)

func init() {
	initTemplates()
	initCache()
}

func initCache() {
	cache = caching.InMemoryCache[models.HNItem]{Timeout: time.Duration(settings.Settings.Timeout)}
	cache.Init()

	cache.SetupTicker(func() {
		fmt.Println("Refreshing cache")

		if s := cache.Get(); s == nil {
			return
		} else {
			stories, err := getStories(len(s), getStoriesForIDsAsync)
			if err != nil {
				return
			}
			cache.Set(stories)
			fmt.Println("Cache refreshed")
		}

	})
}

func initTemplates() {
	storyTemplate = template.Must(template.ParseFiles("templates/index.gohtml"))
}
