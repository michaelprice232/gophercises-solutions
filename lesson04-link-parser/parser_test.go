package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func Test_hrefValue(t *testing.T) {
	attributes := []html.Attribute{
		{Key: "src", Val: "img.jpg"},
		{Key: "title", Val: "x"},
		{Key: "href", Val: "/link"},
	}

	response := hrefValue(attributes)
	assert.Equal(t, "/link", response, "Expected href value should be /link")

	attributes = []html.Attribute{
		{Key: "src", Val: "img.jpg"},
	}

	response = hrefValue(attributes)
	assert.Equal(t, "<not-found>", response, "Expected href value should be <not-found>")
}

func Test_ParseLinks(t *testing.T) {
	tt := []struct {
		name              string
		path              string
		expectedLinks     int
		expectedFirstLink string
		expectedFirstText string
		expectedLastLink  string
		expectedLastText  string
	}{
		{name: "3 links", path: "./testdata/ex1.html", expectedLinks: 3, expectedFirstLink: "/other-page", expectedFirstText: "A link to another page", expectedLastLink: "/bob", expectedLastText: "Bob's page"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			doc, err := ParseHTML(tc.path)
			assert.NoErrorf(t, err, "Expected no error parsing %s", tc.path)

			links := ParseLinks(doc)
			assert.Equal(t, tc.expectedLinks, len(links), "Expected %d links to be returned from %s, but got %d", tc.expectedLinks, tc.path, len(links))

			if tc.expectedLinks == 0 {
				assert.Equal(t, 0, len(links), "Expected no links to be returned, but got %d", len(links))
			}

			if tc.expectedLinks > 0 {
				assert.Greaterf(t, len(links), 0, "Expected > 0 links to be returned, but got %d", len(links))
				assert.Equal(t, tc.expectedFirstLink, links[0].Href, "Expected first link to be %s but got %s", tc.expectedFirstLink, links[0].Href)
				assert.Equal(t, tc.expectedFirstText, links[0].Text, "Expected first link to have text equal to %s but got %s", tc.expectedFirstText, links[0].Text)
			}

			if tc.expectedLinks > 1 {
				assert.Greaterf(t, len(links), 0, "Expected > 1 links to be returned, but got %d", len(links))

				// Calculate the last element
				last := len(links) - 1

				assert.Equal(t, tc.expectedLastLink, links[last].Href, "Expected last link to be %s but got %s", tc.expectedLastLink, links[0].Href)
				assert.Equal(t, tc.expectedLastText, links[last].Text, "Expected first link to have text equal to %s but got %s", tc.expectedLastText, links[0].Text)
			}

		})
	}
}
