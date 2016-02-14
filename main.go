package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // import for side effects

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello/{name}", helloHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

	http.HandleFunc("/health", healthHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}
