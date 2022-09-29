package controllers

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var t time.Time = time.Now()

func Home(w http.ResponseWriter, r *http.Request) {
	if rateLimiter != nil {
		rateLimiter.Wait()
		// fmt.Println("\t--------- got time from rate limiter ---------", t)
	}
	atomic.AddUint64(&counter, 1)

	w.Header().Set("Content-Type", "application/json")
	t = time.Now()

	fmt.Fprintf(w,
		`
		{
			"message": "Hello from net/http",
			"ip": "%v",
			"visited": %d
		}
		`,
		r.RemoteAddr,
		counter,
	)
}
