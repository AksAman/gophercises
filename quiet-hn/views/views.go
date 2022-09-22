package views

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/AksAman/gophercises/quietHN/hnclient"
	"github.com/AksAman/gophercises/quietHN/models"
)

var (
	maxStories    int
	storyTemplate *template.Template
)

func init() {
	flag.IntVar(&maxStories, "n", 30, "Number of stories to fetch")
	flag.Parse()
	initTemplates()
}

type storiesTemplateContext struct {
	Stories []models.HNItem
	Latency time.Duration
}

func initTemplates() {
	storyTemplate = template.Must(template.ParseFiles("templates/index.gohtml"))
}

type NotFoundHandler struct{}

func (h NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "Nothing here bruh!")
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Home!")
}

func Stories(w http.ResponseWriter, r *http.Request) {
	nStoriesInt := maxStories
	queryParams := r.URL.Query()
	nStoriesParam := queryParams.Get("n")

	if n, err := strconv.Atoi(nStoriesParam); err == nil && n > 0 && n <= maxStories {
		nStoriesInt = n
	}

	start := time.Now()

	// Actual code goes here
	client := hnclient.Client{}
	ids, err := client.GetTopItems()
	if err != nil {
		http.Error(w, "Error fetching top stories", http.StatusInternalServerError)
		return
	}
	fmt.Printf("nStoriesInt: %v\n", nStoriesInt)

	requiredIds := ids[:nStoriesInt]
	var stories []models.HNItem

	for _, id := range requiredIds {
		storyStart := time.Now()
		story, err := client.GetItem(id)
		if err != nil {
			continue
		}
		story.Latency = time.Since(storyStart).Round(time.Millisecond)
		stories = append(stories, story)
	}

	templateContext := storiesTemplateContext{
		Stories: stories,
		Latency: time.Since(start).Round(time.Millisecond),
	}

	err = storyTemplate.Execute(w, templateContext)

	if err != nil {
		http.Error(w, "Failed to process the template", http.StatusInternalServerError)
		return
	}

}
