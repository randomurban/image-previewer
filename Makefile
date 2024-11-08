include .env

BIN=./bin/previewer

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/previewer

run: build
	$(BIN) -config ./.env

lint:
	golangci-lint run ./...

test:
	go test -race -count 100 ./internal/...

.PHONY: build run lint test