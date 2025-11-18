# apt-bundle

A declarative package manager for apt, inspired by Homebrew's `brew bundle`.

**[📚 Full Documentation](https://apt-bundle.github.io/apt-bundle/)** | [Installation](#installation) | [Usage](#usage)

## Overview

`apt-bundle` provides a simple, declarative, and shareable way to manage apt packages and repositories on Debian-based systems. Define your system dependencies in an `Aptfile` and install them with a single command.

## Features

- 📦 **Declarative Package Management**: Define packages in a simple text file
- 🔄 **Idempotent Operations**: Safe to run multiple times
- 🔑 **Repository & Key Management**: Add PPAs, custom repositories, and GPG keys
- 📝 **Version Pinning**: Install specific package versions
- 🚀 **Simple CLI**: Easy-to-use command-line interface

## Installation

### Quick Install (Recommended)

Install the latest release using the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/apt-bundle/apt-bundle/main/install.sh | sudo bash
```

### Manual Installation from .deb Package

Download and install the appropriate `.deb` package for your architecture:

```bash
# Detect your architecture
ARCH=$(dpkg --print-architecture)

# Download latest release (replace v1.0.0 with actual version)
VERSION=$(curl -s https://api.github.com/repos/apt-bundle/apt-bundle/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -LO https://github.com/apt-bundle/apt-bundle/releases/download/${VERSION}/apt-bundle_${VERSION#v}_linux_${ARCH}.deb

# Install
sudo dpkg -i apt-bundle_${VERSION#v}_linux_${ARCH}.deb
sudo apt-get install -f  # Install dependencies if needed
```

### From Source

```bash
# Clone the repository
git clone https://github.com/apt-bundle/apt-bundle.git
cd apt-bundle

# Build and install
make build
sudo make install
```

The binary will be installed to `/usr/local/bin/apt-bundle`.

### Building

```bash
# Build the binary
make build

# The binary will be in build/apt-bundle
./build/apt-bundle --help
```

## Usage

### Basic Commands

```bash
# Install packages from Aptfile (default: ./Aptfile)
sudo apt-bundle

# or explicitly
sudo apt-bundle install

# Use a different Aptfile
sudo apt-bundle --file /path/to/Aptfile

# Skip updating package lists (useful in CI/CD)
sudo apt-bundle --no-update

# Check if packages are installed (no root required)
apt-bundle check

# Generate an Aptfile from current system
apt-bundle dump > Aptfile
```

### Aptfile Format

The `Aptfile` is a simple line-oriented text file with the following directives:

#### Install Packages

```aptfile
# Install latest version
apt vim
apt curl
apt git

# Install specific version
apt "nano=2.9.3-2"
```

#### Add PPAs

```aptfile
ppa ppa:ondrej/php
apt php8.1
```

#### Add Custom Repositories

```aptfile
# Add GPG key
key https://download.docker.com/linux/ubuntu/gpg

# Add repository
deb "[arch=amd64] https://download.docker.com/linux/ubuntu focal stable"

# Install packages from the repository
apt docker-ce
apt docker-ce-cli
```

### Complete Example

```aptfile
# Core development tools
apt build-essential
apt curl
apt git
apt vim
apt htop

# Specific version
apt "nano=2.9.3-2"

# PHP from PPA
ppa ppa:ondrej/php
apt php8.1
apt php8.1-cli
apt php8.1-fpm

# Docker
key https://download.docker.com/linux/ubuntu/gpg
deb "[arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
apt docker-ce
apt docker-ce-cli
apt containerd.io
```

## Use Cases

### Developer Onboarding

```bash
# Clone project
git clone https://github.com/myorg/myproject.git
cd myproject

# Install all system dependencies
sudo apt-bundle
```

### Dockerfile

```dockerfile
FROM ubuntu:22.04

# Install apt-bundle
RUN ARCH=$(dpkg --print-architecture) && \
    VERSION=$(curl -s https://api.github.com/repos/apt-bundle/apt-bundle/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/') && \
    curl -LO https://github.com/apt-bundle/apt-bundle/releases/download/${VERSION}/apt-bundle_${VERSION#v}_linux_${ARCH}.deb && \
    dpkg -i apt-bundle_${VERSION#v}_linux_${ARCH}.deb && \
    apt-get install -f -y && \
    rm apt-bundle_${VERSION#v}_linux_${ARCH}.deb

# Copy Aptfile and install dependencies
COPY Aptfile /app/Aptfile
WORKDIR /app
RUN apt-bundle
```

### System Sync

```bash
# On primary workstation
apt-bundle dump > Aptfile

# On new laptop
sudo apt-bundle
```

## Development

### Project Structure

```
apt-bundle/
├── cmd/
│   └── apt-bundle/       # Main entry point
├── internal/
│   ├── aptfile/          # Aptfile parsing
│   ├── commands/         # CLI commands (install, dump, check)
│   ├── apt/              # APT interactions
│   └── config/           # Configuration
├── docs/                 # Requirements and specifications
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

### Building and Testing

```bash
# Format code
make fmt

# Run static analysis
make vet

# Run tests
make test

# Build
make build

# Install locally for testing
sudo make install
```

### Requirements

- Go 1.21 or later
- Debian/Ubuntu-based system (for running the tool)

### Version Management

The project uses a VERSION file for version management:
- The `VERSION` file contains the major.minor version (e.g., `1.0`)
- Patch versions are automatically incremented on each release
- To update the major or minor version, edit the `VERSION` file
- Releases are automatically created when code is merged to the `main` branch

## Technical Details

### Binary Characteristics

- **Self-contained**: The Go binary is statically linked and doesn't require external `.so` or `.dll` files
- **CGO_ENABLED=0**: Ensures pure Go compilation without C dependencies
- **Small size**: Compiled with `-ldflags="-s -w"` to strip debug symbols

### Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework

## Documentation

📚 **[Full Documentation Site](https://apt-bundle.github.io/apt-bundle/)** - Complete user guide, developer documentation, and API reference

For internal specifications:
- [Requirements](specs/requirements.md) - Detailed functional requirements
- [Technical Specification](specs/tech-specs.md) - Aptfile format and implementation details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[Add your license here]

## Roadmap

- [ ] Core functionality (install, dump, check)
- [ ] Package version handling
- [ ] PPA management
- [ ] Custom repository management
- [ ] GPG key management
- [ ] Comprehensive test suite
- [ ] CI/CD pipeline
- [ ] Release automation

