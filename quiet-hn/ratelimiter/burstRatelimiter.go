package ratelimiter

import (
	"fmt"
	"time"
)

type BurstyRateLimiter struct {
	NormalRateLimiter
	currentBurstCount int
	burstCount        int
}

func (brl *BurstyRateLimiter) resetBurst() {
	fmt.Println("\t--------- reseting burst --------- to ", brl.burstCount, "from", brl.currentBurstCount)

	brl.currentBurstCount = brl.burstCount
	brl.ticker.Reset(brl.interval)
}

func (rl *BurstyRateLimiter) run() {
	// this is bursty rate limiter which permits `burstCount` events every `rl.interval` interval
	for {
		if rl.currentBurstCount <= 0 {
			// wait for ticker
			fmt.Println("\t--------- waiting for rate limiter ---------", rl.interval)
			<-rl.ticker.C
			rl.resetBurst()
		}
		rl.C <- time.Now()
		rl.currentBurstCount--
	}
}

func NewBurstyRateLimiter(interval time.Duration, burstCount int) (*BurstyRateLimiter, error) {

	if burstCount <= 0 {
		return nil, fmt.Errorf("burst count should be a positive int")
	}

	brl := &BurstyRateLimiter{
		NormalRateLimiter: NormalRateLimiter{
			interval: interval,
			ticker:   time.NewTicker(interval),
			C:        make(chan time.Time),
		},
		burstCount:        burstCount,
		currentBurstCount: burstCount,
	}
	go brl.run()
	return brl, nil
}
