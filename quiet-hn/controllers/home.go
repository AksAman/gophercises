package controllers

import (
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	atomic.AddUint64(&counter, 1)

	// set headers
	return c.JSON(fiber.Map{
		"message": "Hello From Fiber Server",
		"ip":      c.IP(),
		"visited": counter,
	})
}
