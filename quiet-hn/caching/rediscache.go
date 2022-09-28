package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/go-redis/redis/v9"
	"github.com/google/martian/v3/log"
)

const KEY_PREFIX = "GOPHERCISES:quietHN"

var ctx = context.Background()

type RedisCache[T CacheItem] struct {
	client   *redis.Client
	ItemsKey string

	Timeout time.Duration

	expiry time.Time
	mutex  *sync.Mutex

	ticker           *time.Ticker
	completeItemsKey string
}

func (c *RedisCache[T]) ToString() string {
	return "RedisCache"
}

func (c *RedisCache[T]) Init() error {
	c.mutex = &sync.Mutex{}

	c.completeItemsKey = fmt.Sprintf("%s:%s", KEY_PREFIX, c.ItemsKey)

	c.client = redis.NewClient(&redis.Options{
		Addr:     settings.Settings.RedisAddr,
		Password: settings.Settings.RedisPassword,
		DB:       0,
	})

	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		return err
	}

	fmt.Println("RedisCache initialized and connected")

	c.client.Set(ctx, "test", "test", 0)
	val, err := c.client.Get(ctx, "test").Result()
	if err != nil {
		return err
	}
	fmt.Printf("redis val: %v\n", val)
	return nil
}

func (c *RedisCache[T]) SetupTicker(event func()) {
	c.ticker = time.NewTicker(c.Timeout / 2)

	go func() {
		for {
			event()
			<-c.ticker.C
		}
	}()
}

func (c *RedisCache[T]) IsExpired() bool {
	return time.Now().After(c.expiry)
}

func (c *RedisCache[T]) Set(items []*T) error {
	c.mutex.Lock()
	defer func() {
		c.mutex.Unlock()

	}()

	fmt.Printf("Setting %d items in cache with timeout: %v\n", len(items), c.Timeout)
	serialized, err := json.Marshal(items)
	if err != nil {
		return err
	}

	c.expiry = time.Now().Add(c.Timeout)
	return c.client.Set(ctx, c.completeItemsKey, serialized, c.Timeout).Err()
}

func (c *RedisCache[T]) Get() []*T {
	if c.IsExpired() {
		// fmt.Println("Cache expired, fetching new items")
		return nil
	}
	val, err := c.client.Get(ctx, c.completeItemsKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Errorf("Error getting items from redis: %v", err)
		}
		return nil
	}

	var items []*T
	err = json.Unmarshal([]byte(val), &items)
	if err != nil {
		return nil
	}
	return items
}
