package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"site-map/parser"
)

const maxDepth = 3

// add XML export

func main() {
	baseTarget := flag.String("target", "https://www.calhoun.io", "Base domain to start processing links from")
	flag.Parse()

	// Validate base URL
	u, err := url.Parse(*baseTarget)
	if err != nil {
		log.Fatalf("error whilst parsing the target URL: %s", err)
	}
	if u.Scheme == "" || u.Host == "" {
		log.Fatalf("target must have a scheme and a host e.g. https://www.calhoun.io")
	}
	log.Printf("Base target: %s", *baseTarget)

	// Get links recursively using a breath first search
	l, err := bfs(*baseTarget, maxDepth)
	if err != nil {
		log.Fatalf("error whilst retrieving links: %v", err)
	}
	log.Printf("Found %d links", len(l))

	// Display links
	log.Print("All links:")
	for _, link := range l {
		log.Printf("Link: %s", link)
	}
}

// getLinks retrieves any HTML href links from a page as a slice
func getLinks(urlStr string) ([]string, error) {
	// Retrieve the page and parse HTML links
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}, fmt.Errorf("error whilst doing a HTTP GET against %s: %s", urlStr, err)
	}
	defer resp.Body.Close()

	// The base URL needs to be updated each time we follow a link
	targetDomain := resp.Request.URL.Host
	targetScheme := resp.Request.URL.Scheme

	rawLinks, err := parser.Parse(resp.Body)
	if err != nil {
		return []string{}, fmt.Errorf("error whilst parsing links from %s: %s", urlStr, err)
	}
	//log.Printf("Found %d raw links on %s", len(rawLinks), urlStr)

	// Collect relevant links
	links := make([]string, 0, len(rawLinks))

	for _, link := range rawLinks {
		u, err := url.Parse(link.Href)
		if err != nil {
			return []string{}, fmt.Errorf("error whilst parsing URL %s: %s", link.Href, err)
		}

		// Skip for non-standard URLs such as mailto
		if u.Host == "" && u.Path == "" {
			//log.Printf("Skipping non-standard link: %s", link.Href)
			continue
		}

		// Skip fragments (# references a subsection on the page)
		if u.Fragment != "" {
			//log.Printf("Skipping fragment: %s", link.Href)
			continue
		}

		// Sanitise URL for relative paths so everything is in FQDN format
		if u.Host == "" {
			u.Host = targetDomain
			u.Scheme = targetScheme
		}

		// Strip trailing slash on paths to avoid duplicates
		u.Path = strings.TrimSuffix(u.Path, "/")

		// Skip any links which do not belong to the target domain
		fullURL := u.String()
		if u.Host != targetDomain {
			//log.Printf("Skipping link as it is not in this domain: %s", fullURL)
			continue
		}
		//log.Printf("Adding link: %s", link)
		links = append(links, fullURL)
	}

	sort.Strings(links)
	return links, nil
}

// bfs implements a breath first search to recursively find HTML page links maxDepth levels deep.
// This function was copied as-is from the solution as I was not able to get it working on my own.
// Adding comments whilst re-writing it myself helped me to understand the algorithm better.
func bfs(urlStr string, maxDepth int) ([]string, error) {
	// Count number of HTTP calls made
	httpCalls := 0

	// Pages we have seen, so we only visit each once
	seen := make(map[string]struct{})

	// Current queue - links we are processing as part of this iteration / breath
	q := make(map[string]struct{})

	// Next queue - links we need to process on the next iteration / breath. Add the initial base URL
	nq := map[string]struct{}{
		urlStr: {},
	}

	// Limit the number of levels of links we will retrieve
	for i := 0; i <= maxDepth; i++ {
		log.Printf("Current level: %d", i)

		// Assign the nq to the q and initialise a new nq
		q, nq = nq, make(map[string]struct{})

		// If the queue is empty, we have finished
		if len(q) == 0 {
			break
		}

		// Process the current queue
		for href := range q {
			// Skip if the link has already been visited
			if _, ok := seen[href]; ok {
				continue
			}

			// Retrieve links from the current page
			seen[href] = struct{}{}
			httpCalls++
			links, err := getLinks(href)
			if err != nil {
				log.Printf("error getting links for %s: %s. Skipping page", href, err)
				continue
			}
			for _, link := range links {
				nq[link] = struct{}{}
			}
		}
	}
	log.Printf("Number of HTTP calls made: %d", httpCalls)

	// Return links as a sorted slice
	ret := make([]string, 0, len(seen))
	for link := range seen {
		ret = append(ret, link)
	}
	sort.Strings(ret)

	return ret, nil
}
