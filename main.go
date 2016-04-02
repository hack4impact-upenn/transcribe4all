// Package main initializes a web server.
package main

import (
	"net/http"
	_ "net/http/pprof" // import for side effects

	"github.com/hack4impact/transcribe4all/web"
)

func main() {
	router := web.NewRouter()
	middlewareRouter := web.ApplyMiddleware(router)

	// serve http
	http.Handle("/", middlewareRouter)
	http.Handle("/static/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8080", nil)

}
