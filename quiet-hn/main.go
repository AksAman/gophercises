package main

import (
	"fmt"
	"log"

	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/AksAman/gophercises/quietHN/routing"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/gofiber/fiber/v2"
)

func RunServer() {
	app := fiber.New(
		fiber.Config{
			Views:        views.GetFiberViews(),
			ErrorHandler: controllers.ErrorHandler,
			Prefork:      false,
		},
	)

	app.Static("/", "./static")

	routing.SetupFiberRoutes(app)

	addr := fmt.Sprintf(":%d", settings.Settings.Port)
	log.Fatal(app.Listen(addr))
}

func main() {
	RunServer()
}
