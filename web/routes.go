package web

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hack4impact/transcribe4all/tasks"
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
		"add_job",
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
	route{
		"job_status",
		"GET",
		"/job_status/{id}",
		jobStatusHandler,
	},
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

	io.WriteString(w, "Accepted!")
}

// healthHandler returns a 200 response to the client if the server is healthy.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK :)")
}

// jobStatusHandler returns the status of a task with given id.
func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]

	status := tasks.DefaultTaskExecuter.GetTaskStatus(id)
	io.WriteString(w, status.String())
}
