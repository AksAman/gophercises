package views

import (
	"html/template"
	"time"

	"github.com/AksAman/gophercises/quietHN/models"
	"github.com/gofiber/template/html"
)

func GetFiberViews() *html.Engine {
	templateEngine := html.New("./templates", ".gohtml")

	return templateEngine
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

type ErrorTemplateContext struct {
	StatusCode int
	Message    string
	StackTrace template.HTML
	Debug      bool
}
