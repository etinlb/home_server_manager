package main

import (
	"fmt"
	"net/http"
)

// Adapter wraps an http.Handler with additional
// functionality.
type Adapter func(http.Handler) http.Handler

func Logging() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(r.Method, r.URL.Path)
			fmt.Println("IN HERERE")
			h.ServeHTTP(w, r)
		})
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

// Adapt h with all specified adapters.
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// WithHeader is an Adapter that sets an HTTP handler.
func WithHeader(key, value string) Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add(key, value)
			h.ServeHTTP(w, r)
		})
	}
}

// SupportXHTTPMethodOverride adds support for the X-HTTP-Method-Override
// header.
func SupportXHTTPMethodOverride() Adapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m := r.Header.Get("X-HTTP-Method-Override")
			if len(m) > 0 {
				r.Method = m
			}
			h.ServeHTTP(w, r)
		})
	}
}

func testHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("in renmae route")

	// send response
	w.Write([]byte("Hello"))
}

func otherHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Other")

	// send response
	w.Write([]byte("other"))
}

func main() {
	mux := http.NewServeMux()
	// router := mux.NewRouter()

	// adapt a single route
	mux.Handle("/test", Adapt(http.HandlerFunc(testHandler), WithHeader("X-Something", "Specific")))
	mux.Handle("/other/blah", Adapt(http.HandlerFunc(otherHandler), WithHeader("X-Something", "Specific")))

	// adapt all handlers
	// mux.Handle("/", Adapt(mux,
	//  SupportXHTTPMethodOverride(),
	//  WithHeader("Server", "MyApp v1"),
	//  Logging(),
	// ))
	// http.ListenAndServe(":17901", Log(mux))
	http.ListenAndServe(":17901", Adapt(mux,
		SupportXHTTPMethodOverride(),
		WithHeader("Server", "MyApp v1"),
		Logging(),
	))
}
