package views

import (
	"html/template"
	"time"

	"github.com/AksAman/gophercises/quietHN/settings"
)

var (
	storyTemplate *template.Template
	cache         storyCache
)

func init() {
	initTemplates()
	initCache()
}

func initCache() {
	cache = storyCache{timeout: time.Duration(settings.Settings.Timeout)}
}

func initTemplates() {
	storyTemplate = template.Must(template.ParseFiles("templates/index.gohtml"))
}
