package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"site-map/parser"
)

type results struct {
	numberOfOccurrences int
}

func main() {
	baseTarget := flag.String("target", "https://www.calhoun.io", "Base domain to start processing links from")
	flag.Parse()

	// Find the domain name from the base target
	u, err := url.Parse(*baseTarget)
	if err != nil {
		log.Fatalf("error whilst determining  the base domain: %s", err)
	}
	if u.Scheme == "" || u.Host == "" {
		log.Fatalf("target must have a scheme and a host e.g. https://www.calhoun.io")
	}
	targetDomain := u.Host
	targetScheme := u.Scheme

	log.Printf("Target: %s", *baseTarget)

	// Retrieve the page and parse HTML links
	resp, err := http.Get(*baseTarget)
	if err != nil {
		log.Fatalf("error whilst doing a HTTP GET against %s: %s", *baseTarget, err)
	}
	defer resp.Body.Close()

	rawLinks, err := parser.Parse(resp.Body)
	if err != nil {
		log.Fatalf("error whilst parsing links from %s: %s", *baseTarget, err)
	}
	log.Printf("Found %d raw links", len(rawLinks))

	// Extract links. Use a pointer, so we can update the map values in place
	links := make(map[string]*results)

	for _, link := range rawLinks {
		u, err = url.Parse(link.Href)
		if err != nil {
			log.Fatalf("error whilst parsing URL %s: %s", link.Href, err)
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

		// Add links to the map, so we have a list of unique ones whilst maintaining a duplicate count
		if _, found := links[fullURL]; !found {
			//log.Printf("Adding unique link to map: %s", fullURL)
			links[fullURL] = &results{numberOfOccurrences: 1}
			continue
		}
		links[fullURL].numberOfOccurrences++
	}

	// Print results
	log.Printf("Found %d links which match this domain. Results: ", len(links))
	for _, l := range sorted(links) {
		log.Print(l)
	}
}

func sorted(links map[string]*results) []string {
	ret := make([]string, 0, len(links))
	for k := range links {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}
