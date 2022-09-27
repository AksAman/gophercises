package views

import (
	"fmt"
	"html/template"
	"time"

	"github.com/AksAman/gophercises/quietHN/caching"
	"github.com/AksAman/gophercises/quietHN/models"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/utils"
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

		cachedStories := cache.Get()
		if cachedStories == nil || len(cachedStories) < settings.Settings.MaxStories {
			existingStoriesCount := len(cachedStories)
			existingStoriesCount = utils.Clamp(settings.Settings.MaxStories-existingStoriesCount, 0, settings.Settings.MaxStories)

			stories, err := getStories(existingStoriesCount, getStoriesForIDsAsync, &cache)
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
