package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var t time.Time = time.Now()

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HOME")
	if rateLimiter != nil {
		rateLimiter.Wait()
		// fmt.Println("\t--------- got time from rate limiter ---------", t)
	}
	atomic.AddUint64(&counter, 1)

	w.Header().Set("Content-Type", "application/json")
	t = time.Now()
	jsonData, err := json.MarshalIndent(
		map[string]interface{}{
			"message": "Hello from net/http",
			"ip":      r.RemoteAddr,
			"visited": counter,
		},
		"",
		"  ",
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonData)
}
