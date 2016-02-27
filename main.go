package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // import for side effects
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello/{name}", helloHandler)
	r.HandleFunc("/health", healthHandler)

	// add middleware
	stderrLoggingHandler := func(http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stderr, r)
	}
	middlewareRouter := alice.New(handlers.CompressHandler, stderrLoggingHandler).Then(r)

	// serve http
	http.Handle("/", middlewareRouter)
	http.ListenAndServe(":8080", nil)

	http.HandleFunc("/health", healthHandler)
	http.ListenAndServe(":8080", nil)

	http.HandleFunc("/job_status", jobStatusHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	fmt.Fprintf(w, "Hello %s!", args["name"])
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy!"))
}

// enum - job status
type Status int

const (
	INPROGRESS Status = 1 + iota
	DONE
	ERROR
)

var statuses = [...]string{
	"In progress",
	"Done",
	"Error",
}

func (s Status) String() string {
	return statuses[s-1]
}

func jobStatusHandler(s Status) {
	switch s {
	case INPROGRESS:
		w.Write("Job is in progress.")

	case DONE:
		w.Write("Job is done.")

	case ERROR:
		w.Write("Warning: error!")
	}
}
