package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

// NewRouter creates and returns a mux.Router with default routes
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/hello/{name}", helloHandler)
	router.HandleFunc("/add_job", initiateTranscriptionJobHandler)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// ApplyMiddleware wraps the router in some middleware
func ApplyMiddleware(router http.Handler) http.Handler {
	stderrLoggingHandler := func(http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stderr, router)
	}
	middlewareRouter := alice.New(handlers.CompressHandler, stderrLoggingHandler).Then(router)
	return middlewareRouter
}
