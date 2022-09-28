package ratelimiter

import (
	"fmt"
	"time"
)

type IRateLimiter interface {
	Wait() time.Time
}

type NormalRateLimiter struct {
	interval time.Duration
	ticker   *time.Ticker
	C        chan time.Time
}

func (rl *NormalRateLimiter) Wait() time.Time {
	return <-rl.C
}

func (rl *NormalRateLimiter) run() {
	for {
		fmt.Println("\t--------- waiting for rate limiter ---------")
		rl.C <- <-rl.ticker.C
	}
}

func NewRateLimiter(interval time.Duration) (*NormalRateLimiter, error) {
	rl := &NormalRateLimiter{
		interval: interval,
		ticker:   time.NewTicker(interval),
		C:        make(chan time.Time),
	}

	go rl.run()
	return rl, nil
}
