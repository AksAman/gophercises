package routing

import (
	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/AksAman/gophercises/quietHN/middlewares"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/gofiber/fiber/v2"
)

func SetupFiberRoutes(app *fiber.App) {

	app.Use(middlewares.GetLoggerMiddleware())
	app.Use(middlewares.GetRateLimiterMiddleware())
	app.Use(middlewares.GetRecoveryMiddleware())

	app.Static("/", "./static")
	app.Get("/", controllers.Home)
	app.Get("/stories", controllers.GetStories)

	app.Get("/panic", controllers.FakeError)
	app.Get("/panic-after", controllers.FakeErrorAfter)

	if settings.Settings.Debug {
		internalRoute := app.Group("/__internal__", func(c *fiber.Ctx) error {
			c.Set("private", "true")
			return c.Next()
		})
		internalRoute.Get("/monitor", middlewares.GetMonitorMiddleware())
		internalRoute.Get("/view-source", controllers.ViewSource)
	}
}
