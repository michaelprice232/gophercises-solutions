package main

import (
	"fmt"
	"net/http"

	"urlshort_lesson02"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/test1": "https://mike-price.com",
		"/test2": "https://godoc.org/github.com/gophercises/urlshort",
		"/test3": "https://godoc.org/gopkg.in/yaml.v2",
		"/test4": "https://jayne-price.com",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	//	// Build the YAMLHandler using the mapHandler as the fallback
	//	yaml := `
	//- path: /urlshort
	//  url: https://github.com/gophercises/urlshort
	//- path: /urlshort-final
	//  url: https://github.com/gophercises/urlshort/tree/solution
	//`
	//	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	//	if err != nil {
	//		panic(err)
	//	}
	fmt.Println("Starting the server on :8080")
	_ = http.ListenAndServe(":8080", mapHandler)
	//_ = http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello, world!")
}
