package main

import (
	"fmt"
	"log"

	"github.com/AksAman/gophercises/quietHN/middlewares"
	"github.com/AksAman/gophercises/quietHN/routing"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/gofiber/fiber/v2"
)

func RunServer() {

	app := fiber.New(
		fiber.Config{
			Views: views.GetFiberViews(),
		},
	)
	middlewares.SetupFiberMiddlewares(app)

	routing.SetupFiberRoutes(app)

	addr := fmt.Sprintf(":%d", settings.Settings.Port)
	log.Fatal(app.Listen(addr))
}

func main() {
	RunServer()
}
