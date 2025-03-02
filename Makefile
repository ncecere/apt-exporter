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
	rm -f $(BINARY_NAME)_*

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
cover:
	$(GOTEST) -v -cover ./...

# Generate coverage report
cover-html:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Update dependencies
deps:
	$(GOMOD) tidy

# Build for Linux (amd64)
build-linux-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_amd64 -v ./cmd/apt-exporter

# Build for Linux (arm64)
build-linux-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_arm64 -v ./cmd/apt-exporter

# Build for Linux (arm)
build-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_arm -v ./cmd/apt-exporter

# Build for Linux (386)
build-linux-386:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME)_linux_386 -v ./cmd/apt-exporter

# Build for all platforms
build-all: build-linux-amd64 build-linux-arm64 build-linux-arm build-linux-386

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

# Create release archives
release: build-all
	mkdir -p release
	tar -czf release/$(BINARY_NAME)_$(VERSION)_linux_amd64.tar.gz $(BINARY_NAME)_linux_amd64 LICENSE README.md config.yml
	tar -czf release/$(BINARY_NAME)_$(VERSION)_linux_arm64.tar.gz $(BINARY_NAME)_linux_arm64 LICENSE README.md config.yml
	tar -czf release/$(BINARY_NAME)_$(VERSION)_linux_arm.tar.gz $(BINARY_NAME)_linux_arm LICENSE README.md config.yml
	tar -czf release/$(BINARY_NAME)_$(VERSION)_linux_386.tar.gz $(BINARY_NAME)_linux_386 LICENSE README.md config.yml

.PHONY: all build clean test cover cover-html deps build-linux-amd64 build-linux-arm64 build-linux-arm build-linux-386 build-all run install fmt lint release
