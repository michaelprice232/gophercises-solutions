package main

import (
	"flag"
	"log"
	"os"

	"link-parser"
)

func main() {
	filePath := flag.String("file-path", "./testdata/ex1.html", "Path to the HTML file to parse")
	flag.Parse()

	doc, err := parser.ParseHTML(*filePath)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	links := parser.ParseLinks(doc)

	if len(links) == 0 {
		log.Printf("No links found in %s", *filePath)
		os.Exit(0)
	}

	log.Printf("Links found in %s:", *filePath)
	for _, link := range links {
		log.Printf("URL: %s, Text: %s", link.Href, link.Text)
	}
}
