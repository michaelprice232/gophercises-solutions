package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

// MapHandler will return a http.HandlerFunc (which also implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {

	mux := http.NewServeMux()

	for path, target := range pathsToUrls {
		// Avoid the inner function reading the last written outer variable
		targetInternal := target
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, targetInternal, http.StatusMovedPermanently)
		})
	}

	// Decide which mux to call based on whether the path is registered
	return func(w http.ResponseWriter, r *http.Request) {
		if _, ok := pathsToUrls[r.URL.Path]; ok {
			mux.ServeHTTP(w, r)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}

type redirect struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

type redirects []redirect

// YAMLHandler will parse the provided YAML and then return a http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding URL.
// If the path is not provided in the YAML, then the fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	red, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(red)

	return MapHandler(pathMap, fallback), nil
}

// parseYAML reads a YAML string and returns a redirects
func parseYAML(yml []byte) (redirects, error) {
	r := make(redirects, 0)

	err := yaml.Unmarshal(yml, &r)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %v", err)
	}

	return r, err
}

// buildMap returns redirects as a map, so we have a common data format for re-use
func buildMap(r redirects) map[string]string {
	m := make(map[string]string)
	for _, red := range r {
		m[red.Path] = red.URL
	}

	return m
}

// JSONHandler is the same as YAMLHandler although handles JSON source data
func JSONHandler(jsonStr []byte, fallback http.Handler) (http.HandlerFunc, error) {
	red, err := parseJSON(jsonStr)
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(red)

	return MapHandler(pathMap, fallback), nil
}

// parseJSON reads a JSON string and returns a redirects
func parseJSON(jsonStr []byte) (redirects, error) {
	r := make(redirects, 0)

	err := json.Unmarshal(jsonStr, &r)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %v", err)
	}

	return r, err
}
