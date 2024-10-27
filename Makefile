include .env

BIN=./bin

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/previewer

lint:
	golangci-lint run ./...
