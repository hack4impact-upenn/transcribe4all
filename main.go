package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	_ "net/http/pprof" // import for side effects

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello/{name}", helloHandler)
	r.HandleFunc("/add_job", jobHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}

func jobHandler(w http.ResponseWriter, r *http.Request) {
	var d audioData

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(d)
}

type audioData struct {
	AudioURL string `json:"audioURL"`
}
