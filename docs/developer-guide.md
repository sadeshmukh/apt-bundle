---
layout: default
title: Developer Guide
nav_order: 5
---

# Developer Guide

This guide is for developers who want to contribute to apt-bundle or understand its internals.

## Project Structure

```
apt-bundle/
├── cmd/
│   └── apt-bundle/       # Main entry point
├── internal/
│   ├── aptfile/          # Aptfile parsing
│   ├── commands/         # CLI commands (install, dump, check)
│   ├── apt/              # APT interactions
│   └── config/           # Configuration
├── docs/                 # Documentation
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

## Prerequisites

- Go 1.21 or later
- Debian/Ubuntu-based system (for testing)
- `make` (for build automation)

## Building

### Build the Binary

```bash
make build
```

The binary will be created at `build/apt-bundle`.

### Install Locally

```bash
sudo make install
```

This installs to `/usr/local/bin/apt-bundle` by default. You can customize:

```bash
INSTALL_DIR=$HOME/.local/bin USE_SUDO="" make install
```

### Build Options

The Makefile uses the following build flags:
- `CGO_ENABLED=0`: Pure Go compilation (no C dependencies)
- `-ldflags="-s -w"`: Strip debug symbols for smaller binary size

## Development Workflow

### Format Code

```bash
make fmt
```

This runs `go fmt ./...` to format all Go code.

### Run Static Analysis

```bash
make vet
```

This runs `go vet ./...` to catch common errors.

### Run Linter

```bash
make lint
```

This runs `golangci-lint` for comprehensive code analysis.

### Run Tests

```bash
make test
```

This runs all tests with verbose output.

### Clean Build Artifacts

```bash
make clean
```

Removes the `build/` directory.

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests for a specific package
go test ./internal/aptfile/...

# Run tests with coverage
go test -cover ./...
```

### Writing Tests

- Place test files alongside source files with `_test.go` suffix
- Use standard Go testing patterns
- Mock external commands (apt-get, add-apt-repository) for unit tests
- Use integration tests for end-to-end scenarios

### Test Requirements

- Tests should be idempotent when possible
- Avoid requiring root privileges for unit tests
- Use temporary directories for file operations
- Clean up after tests

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting (enforced by `make fmt`)
- Follow the guidelines enforced by `golangci-lint`
- Write clear, self-documenting code
- Add comments for exported functions and types

## Architecture

### Command Structure

The CLI uses [Cobra](https://github.com/spf13/cobra) for command parsing:

- `cmd/apt-bundle/main.go`: Entry point
- `internal/commands/`: Command implementations
  - `install.go`: Install command
  - `check.go`: Check command
  - `dump.go`: Dump command

### Aptfile Parsing

- `internal/aptfile/`: Parses Aptfile format
  - Tokenizes directives
  - Validates syntax
  - Provides structured representation

### APT Integration

- `internal/apt/`: Wraps apt-get and add-apt-repository
  - Executes system commands
  - Checks for existing packages/repositories
  - Handles errors gracefully

### Configuration

- `internal/config/`: Manages configuration
  - Default file paths
  - Command-line flags
  - Environment variables

## Contributing

### Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/yourusername/apt-bundle.git`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `make test`
6. Format code: `make fmt`
7. Run linter: `make lint`
8. Commit your changes
9. Push to your fork
10. Open a Pull Request

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Include tests for new functionality
- Update documentation as needed
- Ensure all tests pass
- Follow the existing code style

### Commit Messages

Write clear, descriptive commit messages:

```
Add support for version pinning in Aptfile

- Parse version strings from apt directives
- Pass version to apt-get install
- Add tests for version parsing
```

## Dependencies

### External Dependencies

- [spf13/cobra](https://github.com/spf13/cobra): CLI framework

### Development Dependencies

- [golangci/golangci-lint](https://github.com/golangci/golangci-lint): Linter

### Managing Dependencies

```bash
# Download dependencies
make deps

# Add a new dependency
go get <package>

# Update dependencies
go get -u ./...
go mod tidy
```

## Binary Characteristics

The apt-bundle binary is designed to be:

- **Self-contained**: Statically linked, no external `.so` or `.dll` files
- **Portable**: Works on any Linux system with the same architecture
- **Small**: Stripped debug symbols for minimal size
- **Fast**: Minimal overhead over native apt commands

## Release Process

1. Update version in code (if versioning is implemented)
2. Update CHANGELOG.md
3. Create a git tag
4. Build release binaries
5. Create GitHub release
6. Upload binaries

## Troubleshooting

### Build Issues

- Ensure Go 1.21+ is installed: `go version`
- Check `go.mod` is up to date: `go mod tidy`
- Clean and rebuild: `make clean && make build`

### Test Issues

- Some tests may require a Debian/Ubuntu system
- Integration tests may need root privileges (use sudo)
- Check that required system commands are available

### Development Environment

- Use `direnv` or similar for environment management
- Consider using Docker for isolated testing
- Use `go run` for quick iteration: `go run ./cmd/apt-bundle`

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Cobra Documentation](https://github.com/spf13/cobra/blob/main/site/content/introduction.md)
- [Go Testing](https://golang.org/pkg/testing/)
- [Project Requirements](https://github.com/apt-bundle/apt-bundle/blob/main/specs/requirements.md) - Internal specification
- [Technical Specification](https://github.com/apt-bundle/apt-bundle/blob/main/specs/tech-specs.md) - Internal specification

## Getting Help

- Open an issue on GitHub for bugs or feature requests
- Check existing issues and discussions
- Review the [API Reference](api-reference.html) for code documentation

