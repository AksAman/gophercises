package linkparser

import (
	"errors"
	"io"

	"golang.org/x/net/html"
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

func Parse(r io.Reader) ([]Link, error) {
	links := []Link{}

	node, err := html.Parse(r)
	if err != nil {
		return []Link{}, err
	}

	parseNodes(node, &links)

	return links, nil
}

func parseNodes(node *html.Node, links *[]Link) {
	// Do something with n...
	if link, err := parseNode(node); err == nil {
		*links = append(*links, link)
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		parseNodes(child, links)
	}
}

func parseNode(node *html.Node) (Link, error) {
	if node.Type == html.ElementNode && node.Data == "a" {
		var href string
		for _, attr := range node.Attr {
			if attr.Key == "href" {
				href = attr.Val
				break
			}
		}
		textData := findTextOfAnchorNode(node)

		if href != "" && textData != "" {
			return Link{Href: href, Text: textData}, nil
		} else {
			return Link{}, errors.New("href attr not found in a tag")
		}
	}
	return Link{}, errors.New("not an element node")
}

func findTextOfAnchorNode(node *html.Node) string {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			return child.Data
		}
	}
	return ""
}
