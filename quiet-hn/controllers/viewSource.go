package controllers

import (
	"github.com/AksAman/gophercises/quietHN/devtools"
	"github.com/AksAman/gophercises/quietHN/utils"
	"github.com/gofiber/fiber/v2"
)

func ViewSource(c *fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		return c.Status(fiber.StatusBadRequest).SendString("path query param is required")
	}
	lineNumber := utils.GetQueryParam(c, "line", -1)

	highlightedCode, err := devtools.GetHighlightedSourceCode(filePath, lineNumber, "monokai", -1)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	c.Response().Header.Set("Content-Type", "text/html")
	return c.SendString(highlightedCode)

}
