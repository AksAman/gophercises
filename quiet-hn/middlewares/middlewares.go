package middlewares

import (
	"time"

	"github.com/AksAman/gophercises/quietHN/ratelimiter"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/gin-gonic/gin"
)

var (
	rateLimiter ratelimiter.IRateLimiter
)

func init() {
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

func SetupGinMiddlewares(app *gin.Engine) {
	app.Use(rateLimiterMiddleware())
}

func rateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if rateLimiter != nil {
			rateLimiter.Wait()
		}

		c.Next()
	}
}
