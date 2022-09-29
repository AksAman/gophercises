package main

import (
	"fmt"

	"github.com/AksAman/gophercises/quietHN/middlewares"
	"github.com/AksAman/gophercises/quietHN/routing"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/labstack/echo/v4"
)

func RunServer() {

	e := echo.New()
	e.Renderer = views.GetEchoTemplateRenderer()
	middlewares.SetupEchoMiddlewares(e)
	routing.SetupRoutes(e)

	addr := fmt.Sprintf(":%d", settings.Settings.Port)
	e.Start(addr)
}

func main() {
	RunServer()
}
