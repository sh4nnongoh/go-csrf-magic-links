.PHONY: build test run lint

build:
	go get && go mod tidy && go build ./...

test:
	go test ./...

run:
	./m

lint:
	golangci-lint run

clean:
	go clean
