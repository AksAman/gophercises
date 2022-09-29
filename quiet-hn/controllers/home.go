package controllers

import (
	"net/http"
	"sync/atomic"

	"github.com/gin-gonic/gin"
)

func Home(c *gin.Context) {
	atomic.AddUint64(&counter, 1)

	// set headers
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello From Fiber Server",
		"ip":      c.ClientIP(),
		"visited": counter,
	})
}
