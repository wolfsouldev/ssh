.PHONY: build clean test lint fmt vet install uninstall help

BINARY_NAME=sshh
VERSION?=0.1.0
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

## help: Show this help message
help:
	@echo.
	@echo  SSHH - Build Commands
	@echo  =====================
	@echo.
	@echo  make build      Build the binary
	@echo  make clean      Remove build artifacts
	@echo  make test       Run tests
	@echo  make lint       Run linters
	@echo  make fmt        Format code
	@echo  make vet        Run go vet
	@echo  make install    Build and install to PATH
	@echo  make all        Format, lint, test, and build
	@echo.

## build: Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME).exe .

## clean: Remove build artifacts
clean:
	go clean
	del /f $(BINARY_NAME).exe 2>nul

## test: Run all tests
test:
	go test -v -race -coverprofile=coverage.out ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format Go code
fmt:
	gofmt -s -w .
	goimports -w .

## vet: Run go vet
vet:
	go vet ./...

## install: Build and install
install: build
	.\$(BINARY_NAME).exe install

## uninstall: Uninstall from PATH
uninstall:
	.\$(BINARY_NAME).exe uninstall

## all: Format, lint, test, and build
all: fmt vet lint test build
