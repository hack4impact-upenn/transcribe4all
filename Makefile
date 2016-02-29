NON_VENDOR := $(shell go list ./... | grep -v /vendor/)

all: test

test:
	go test -v $(NON_VENDOR)
