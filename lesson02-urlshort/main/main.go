package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"urlshort_lesson02"
)

const httpPort = 8080

func main() {
	yamlConfigPath := flag.String("yaml-config", "", "path to a YAML config file e.g. ./config/config.yaml")
	jsonConfigPath := flag.String("json-config", "", "path to a JSON config file e.g. ./config/config.json")
	readFromDatabase := flag.Bool("read-from-db", false, "read from database instead of YAML/JSON config file")
	flag.Parse()

	var err error
	var config []byte
	if !*readFromDatabase {
		config, err = parseConfigFromFile(yamlConfigPath, jsonConfigPath)
		if err != nil {
			slog.Error(fmt.Sprintf("unable to load config from file: %v", err))
			os.Exit(1)
		}
	}

	mux := defaultMux()

	// Build the MapHandler using the default mux as the fallback
	pathsToUrls := map[string]string{
		"/test1": "https://mike-price.com",
		"/test2": "https://godoc.org/github.com/gophercises/urlshort",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler or JSONHandler using the mapHandler as the fallback
	var handler http.HandlerFunc
	handler, err = callHandler(yamlConfigPath, jsonConfigPath, readFromDatabase, config, mapHandler)
	if err != nil {
		slog.Error(fmt.Sprintf("error whilst reading config from file: %s", err))
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Starting the server on port %d", httpPort))
	err = http.ListenAndServe(fmt.Sprintf(":%d", httpPort), handler)
	if err != nil {
		slog.Error(fmt.Sprintf("error whilst running web server: %v", err))
		os.Exit(1)
	}
}

// callHandler determines which type of config file has been passed and calls the respective handler
func callHandler(ymlFile, jsonFile *string, dbFlag *bool, config []byte, fallback http.HandlerFunc) (http.HandlerFunc, error) {
	var handler http.HandlerFunc
	var err error

	if *dbFlag {
		handler, err = urlshort.DBHandler(fallback)
		if err != nil {
			return nil, err
		}

		return handler, nil
	}

	if *ymlFile != "" {
		handler, err = urlshort.YAMLHandler(config, fallback)
		if err != nil {
			return nil, err
		}
		return handler, err
	}
	if *jsonFile != "" {
		handler, err = urlshort.JSONHandler(config, fallback)
		if err != nil {
			return nil, err
		}
	}
	return handler, err
}

// defaultMux returns a default mux to be served when no other routes match the request path
func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Hello, world!")
}

// readFile returns a local file as a slice of bytes
func readFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %v", path, err)
	}
	return b, err
}

// parseConfigFromFile validates the config file parameters and then returns the respective file contents
func parseConfigFromFile(ymlFile, jsonFile *string) ([]byte, error) {
	if (*ymlFile == "" && *jsonFile == "") || (*ymlFile != "" && *jsonFile != "") {
		return nil, fmt.Errorf("exactly one config file must be set (yaml-config or json-config)")
	}

	var configPath string
	if *ymlFile != "" {
		configPath = *ymlFile
	} else {
		configPath = *jsonFile
	}

	config, err := readFile(configPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}
