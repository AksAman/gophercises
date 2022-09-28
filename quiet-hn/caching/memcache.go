package caching

import (
	"fmt"
	"sync"
	"time"
)

type InMemoryCache[T CacheItem] struct {
	Timeout time.Duration

	items  []*T
	expiry time.Time
	mutex  *sync.Mutex

	ticker *time.Ticker
}

func (c *InMemoryCache[T]) ToString() string {
	return "InMemoryCache"
}

func (c *InMemoryCache[T]) Init() error {
	c.mutex = &sync.Mutex{}

	fmt.Println("InMemoryCache initialized")
	return nil
}

func (c *InMemoryCache[T]) SetupTicker(event func()) {
	c.ticker = time.NewTicker(c.Timeout / 2)

	go func() {
		for {
			event()
			<-c.ticker.C
		}
	}()
}

func (c *InMemoryCache[T]) IsExpired() bool {
	return time.Now().After(c.expiry)
}

func (c *InMemoryCache[T]) Set(items []*T) error {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()

	}()
	fmt.Printf("Setting %d items in cache with timeout: %v\n", len(items), c.Timeout)
	c.items = items

	c.expiry = time.Now().Add(c.Timeout)
	return nil
}

func (c *InMemoryCache[T]) Get() []*T {
	if c.IsExpired() {
		fmt.Println("Cache expired, fetching new items")
		return nil
	}
	fmt.Printf("Returning %d items from cache\n", len(c.items))
	return c.items
}
