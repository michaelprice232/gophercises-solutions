package parser

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

// ParseLinks recursively parses an HTML tree for links and extracts any text from the child elements.
// Expects a html.Node pointer parsed from the x/net/html module.
func ParseLinks(rootNode *html.Node) []Link {
	var f func(*html.Node)
	links := make([]Link, 0)

	var linkValue string
	var textForLink string

	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {

			textForLink = ""
			linkValue = hrefValue(n.Attr)

			// Process only direct children, looking for text
			for ic := n.FirstChild; ic != nil; ic = ic.Parent.NextSibling {
				if ic.Type == html.TextNode {
					// Remove newlines and whitespace
					d := strings.TrimSpace(ic.Data)

					// Remove double spaces
					d = strings.ReplaceAll(d, "  ", " ")

					textForLink += d
				}
			}

			links = append(links, Link{Href: linkValue, Text: textForLink})
		}

		// Recursive through all nodes in the tree
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(rootNode)

	return links
}

// hrefValue retrieves the href value from a slice of attributes returned for a html.Node.
func hrefValue(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == "href" {
			return a.Val
		}
	}
	return "<not-found>"
}

// ParseHTML returns a node tree based on an HTML source file.
func ParseHTML(path string) (*html.Node, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	reader := bytes.NewReader(b)

	doc, err := html.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML: %v", err)
	}

	return doc, nil
}
