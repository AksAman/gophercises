package controllers

import (
	"html/template"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

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
		err = c.Status(statusCode).Render("error", views.ErrorTemplateContext{
			Path:          c.Path(),
			Method:        c.Method(),
			URL:           c.BaseURL() + c.Path(),
			ServerTime:    time.Now().Format(time.RFC1123Z),
			StatusCode:    statusCode,
			Message:       strings.TrimSpace(err.Error()),
			StackTrace:    template.HTML(stackTrace),
			Debug:         settings.Settings.Debug,
			GolangVersion: runtime.Version(),
		})

	} else {
		stackTrace = ""
		err = c.Status(statusCode).Render("error", views.ErrorTemplateContext{
			Debug:      settings.Settings.Debug,
			Path:       c.Path(),
			StatusCode: statusCode,
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return nil
}
