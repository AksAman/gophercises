package middlewares

import (
	"time"

	"github.com/AksAman/gophercises/quietHN/ratelimiter"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
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

func GetRateLimiterMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if rateLimiter != nil {
			rateLimiter.Wait()
		}
		return c.Next()
	}
}

func GetMonitorMiddleware() fiber.Handler {
	return monitor.New(
		monitor.Config{
			Title: "Monitor for QuietHN",
		},
	)
}
