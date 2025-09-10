# Prompt Builder Makefile
# A standalone Go application for building structured prompts

# Variables
BINARY_NAME=prompt-builder
MAIN_PATH=./cmd/prompt-builder
PACKAGE_PATH=./internal/prompt-builder
MODULE_NAME=github.com/nnikolov3/prompt-builder
BIN_DIR=$(HOME)/bin

# Go build flags
LDFLAGS=-ldflags "-X main.Version=$(shell git describe --tags --always --dirty) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# Default target
.PHONY: all
all: clean build test lint

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for multiple platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BIN_DIR)/$(BINARY_NAME) -h

# Test targets
.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

.PHONY: test-race
test-race:
	@echo "Running tests with race detection..."
	go test -race -v ./...

.PHONY: test-benchmark
test-benchmark:
	@echo "Running benchmarks..."
	go test -bench=. -benchmem ./...

# Linting targets
.PHONY: lint
lint: lint-golangci-lint lint-gofmt

.PHONY: lint-golangci-lint
lint-golangci-lint:
	@echo "Running golangci-lint..."
	golangci-lint run --fix

.PHONY: lint-gofmt
lint-gofmt:
	@echo "Checking code formatting..."
	gofmt -d .

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	gofmt -w .

# Code quality targets
.PHONY: tidy
tidy:
	@echo "Tidying go.mod..."
	go mod tidy

.PHONY: verify
verify:
	@echo "Verifying dependencies..."
	go mod verify

.PHONY: download
download:
	@echo "Downloading dependencies..."
	go mod download

# Clean targets
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BIN_DIR)/$(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out
	rm -f coverage.html
	go clean -cache -testcache

.PHONY: clean-all
clean-all: clean
	@echo "Cleaning all artifacts..."
	rm -rf dist/
	rm -rf vendor/

# Development targets
.PHONY: dev
dev:
	@echo "Starting development mode..."
	@echo "Watching for changes and rebuilding..."
	@while true; do \
		inotifywait -r -e modify,create,delete .; \
		make build test; \
	done

.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME)..."
	@echo "Binary already installed in $(BIN_DIR)/$(BINARY_NAME)"

.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(BIN_DIR)/$(BINARY_NAME)

# Documentation targets
.PHONY: docs
docs:
	@echo "Generating documentation..."
	godoc -http=:6060

.PHONY: readme
readme:
	@echo "Checking README..."
	@if [ ! -f README.md ]; then \
		echo "README.md not found!"; \
		exit 1; \
	fi

# Release targets
.PHONY: release
release: clean build-all test lint
	@echo "Creating release..."
	mkdir -p dist
	cp $(BINARY_NAME)-* dist/
	cp README.md dist/
	@echo "Release artifacts created in dist/"

.PHONY: release-tag
release-tag:
	@echo "Creating git tag for release..."
	@read -p "Enter version tag (e.g., v1.0.0): " version; \
	git tag -a $$version -m "Release $$version"; \
	git push origin $$version

# Docker targets
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	docker run --rm -it $(BINARY_NAME):latest

# Security targets
.PHONY: security
security: lint-gosec
	@echo "Security check completed"

.PHONY: audit
audit:
	@echo "Auditing dependencies..."
	go list -json -deps ./... | nancy sleuth

# Performance targets
.PHONY: profile
profile:
	@echo "Running performance profiling..."
	go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...

.PHONY: profile-cpu
profile-cpu:
	@echo "Analyzing CPU profile..."
	go tool pprof cpu.prof

.PHONY: profile-mem
profile-mem:
	@echo "Analyzing memory profile..."
	go tool pprof mem.prof

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-all      - Build for Linux, macOS, and Windows"
	@echo "  run            - Build and run the application"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-benchmark - Run benchmarks"
	@echo "  lint           - Run all linters"
	@echo "  fmt            - Format code with gofmt"
	@echo "  tidy           - Tidy go.mod"
	@echo "  clean          - Clean build artifacts"
	@echo "  install        - Install to /usr/local/bin"
	@echo "  uninstall      - Remove from /usr/local/bin"
	@echo "  release        - Create release artifacts"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  security       - Run security checks"
	@echo "  profile        - Run performance profiling"
	@echo "  help           - Show this help message"

# Default target
.DEFAULT_GOAL := help
