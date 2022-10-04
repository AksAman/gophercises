package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func ProblematicFunc() {
	panic(errors.New("some Error"))
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
