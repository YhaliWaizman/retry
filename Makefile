ARTIFACT_NAME := retry
BINARY_PATH := bin/$(ARTIFACT_NAME)
GO_FILES := $(shell find . -name '*.go' -type f)

.PHONY: build run test clean install lint help fmt vet check test-unit test-integration

build: $(BINARY_PATH)

$(BINARY_PATH): $(GO_FILES)
	@echo "Building $(ARTIFACT_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) cmd/$(ARTIFACT_NAME)/main.go
	@echo "✓ Build complete: $(BINARY_PATH)"


run:
	@go run cmd/$(ARTIFACT_NAME)/main.go 3 echo "Hello, Retry!"


test:
	@echo "Running unit tests..."
	@go test -v ./...
	@echo ""
	@echo "Running integration tests..."
	@./tests/integration-tests.sh


test-unit:
	@go test -v ./...


test-integration:
	@./tests/integration-tests.sh


clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@echo "✓ Clean complete"


install: build
	@echo "Installing to $(GOPATH)/bin..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/
	@echo "✓ Installed to $(GOPATH)/bin/$(ARTIFACT_NAME)"


lint:
	@if command -v golangci-lint > /dev/null; then \
		echo "Running linter..."; \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi


fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "✓ Code formatted"


vet:
	@echo "Running go vet..."
	@go vet ./...
	@echo "✓ Vet complete"


check: fmt vet test
	@echo "✓ All checks passed"

.DEFAULT_GOAL := build
 
