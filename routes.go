package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"github.com/gorilla/mux"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type transcriptionJobData struct {
	AudioURL string `json:"audioURL"`
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
		"DELETE",
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
		"email",
		"GET",
		"/email",
		emailHandler,
	},
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}

// initiateTranscriptionJobHandle takes a POST request containing a json object,
// decodes it into an audioData struct, and returns the struct json-encoded.
func initiateTranscriptionJobHandler(w http.ResponseWriter, r *http.Request) {
	var jsonData transcriptionJobData

	// unmarshal from the response body directly into our struct
	if err := json.NewDecoder(r.Body).Decode(&jsonData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return the struct encoded as json
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(jsonData)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy!"))
}

func emailHandler(w http.ResponseWriter, r *http.Request) {
	// Set up authentication information.
	password := os.Getenv("MAIL_PASSWORD")

	auth := smtp.PlainAuth(
		"",
		"test4impact@gmail.com",
		password,
		"smtp.gmail.com",
	)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(
		"smtp.gmail.com:25",
		auth,
		"test4impact@gmail.org",
		[]string{"yoninachmany@gmail.com"},
		[]byte("This is the email body."),
	)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte("email!"))
}
