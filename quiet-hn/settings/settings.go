package settings

import (
	"flag"
	"fmt"
	"time"
)

type GlobalSettings struct {
	Port                 int
	MaxStories           int
	Timeout              time.Duration
	RedisHost            string
	RedisPort            int
	RedisAddr            string
	RedisPassword        string
	CachingStrategy      string
	RateLimitingType     string
	RateLimitingInterval time.Duration
	BurstRateCount       int
}

var Settings *GlobalSettings

const (
	MemCacheStrategy   = "mem"
	RedisCacheStrategy = "redis"
	NoCaching          = "none"
)

const (
	NormalRateLimting  = "normal"
	BurstyRateLimiting = "burst"
	NoRateLimiting     = "none"
)

func init() {
	Settings = &GlobalSettings{
		Port:                 8080,
		MaxStories:           30,
		Timeout:              time.Second * 10,
		RedisHost:            "localhost",
		RedisPort:            6379,
		RedisPassword:        "",
		CachingStrategy:      MemCacheStrategy,
		RateLimitingType:     NormalRateLimting,
		RateLimitingInterval: time.Second * 5,
		BurstRateCount:       5,
	}

	flag.IntVar(&Settings.Port, "port", Settings.Port, "Port to start server on")
	flag.IntVar(&Settings.MaxStories, "n", Settings.MaxStories, "Number of stories to fetch")
	flag.DurationVar(&Settings.Timeout, "timeout", Settings.Timeout, "Timeout for cache stories")
	flag.StringVar(&Settings.RedisHost, "redis-host", Settings.RedisHost, "Redis host")
	flag.IntVar(&Settings.RedisPort, "redis-port", Settings.RedisPort, "Redis port")
	flag.StringVar(&Settings.RedisPassword, "redis-password", Settings.RedisPassword, "Redis password")
	flag.Func("caching", "Caching strategy to use", func(s string) error {
		choices := []string{MemCacheStrategy, RedisCacheStrategy, NoCaching}
		for _, choice := range choices {
			if s == choice {
				Settings.CachingStrategy = s
				return nil
			}
		}
		return fmt.Errorf("\ninvalid caching strategy: %s, valid choices are: %v", s, choices)
	})

	flag.Func("rate-type", "Rate limiting strategy to use", func(s string) error {
		choices := []string{NormalRateLimting, BurstyRateLimiting, NoRateLimiting}
		for _, choice := range choices {
			if s == choice {
				Settings.RateLimitingType = s
				return nil
			}
		}
		return fmt.Errorf("\ninvalid rate limiting strategy: %s, valid choices are: %v", s, choices)
	})

	flag.DurationVar(&Settings.RateLimitingInterval, "rate-interval", Settings.RateLimitingInterval, "Rate limiting interval")
	flag.IntVar(&Settings.BurstRateCount, "rate-burst", Settings.BurstRateCount, "Burst rate count")

	flag.Parse()

	Settings.RedisAddr = fmt.Sprintf("%s:%d", Settings.RedisHost, Settings.RedisPort)
}
