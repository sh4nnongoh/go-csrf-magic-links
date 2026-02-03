.PHONY: build test run lint

build:
	go get && go mod tidy && go tool templ generate && go build

test:
	go test ./...

run:
	./go-csrf-magic-links

lint:
	golangci-lint run

clean:
	go clean
