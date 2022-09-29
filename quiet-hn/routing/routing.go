package routing

import (
	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/gin-gonic/gin"
)

func SetupGinRoutes(app *gin.Engine) {
	app.GET("/", controllers.Home)
	app.GET("/stories", controllers.GetStories)
}
