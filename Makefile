NON_VENDOR := $(shell go list ./... | grep -v /vendor/)

all: test

list:
	go list ./...

test:
	go test -v $(NON_VENDOR)
