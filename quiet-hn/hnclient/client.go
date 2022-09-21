package hnclient

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AksAman/gophercises/quietHN/models"
)

const (
	BASE_URL string = "https://hacker-news.firebaseio.com/v0"
)

type Client struct {
	baseURL string
}

// defaultify : client code can call this method to fill up the zero value of the client
func (c *Client) defaultify() {
	if c.baseURL == "" {
		c.baseURL = BASE_URL
	}
}

func (c *Client) getTopStoriesEndpoint() string {
	return fmt.Sprintf("%s/topstories.json", c.baseURL)
}

func (c *Client) getStoryEndpoint(id int) string {
	return fmt.Sprintf("%s/item/%d.json", c.baseURL, id)
}

func (c *Client) GetTopItems() ([]int, error) {
	c.defaultify()

	resp, err := http.Get(c.getTopStoriesEndpoint())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&ids)
	if err != nil {
		return nil, err
	}

	return ids, err
}

func (c *Client) GetItem(id int) (models.HNItem, error) {
	c.defaultify()

	var item models.HNItem
	endpoint := c.getStoryEndpoint(id)

	resp, err := http.Get(endpoint)
	if err != nil {
		return item, err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&item)

	if err != nil {
		return item, err
	}

	return item, nil

}
