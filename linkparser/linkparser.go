package linkparser

import (
	"io"
)

/*
	Link{
	  Href: "/dog",
	  Text: "Something in a span Text not in a span Bold text!",
	}
*/
type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

func Parse(r io.Reader) []Link {
	links := []Link{}

	return links
}
