all: test

test:
	go test -v $(go list ./... | grep -v /vendor/)
