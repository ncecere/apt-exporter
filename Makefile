# APT Exporter Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=apt-exporter
BINARY_UNIX=$(BINARY_NAME)_unix

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS = -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build the project
all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/apt-exporter

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
cover:
	$(GOTEST) -v -cover ./...

# Update dependencies
deps:
	$(GOMOD) tidy

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v ./cmd/apt-exporter

# Run the application
run:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/apt-exporter
	./$(BINARY_NAME)

# Install the application
install:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v ./cmd/apt-exporter
	mv $(BINARY_NAME) $(GOPATH)/bin/

# Format code
fmt:
	gofmt -s -w .

# Lint code
lint:
	golangci-lint run ./...

.PHONY: all build clean test cover deps build-linux run install fmt lint
