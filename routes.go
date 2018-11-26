package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Return a http.ServeMux with all the routes bound to it
func register_routes() *mux.Router {
	// r := mux
	r := mux.NewRouter()

	wapi := r.PathPrefix("/wapi").Subrouter()
	wapi.HandleFunc("/", websocketApiHandler)

	// TODO: The static router should be done in nginx I think
	staticFileHandler := getHTTPFileHandler("dist/")
	r.PathPrefix("/static/").Handler(staticFileHandler)

	// Server up the index and favicon.ico
	r.Handle("/", staticFileHandler)

	return r
}

// Registers the directory string as a generic file handler. Checks if the
// directory that we are trying to serve exists relative to the current
// directory  we are in.
// TODO: Bubble up error
func getHTTPFileHandler(directory string) http.Handler {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Everything checks out, server up the static page
	fs := http.Dir(dir + "/" + directory)
	fileHandler := http.FileServer(fs)
	return fileHandler
}
