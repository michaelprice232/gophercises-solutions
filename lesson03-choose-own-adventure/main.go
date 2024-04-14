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
	htmlTemplate      = "./templates/page.html"
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
	adventure    map[string]*Arc
	htmlTemplate *template.Template
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

	// Render the page if the request path correlates with an arc name
	if data, found := h.adventure[requestPath]; found {
		err := h.renderPage(w, data)
		if err != nil {
			// Only print the real error to stdout
			log.Printf("error rendering page: %s", err)
		}
		return
	}

	// Not found arc
	http.Error(w, "Arc not found", http.StatusNotFound)
}

func (h handler) renderPage(w http.ResponseWriter, data *Arc) error {
	err := h.htmlTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
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

func loadHTMLTemplate(templatePath string) (*template.Template, error) {
	bodyBytes, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("error reading template source file from '%s': %s", templatePath, err)
	}

	tmpl, err := template.New("page").Parse(string(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("error parsing HTML template using source file '%s': %s", templatePath, err)
	}
	return tmpl, nil
}

func main() {
	adventure, err := loadAdventure(adventureFilePath)
	if err != nil {
		log.Fatalf("loading adventure: %v", err)
	}
	log.Printf("Loaded %d stories", len(adventure))

	tmpl, err := loadHTMLTemplate(htmlTemplate)
	if err != nil {
		log.Fatalf("loading HTML template: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler{adventure: adventure, htmlTemplate: tmpl})

	log.Printf("HTTP server listening on port 8080")
	if err = http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
