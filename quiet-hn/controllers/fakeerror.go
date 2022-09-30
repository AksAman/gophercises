package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ProblematicFunc() {
	panic(fmt.Errorf("some Error"))
}

func FakeError(c *fiber.Ctx) error {
	ProblematicFunc()
	return nil
}

func FakeErrorAfter(c *fiber.Ctx) error {
	c.SendString("Hello")
	ProblematicFunc()
	return nil
}
