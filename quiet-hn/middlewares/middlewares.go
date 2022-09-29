package middlewares

import (
	"time"

	"github.com/AksAman/gophercises/quietHN/ratelimiter"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/labstack/echo/v4"
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

func SetupEchoMiddlewares(app *echo.Echo) {
	app.Use(rateLimiterMiddleware)
}

func rateLimiterMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if rateLimiter != nil {
			rateLimiter.Wait()
		}

		return next(c)
	}
}
