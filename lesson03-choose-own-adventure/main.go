package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	adventureFilePath = "story.json"
	defaultArc        = "intro"
)

type Arc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options"`
}

// handler is used to serve HTTP requests and satisfies http.Handler
type handler struct {
	adventure map[string]*Arc
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Raw request path with leading slash
	rawRequestPath := r.URL.Path

	// Send user to intro if no arc is present in the path
	requestPath := ""
	if rawRequestPath == "/" {
		requestPath = defaultArc
	} else {
		// Remove the leading slash and any additional sub paths
		requestPath = strings.Split(rawRequestPath, "/")[1]
	}
	log.Printf("Raw URL: %s, Extracted arc: %s", rawRequestPath, requestPath)

	// Load the title if the request path correlates with an arc name
	if data, found := h.adventure[requestPath]; found {
		err := renderPage(w, data)
		if err != nil {
			log.Printf("error rendering page: %s", err)
		}
		return
	}

	// Not found arc
	w.WriteHeader(http.StatusNotFound)
	_, err := w.Write([]byte("<h1>Arc not found</h1>"))
	if err != nil {
		log.Printf("error writing HTTP response: %v", err)
	}
}

func renderPage(w http.ResponseWriter, data *Arc) error {
	bodyBytes, err := os.ReadFile("templates/page.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("error reading template source file: %s", err)
	}

	tmpl, err := template.New("page").Parse(string(bodyBytes))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("error parsing HTML template: %s", err)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("error executing HTML template: %s", err)
	}

	return nil
}

func loadAdventure(sourceFilepath string) (map[string]*Arc, error) {
	bodyBytes, err := os.ReadFile(sourceFilepath)
	if err != nil {
		return nil, fmt.Errorf("error reading file '%s': %s", sourceFilepath, err)
	}

	adventure := make(map[string]*Arc)
	err = json.Unmarshal(bodyBytes, &adventure)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json: %s", err)
	}

	return adventure, nil
}

func main() {
	adventure, err := loadAdventure(adventureFilePath)
	if err != nil {
		log.Fatalf("loading adventure: %v", err)
	}
	log.Printf("Loaded %d stories", len(adventure))

	mux := http.NewServeMux()
	mux.Handle("/", handler{adventure: adventure})

	log.Printf("HTTP server listening on port 8080")
	if err = http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
