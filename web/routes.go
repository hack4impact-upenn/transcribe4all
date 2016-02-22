package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type transcriptionJobData struct {
	AudioURL       string   `json:"audioURL"`
	EmailAddresses []string `json:"emailAddresses"`
}

var routes = []route{
	route{
		"hello",
		"GET",
		"/hello/{name}",
		helloHandler,
	},
	route{
		"initiateTranscriptionJob",
		"POST",
		"/add_job",
		initiateTranscriptionJobHandler,
	},
	route{
		"health",
		"GET",
		"/health",
		healthHandler,
	},
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}

// initiateTranscriptionJobHandle takes a POST request containing a json object,
// decodes it into an audioData struct, and returns appropriate message.
func initiateTranscriptionJobHandler(w http.ResponseWriter, r *http.Request) {
	var jsonData transcriptionJobData

	// unmarshal from the response body directly into our struct
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Accepted!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}
