package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello/{name}", helloHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}
