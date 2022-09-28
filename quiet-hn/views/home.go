package views

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

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Hello %v, visited %d times, time: %v", r.RemoteAddr, counter, time.Since(t).Round(time.Second))
	t = time.Now()
}
