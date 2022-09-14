package linkparser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

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

var nodeTypes map[html.NodeType]string = make(map[html.NodeType]string)

func init() {
	nodeTypes[html.ErrorNode] = "Type(ErrorNode)"
	nodeTypes[html.TextNode] = "Type(TextNode)"
	nodeTypes[html.DocumentNode] = "Type(ElementNode)"
	nodeTypes[html.ElementNode] = "Type(ElementNode)"
	nodeTypes[html.CommentNode] = "Type(CommentNode)"
	nodeTypes[html.DoctypeNode] = "Type(DoctypeNode)"
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

func ParseBytes(data []byte) ([]Link, error) {
	links := []Link{}
	node, err := html.Parse(bytes.NewReader(data))
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
		textData := findTextOfAnchorNode(node, "")
		textData = cleanString(textData)
		if href != "" && textData != "" {
			return Link{Href: href, Text: textData}, nil
		} else {
			return Link{}, errors.New("href attr not found in a tag")
		}
	}
	return Link{}, errors.New("not an element node")
}

func findTextOfAnchorNode(node *html.Node, padding string) string {

	var text string
	// printNode(node, padding)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			text += child.Data
		} else {
			text += " " + findTextOfAnchorNode(child, padding+"  ")
		}
	}
	// fmt.Printf("%v text: %v\n", padding+" -", text)
	return text
}

func printNode(n *html.Node, padding string) {
	fmt.Printf("%v node data: %q, of %v\n", padding, cleanString(n.Data), nodeTypes[n.Type])
}

func cleanString(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
