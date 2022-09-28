package views

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/AksAman/gophercises/quietHN/caching"
	"github.com/AksAman/gophercises/quietHN/models"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/utils"
)

var (
	storyTemplate *template.Template
	cache         caching.Cache[models.HNItem]
)

func init() {
	initTemplates()

	if settings.Settings.CachingStrategy == settings.MemCacheStrategy {
		cache = &caching.InMemoryCache[models.HNItem]{Timeout: time.Duration(settings.Settings.Timeout)}
	} else if settings.Settings.CachingStrategy == settings.RedisCacheStrategy {
		cache = &caching.RedisCache[models.HNItem]{Timeout: time.Duration(settings.Settings.Timeout), ItemsKey: "stories"}
	} else {
		cache = &caching.NoCache[models.HNItem]{}
	}

	initCache(cache)
}

func initCache(cache caching.Cache[models.HNItem]) {
	// cache = caching.InMemoryCache[models.HNItem]{Timeout: time.Duration(settings.Settings.Timeout)}
	err := cache.Init()
	if err != nil {
		log.Fatal(err)
	}

	cache.SetupTicker(func() {
		fmt.Println("Refreshing cache", cache.ToString())

		cachedStories := cache.Get()
		if cachedStories == nil || len(cachedStories) < settings.Settings.MaxStories {
			existingStoriesCount := len(cachedStories)
			existingStoriesCount = utils.Clamp(settings.Settings.MaxStories-existingStoriesCount, 0, settings.Settings.MaxStories)

			stories, err := getStories(existingStoriesCount, getStoriesForIDsAsync, cache)
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
