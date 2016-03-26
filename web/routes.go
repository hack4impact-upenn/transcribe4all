package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
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
		"hello",
		"GET",
		"/hello/{name}",
		helloHandler,
	},
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
	route{
		"form",
		"GET",
		"/form",
		formHandler,
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

	log.Print(w, "Accepted task!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

// jobStatusHandler returns the status of a task with given id.
func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id := args["id"]

	status := tasks.DefaultTaskExecuter.GetTaskStatus(id)
	w.Write([]byte(status.String()))
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	const tpl = `
	<!DOCTYPE html>
	<html>
		<head>
		  <title></title>
		</head>
	<body>
	  <form action="/add_job" method="POST" enctype="application/json">
	    <div>URL:<input type="url" name="AudioURL" required></div>
			<div>Email Addresses:<input type="email" name="EmailAddresses" multiple required></div>
			<div><input type="submit" value="Submit"></div>
	  </form>
	</form>
	</body>
	</html>
	`
	t, _ := template.New("webpage").Parse(tpl)
	_ = t.Execute(w, transcriptionJobData{})
}
