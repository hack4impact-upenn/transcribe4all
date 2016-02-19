package tasks

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

// Status is the status of the task.
type Status int

// TaskExecutor executes a series of task functions.
type TaskExecutor interface {
	QueueTask(task func() error) string
	GetTaskStatus(id string) (Status, error)
	completeTask(id string, task func() error)
}

type defaultExecuter struct {
	sync.RWMutex
	m map[string]Status
}

// These are some enumerated Status constants.
// INPROGRESS: Task is still in progress.
// SUCCESS: Task finished successfully.
// FAILURE: Task finished unsuccessfully.
// NOTFOUND: Task could not be found.
const (
	INPROGRESS Status = iota
	SUCCESS
	FAILURE
	NOTFOUND
)

// NewTaskExectuer returns a TaskExecutor ready to execute.
func NewTaskExectuer() TaskExecutor {
	return &defaultExecuter{m: make(map[string]Status)}
}

// QueueTask initializes a new task. It takes a generic task function.
func (ex *defaultExecuter) QueueTask(task func() error) string {
	id := generateID(20)
	ex.Lock()
	ex.m[id] = INPROGRESS
	ex.Unlock()
	go ex.completeTask(id, task)
	return id
}

// GetTaskStatus gets the current status of a task.
func (ex *defaultExecuter) GetTaskStatus(id string) (Status, error) {
	ex.RLock()
	defer ex.RUnlock()

	if status, ok := ex.m[id]; ok {
		return status, nil
	}
	return NOTFOUND, errors.New("Invalid id")
}

func (ex *defaultExecuter) completeTask(id string, task func() error) {
	defer func() {
		if r := recover(); r != nil {
			ex.Lock()
			ex.m[id] = FAILURE
			ex.Unlock()
		}
	}()

	// Run the task.
	if err := task(); err != nil {
		ex.Lock()
		ex.m[id] = FAILURE
		ex.Unlock()
	}

	ex.Lock()
	ex.m[id] = SUCCESS
	ex.Unlock()
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
