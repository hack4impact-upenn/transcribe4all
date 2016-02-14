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
	r.HandleFunc("/add_job", initiateTranscriptionJobHandler)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}

// initiateTranscriptionJobHandle takes a POST request containing a json object,
// decodes it into an audioData struct, and returns the struct json-encoded.
func initiateTranscriptionJobHandler(w http.ResponseWriter, r *http.Request) {
	var d transcriptionJobData

	// unmarshal from the response body directly into our struct
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the struct encoded as json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(d)
}

type transcriptionJobData struct {
	AudioURL string `json:"audioURL"`
}
