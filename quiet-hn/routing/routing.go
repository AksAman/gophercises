package routing

import (
	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(app *echo.Echo) {
	app.GET("/", controllers.Home)
	app.GET("/stories", controllers.GetStories)
}
