package views

import (
	"time"

	"github.com/AksAman/gophercises/quietHN/models"
)

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
