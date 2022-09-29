package main

import (
	"fmt"

	"github.com/AksAman/gophercises/quietHN/middlewares"
	"github.com/AksAman/gophercises/quietHN/routing"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/gin-gonic/gin"
)

func RunServer() {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.LoadHTMLGlob("templates/*")
	middlewares.SetupGinMiddlewares(app)
	routing.SetupGinRoutes(app)

	addr := fmt.Sprintf(":%d", settings.Settings.Port)
	app.Run(addr)
}

func main() {
	RunServer()
}
