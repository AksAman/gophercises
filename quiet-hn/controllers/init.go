package controllers

import (
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/AksAman/gophercises/quietHN/caching"
	"github.com/AksAman/gophercises/quietHN/models"
	"github.com/AksAman/gophercises/quietHN/ratelimiter"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/utils"
)

var (
	storyTemplate *template.Template
	cache         caching.Cache[models.HNItem]
	rateLimiter   ratelimiter.IRateLimiter
	counter       uint64
)

func init() {
	initTemplates()

	getProperCache()
	initCache(cache)

	initRateLimiter()

}

func initRateLimiter() {
	if settings.Settings.RateLimitingType == settings.NormalRateLimting {
		rateLimiter, _ = ratelimiter.NewRateLimiter(time.Duration(settings.Settings.RateLimitingInterval))
	} else if settings.Settings.RateLimitingType == settings.BurstyRateLimiting {
		rateLimiter, _ = ratelimiter.NewBurstyRateLimiter(time.Duration(settings.Settings.RateLimitingInterval), settings.Settings.BurstRateCount)
	} else {
		rateLimiter = nil
	}
}

func getProperCache() {
	if settings.Settings.CachingStrategy == settings.MemCacheStrategy {
		cache = &caching.InMemoryCache[models.HNItem]{Timeout: time.Duration(settings.Settings.Timeout)}
	} else if settings.Settings.CachingStrategy == settings.RedisCacheStrategy {
		cache = &caching.RedisCache[models.HNItem]{Timeout: time.Duration(settings.Settings.Timeout), ItemsKey: "stories"}
	} else {
		cache = &caching.NoCache[models.HNItem]{}
	}
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
