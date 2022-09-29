package routing

import (
	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/AksAman/gophercises/quietHN/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupFiberRoutes(app *fiber.App) {

	app.Use(middlewares.GetRateLimiterMiddleware())

	app.Get("/", controllers.Home)
	app.Get("/stories", controllers.GetStories)

	internalRoute := app.Group("/__internal__", func(c *fiber.Ctx) error {
		c.Set("private", "true")
		return c.Next()
	})
	internalRoute.Get("/monitor", middlewares.GetMonitorMiddleware())
}
