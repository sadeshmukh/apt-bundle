.PHONY: build install clean test fmt vet lint help

BINARY_NAME=apt-bundle
BUILD_DIR=build
GO=go
GOFLAGS=-ldflags="-s -w"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/apt-bundle
	@echo "✓ Binary built at $(BUILD_DIR)/$(BINARY_NAME)"

# Install the binary to /usr/local/bin (requires sudo)
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) installed successfully"

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✓ $(BINARY_NAME) uninstalled"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

# Run golangci-lint (if installed)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  install    - Install the binary to /usr/local/bin (requires sudo)"
	@echo "  uninstall  - Remove the binary from /usr/local/bin"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run golangci-lint"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  help       - Show this help message"

