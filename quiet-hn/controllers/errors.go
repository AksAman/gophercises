package controllers

import (
	"html/template"
	"runtime/debug"

	"github.com/AksAman/gophercises/quietHN/devtools"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	statusCode := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		statusCode = e.Code
	}

	var stackTrace string
	if settings.Settings.Debug {
		stackTrace = devtools.ParseStackTraceToHTML(string(debug.Stack()))
	} else {
		stackTrace = ""
	}

	err = c.Status(statusCode).Render("error", views.ErrorTemplateContext{
		StatusCode: statusCode,
		Message:    err.Error(),
		StackTrace: template.HTML(stackTrace),
		Debug:      settings.Settings.Debug,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}
