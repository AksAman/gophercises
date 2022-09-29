package controllers

import (
	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		statusCode = e.Code
	}

	err = c.Status(statusCode).Render("error", views.ErrorTemplateContext{StatusCode: statusCode, Message: err.Error()})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}
