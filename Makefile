# Makefile for the prompt-builder component providing build, lint, and test workflows.
GO_PACKAGES := ./...
SERVICE_NAME := prompt-builder
SERVICE_ENTRYPOINT := ./cmd/prompt-builder
BINARY_DIRECTORY := $(HOME)/bin
BINARY_PATH := $(BINARY_DIRECTORY)/$(SERVICE_NAME)

.PHONY: build test test-cover test-race clean fmt vet lint run install help

build:
	mkdir -p $(BINARY_DIRECTORY)
	go build -o $(BINARY_PATH) $(SERVICE_ENTRYPOINT)

test:
	go test -v $(GO_PACKAGES)

test-cover:
	go test -coverprofile=coverage.out $(GO_PACKAGES)
	go tool cover -html=coverage.out

test-race:
	go test -race $(GO_PACKAGES)

clean:
	rm -f $(BINARY_PATH)
	rm -f coverage.out

fmt:
	gofmt -s -w .

vet:
	go vet $(GO_PACKAGES)

lint:
	golangci-lint run --fix ./...
	golangci-lint cache clean
	go clean -r -cache

run:
	go run $(SERVICE_ENTRYPOINT)

install:
	go mod tidy

help:
	@echo "Available targets:"
	@echo "  build        - Build the application binary into $(BINARY_DIRECTORY)"
	@echo "  test         - Run unit tests"
	@echo "  test-cover   - Run unit tests with coverage"
	@echo "  test-race    - Run unit tests with the race detector"
	@echo "  clean        - Remove generated binaries and coverage artifacts"
	@echo "  fmt          - Format Go source files"
	@echo "  vet          - Run go vet on the module"
	@echo "  lint         - Run golangci-lint and clean caches"
	@echo "  run          - Run the application"
	@echo "  install      - Synchronize module dependencies"
	@echo "  help         - Show this help message"
