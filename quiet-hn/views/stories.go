package views

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/AksAman/gophercises/quietHN/caching"
	"github.com/AksAman/gophercises/quietHN/hnclient"
	"github.com/AksAman/gophercises/quietHN/models"
	"github.com/AksAman/gophercises/quietHN/settings"
	"github.com/AksAman/gophercises/quietHN/utils"
)

type storiesTemplateContext struct {
	Strategy      string
	RequiredCount int
	Stories       []*models.HNItem
	Latency       time.Duration
	TotalLatency  time.Duration
}

type storyChanResult struct {
	index int
	story *models.HNItem
	err   error
}

type storyGetterStrategy func(client *hnclient.Client, ids []int) []*models.HNItem

func (c *storiesTemplateContext) CalculateTotalLatency() {
	var totalLatency time.Duration
	for _, story := range c.Stories {
		totalLatency += story.Latency
	}
	c.TotalLatency = totalLatency
}

func Stories(w http.ResponseWriter, r *http.Request) {
	if rateLimiter != nil {
		rateLimiter.Wait()
	}
	fmt.Println("\nGETTING STORIES")

	requiredStoriesCount := utils.GetQueryParam(r, "n", settings.Settings.MaxStories)
	if requiredStoriesCount > 0 && requiredStoriesCount > settings.Settings.MaxStories {
		requiredStoriesCount = settings.Settings.MaxStories
	}

	getStrategy := getStoriesForIDsAsync
	strategyName := "Async"

	if strings.ToLower(utils.GetQueryParam(r, "strategy", "async")) != "async" {
		getStrategy = getStoriesForIDsSync
		strategyName = "Sync"
	}

	start := time.Now()

	// Actual code goes here

	stories, err := getStories(requiredStoriesCount, getStrategy, cache)
	if err != nil {
		http.Error(w, "Error fetching top stories", http.StatusInternalServerError)
		return
	}
	templateContext := storiesTemplateContext{
		RequiredCount: requiredStoriesCount,
		Stories:       stories,
		Latency:       time.Since(start).Round(time.Nanosecond),
		Strategy:      strategyName,
	}
	templateContext.CalculateTotalLatency()

	err = storyTemplate.Execute(w, templateContext)

	if err != nil {
		http.Error(w, "Failed to process the template", http.StatusInternalServerError)
		return
	}

}

func getStories(requiredStoriesCount int, getStrategy storyGetterStrategy, cache caching.Cache[models.HNItem]) ([]*models.HNItem, error) {
	fmt.Printf("Getting stories count:%d\n", requiredStoriesCount)
	var stories []*models.HNItem
	client := hnclient.Client{}

	cachedStories := cache.Get()
	if cachedStories != nil {
		if len(cachedStories) >= requiredStoriesCount {
			stories = cachedStories[:requiredStoriesCount]
		} else {
			stories = append(stories, cachedStories...)
		}
	}

	if len(stories) < requiredStoriesCount {
		currentStartIdx := len(stories)
		ids, err := client.GetTopItems()
		if err != nil {
			return nil, err
		}
		for len(stories) < requiredStoriesCount {
			needed := (requiredStoriesCount - len(stories)) * 2
			fmt.Printf("have %d stories, need %d more to reach %d stories\n", len(stories), needed, requiredStoriesCount)
			nextIDs := ids[currentStartIdx : currentStartIdx+needed]

			newStories := getStrategy(&client, nextIDs)
			stories = append(stories, newStories...)
			currentStartIdx += needed
		}

		stories = stories[:requiredStoriesCount]
		cache.Set(stories)
	}
	return stories, nil
}

func getStoriesForIDsAsync(client *hnclient.Client, ids []int) []*models.HNItem {
	time.Sleep(time.Second)
	fmt.Println("\tgetStoriesForIDsAsync", len(ids), ids)
	storyChan := make(chan storyChanResult)

	for index, id := range ids {
		go func(index, id int) {
			storyStart := time.Now()
			item, err := client.GetItem(id)
			if err != nil {
				log.Printf("\t\t\tError fetching story with id %d: %v", id, err)
				storyChan <- storyChanResult{index: index, err: err}
				return
			}
			if !item.IsStory() {
				fmt.Printf("\t\t\tItem with id %d is not a valid story\n", id)
				storyChan <- storyChanResult{err: fmt.Errorf("item with id %d is not a valid story", id)}
				return
			}
			item.Latency = time.Since(storyStart).Round(time.Nanosecond)

			storyChan <- storyChanResult{index: index, story: &item}
		}(index, id)
	}

	var chanResults []storyChanResult

	for i := 0; i < len(ids); i++ {
		result := <-storyChan
		if result.err != nil {
			continue
		}
		fmt.Printf("\t\t%d: Found result.index:%d from id:%d\n", i, result.index, result.story.ID)
		chanResults = append(chanResults, result)
	}

	fmt.Println("\t\tFound", len(chanResults), "stories")

	sort.Slice(chanResults, func(i, j int) bool {
		return chanResults[i].index < chanResults[j].index
	})

	var stories []*models.HNItem
	for _, result := range chanResults {
		stories = append(stories, result.story)
	}

	return stories
}

func getStoriesForIDsSync(client *hnclient.Client, ids []int) []*models.HNItem {
	var stories []*models.HNItem

	for _, id := range ids {
		storyStart := time.Now()
		item, err := client.GetItem(id)
		if err != nil {
			log.Printf("\tError fetching story with id %d: %v", id, err)
			continue
		}
		if !item.IsStory() {
			// fmt.Printf("\t\tItem with id %d is not a valid story\n", id)
			continue
		}
		item.Latency = time.Since(storyStart).Round(time.Millisecond)

		stories = append(stories, &item)
	}
	return stories
}
