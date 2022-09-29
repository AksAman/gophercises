package routing

import (
	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetupFiberRoutes(app *fiber.App) {
	app.Get("/", controllers.Home)
	app.Get("/stories", controllers.GetStories)
}
