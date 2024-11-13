include .env

BIN=./bin/previewer

.PHONY: build
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/previewer

.PHONY: run
run: build
	$(BIN) -config ./.env

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test -race -count 100 ./internal/...

.PHONY: integration-test-build
integration-test-build:
	docker compose -f ./tests/compose.yaml build

.PHONY: integration-test-up
integration-test-up: integration-test-build
	docker compose -f ./tests/compose.yaml up

.PHONY: integration-test
integration-test: integration-test-build
	docker compose -f ./tests/compose.yaml up --exit-code-from integration-test

