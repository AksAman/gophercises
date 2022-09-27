package views

import (
	"fmt"
	"sync"
	"time"

	"github.com/AksAman/gophercises/quietHN/models"
)

var cacheMutex sync.Mutex

type storyCache struct {
	stories []*models.HNItem

	timeout time.Duration
	expiry  time.Time
}

func (c *storyCache) IsExpired() bool {
	return time.Now().After(c.expiry)
}

func (c *storyCache) Set(stories []*models.HNItem) {
	fmt.Println("LOCKING")
	cacheMutex.Lock()
	defer func() {
		fmt.Println("UN-LOCKING")
		cacheMutex.Unlock()

	}()
	fmt.Printf("Setting %d stories in cache with timeout: %v\n", len(stories), c.timeout)
	c.stories = stories

	c.expiry = time.Now().Add(c.timeout)
}

func (c *storyCache) Get() []*models.HNItem {
	// cacheMutex.Lock()
	// defer cacheMutex.Unlock()
	if c.IsExpired() {
		fmt.Println("Cache expired, fetching new stories")
		return nil
	}
	fmt.Printf("Returning %d stories from cache", len(c.stories))
	// for _, story := range c.stories {
	// 	story.Latency = 0
	// }
	return c.stories
}
