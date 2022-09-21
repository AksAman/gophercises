package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/AksAman/gophercises/quietHN/hnclient"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	client := hnclient.Client{}
	topStories, err := client.GetTopItems()
	if err != nil {
		panic(err)
	}
	fmt.Printf("topStories: %v\n", topStories[:50])

	// generate a random id
	randomIdx := rand.Intn(len(topStories))
	randomID := topStories[randomIdx]

	item, err := client.GetItem(randomID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("item:  %v\n", item.String())

}
