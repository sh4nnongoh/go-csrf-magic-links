.PHONY: build test run lint

build:
	go build ./...

test:
	go test ./...

run:
	./m

lint:
	golangci-lint run
