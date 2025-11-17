.PHONY: build install clean test test-coverage test-coverage-html fmt vet lint help

BINARY_NAME=apt-bundle
BUILD_DIR=build
GO=go
GOFLAGS=-ldflags="-s -w"
INSTALL_DIR ?= /usr/local/bin
USE_SUDO ?= sudo

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/apt-bundle
	@echo "✓ Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

# Install the binary to $(INSTALL_DIR) (may require sudo)
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	@$(USE_SUDO) cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	@$(USE_SUDO) chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) installed successfully to $(INSTALL_DIR)"

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME) from $(INSTALL_DIR)..."
	@$(USE_SUDO) rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) uninstalled from $(INSTALL_DIR)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@echo "Coverage report:"
	$(GO) tool cover -func=coverage.out

# Run tests with coverage and generate HTML report
test-coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated at coverage.html"

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run golangci-lint (managed as a Go module dependency)
lint:
	@echo "Running linter..."
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build               - Build the binary"
	@echo "  install             - Install the binary to $(INSTALL_DIR) (may require sudo)"
	@echo "  uninstall           - Remove the binary from $(INSTALL_DIR)"
	@echo "  clean               - Remove build artifacts and coverage reports"
	@echo "  test                - Run tests"
	@echo "  test-coverage       - Run tests with coverage report"
	@echo "  test-coverage-html  - Run tests with HTML coverage report"
	@echo "  fmt                 - Format code"
	@echo "  vet                 - Run go vet"
	@echo "  lint                - Run golangci-lint"
	@echo "  deps                - Download and tidy dependencies"
	@echo "  help                - Show this help message"
	@echo ""
	@echo "Environment variables:"
	@echo "  INSTALL_DIR - Installation directory (default: /usr/local/bin)"
	@echo "  USE_SUDO    - Command prefix for install/uninstall (default: sudo, set to empty for no sudo)"

