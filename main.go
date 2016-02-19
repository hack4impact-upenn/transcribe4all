package main

import (
	"net/http"
	_ "net/http/pprof" // import for side effects
)

func main() {
	router := NewRouter()
	middlewareRouter := ApplyMiddleware(router)

	// serve http
	http.Handle("/", middlewareRouter)
	http.ListenAndServe(":8080", nil)
}
