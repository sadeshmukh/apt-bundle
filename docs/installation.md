---
layout: default
title: Installation
nav_order: 2
---

# Installation

## Quick Install (Recommended)

Install the latest release using the install script:

```bash
curl -fsSL https://raw.githubusercontent.com/apt-bundle/apt-bundle/main/install.sh | sudo bash
```

This script automatically detects your system architecture and installs the appropriate `.deb` package.

## Install via APT Repository

For production environments or when you need version control, install from the official APT repository:

```bash
# Add the repository
echo "deb [arch=amd64,arm64,armhf,i386] [trusted=yes] https://apt-bundle.org/repo/ stable main" | sudo tee /etc/apt/sources.list.d/apt-bundle.list

# Update package lists
sudo apt-get update

# Install apt-bundle
sudo apt-get install -y apt-bundle
```

### Benefits

- ✅ Version control - Can pin to specific versions
- ✅ Better caching - APT caches packages efficiently
- ✅ Multi-architecture - Supports amd64, arm64, armhf, i386
- ✅ Production-ready - Standard APT package management
- ✅ Reproducible - Same version every time

## From Source

If you want to build from source or contribute to the project:

### Prerequisites

- Go 1.21 or later
- Debian/Ubuntu-based system (for running the tool)
- `make` (usually pre-installed)

### Build and Install

```bash
# Clone the repository
git clone https://github.com/apt-bundle/apt-bundle.git
cd apt-bundle

# Build the binary
make build

# Install to /usr/local/bin (requires sudo)
sudo make install
```

The binary will be installed to `/usr/local/bin/apt-bundle`.

### Build Only

If you want to build without installing:

```bash
make build
```

The binary will be in `build/apt-bundle`. You can test it:

```bash
./build/apt-bundle --help
```

### Custom Installation Directory

You can specify a custom installation directory:

```bash
INSTALL_DIR=/opt/bin sudo make install
```

### Without sudo

If you don't want to use sudo (e.g., installing to a user directory):

```bash
INSTALL_DIR=$HOME/.local/bin USE_SUDO="" make install
```

## Uninstallation

To remove the installed binary:

```bash
sudo make uninstall
```

Or manually:

```bash
sudo rm /usr/local/bin/apt-bundle
```

## Verification

After installation, verify that `apt-bundle` is available:

```bash
apt-bundle --version
apt-bundle --help
```

## Next Steps

- Learn how to use apt-bundle in the [Usage Guide](usage.html)
- Understand the [Aptfile Format](aptfile-format.html)
- Check out [Use Cases](index.html#use-cases) for common scenarios

