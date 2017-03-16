package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Return a http.ServeMux with all the routes bound to it
func register_routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/rename", renameHandler)
	mux.HandleFunc("/api/", messageHandler)

	staticFileHandler := getHTTPFileHandler("static", "/")
	mux.Handle("/static", staticFileHandler)
	// registerFileHandlerRoute("static", "/", mux)

	return mux
}

// Registers the directory string as a generic file handler. Checks if the
// directory that we are trying to serve exists relative to the current
// directory  we are in.
// TODO: Bubble up error
func getHTTPFileHandler(directory string) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	// Everything checks out, server up the static page
	fs := http.Dir(dir + "/" + directory)
	fileHandler := http.FileServer(fs)
	return fileHandler
	// mux.Handle(pattern, fileHandler)
}
