package controllers

import (
	"fmt"

	"github.com/AksAman/gophercises/quietHN/devtools"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/utils"
	"github.com/gofiber/fiber/v2"
)

func ViewSource(c *fiber.Ctx) error {
	filePath := c.Query("path")
	if filePath == "" {
		return c.Status(fiber.StatusBadRequest).SendString("path query param is required")
	}
	lineNumber := utils.GetQueryParam(c, "line", -1)

	htmlString := fmt.Sprintf(`
	<div id="explanation" style="padding: 10px; background: #282c34; color: #f00;">
            Debug Mode is ON
    </div>
	<h4 style="margin: 10px">Source code for %s</h4>
`, filePath)

	highlightedCode, err := devtools.GetHighlightedSourceCode(filePath, lineNumber, settings.Settings.StackTraceTheme, -1)

	htmlString += highlightedCode
	htmlString += `<style>
		body {
			background-color: #111;
			margin: 0;
			color : #fff;
			font-family: monospace;
		}
	</style>
	`

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	c.Response().Header.Set("Content-Type", "text/html")
	return c.SendString(htmlString)

}
