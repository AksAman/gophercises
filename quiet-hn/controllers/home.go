package controllers

import (
	"net/http"
	"sync/atomic"

	"github.com/labstack/echo/v4"
)

func Home(c echo.Context) error {
	atomic.AddUint64(&counter, 1)

	// set headers
	return c.JSON(http.StatusOK, map[string]any{
		"message": "Hello From ECHO Server",
		"ip":      c.RealIP(),
		"visited": counter,
	})
}
