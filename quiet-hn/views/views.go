package views

import (
	"html/template"
	"io"
	"time"

	"github.com/AksAman/gophercises/quietHN/models"
	"github.com/labstack/echo/v4"
)

type EchoTemplate struct {
	templates *template.Template
}

func (t *EchoTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func GetEchoTemplateRenderer() *EchoTemplate {
	templates := template.Must(template.ParseGlob("templates/*.gohtml"))
	return &EchoTemplate{
		templates: templates,
	}
}

type StoriesTemplateContext struct {
	Strategy      string
	RequiredCount int
	Stories       []*models.HNItem
	Latency       time.Duration
	TotalLatency  time.Duration
}

func (c *StoriesTemplateContext) CalculateTotalLatency() {
	var totalLatency time.Duration
	for _, story := range c.Stories {
		totalLatency += story.Latency
	}
	c.TotalLatency = totalLatency
}
