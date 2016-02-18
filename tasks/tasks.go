package tasks

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// Status is the status of the task.
type Status int

const (
	INPROGRESS Status = iota
	SUCCESS
	FAILURE
)

var statuses = struct {
	sync.RWMutex
	m map[string]Status
}{m: make(map[string]Status)}

// QueueTask initializes a new task.
func QueueTask() string { // TODO: add needed arguments eventually
	id := generateID(20)
	go completeTask(id)
	return id
}

// TaskStatus gets the current status of a task.
func TaskStatus(id string) (Status, error) { // TODO: add needed arguments eventually
	statuses.RLock()
	defer statuses.RUnlock()

	if status, ok := statuses.m[id]; ok {
		return status, nil
	}
	return -1, errors.New("Invalid id")

}

func completeTask(id string) {
	defer func() {
		if r := recover(); r != nil {
			statuses.Lock()
			statuses.m[id] = FAILURE
			statuses.Unlock()
		}
	}()
	statuses.Lock()
	statuses.m[id] = INPROGRESS
	statuses.Unlock()

	// DO THINGS

	statuses.Lock()
	statuses.m[id] = SUCCESS
	statuses.Unlock()
}

// Borrowed from https://siongui.github.io/2015/04/13/go-generate-random-string/
func generateID(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
