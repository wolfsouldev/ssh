.PHONY: build clean test lint fmt vet install uninstall help

BINARY_NAME := sshh
VERSION     ?= 1.0.0
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE        := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
MODULE      := github.com/wolfsouldev/ssh/internal/version

LDFLAGS := -ldflags "\
  -s -w \
  -X $(MODULE).Version=$(VERSION) \
  -X $(MODULE).Commit=$(COMMIT) \
  -X $(MODULE).Date=$(DATE)"

# Detect OS for binary extension
ifeq ($(OS),Windows_NT)
    BINARY := $(BINARY_NAME).exe
    RM_CMD := del /f $(BINARY) 2>nul
else
    BINARY := $(BINARY_NAME)
    RM_CMD := rm -f $(BINARY)
endif

## help: Show this help message
help:
	@echo ""
	@echo "  SSHH v$(VERSION) - Build Commands"
	@echo "  ================================="
	@echo ""
	@echo "  make build      Build the binary"
	@echo "  make clean      Remove build artifacts"
	@echo "  make test       Run tests"
	@echo "  make lint       Run linters"
	@echo "  make fmt        Format code"
	@echo "  make vet        Run go vet"
	@echo "  make install    Build and install to /usr/local/bin (Linux/macOS)"
	@echo "  make uninstall  Remove from /usr/local/bin"
	@echo "  make all        Format, lint, test, and build"
	@echo ""

## build: Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY) .

## clean: Remove build artifacts
clean:
	go clean
	$(RM_CMD)

## test: Run all tests
test:
	go test -v -race -coverprofile=coverage.out ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format Go code
fmt:
	gofmt -s -w .

## vet: Run go vet
vet:
	go vet ./...

## install: Build and install to system
install: build
ifeq ($(OS),Windows_NT)
	.\$(BINARY) install
else
	sudo cp $(BINARY) /usr/local/bin/$(BINARY_NAME)
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"
endif

## uninstall: Remove from system
uninstall:
ifeq ($(OS),Windows_NT)
	.\$(BINARY) uninstall
else
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Removed /usr/local/bin/$(BINARY_NAME)"
endif

## all: Format, lint, test, and build
all: fmt vet lint test build
