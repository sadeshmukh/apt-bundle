# Dockerfile Examples with Aptfile

This directory contains practical examples demonstrating how to use `apt-bundle` with Dockerfiles in various scenarios.

## Quick Example

Here's a basic example of using `apt-bundle` in a Dockerfile:

```dockerfile
FROM ubuntu:22.04

# Prevent interactive prompts during apt operations
ENV DEBIAN_FRONTEND=noninteractive

# Install prerequisites for apt-bundle installation
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install apt-bundle
RUN ARCH=$(dpkg --print-architecture) && \
    VERSION=$(curl -s https://api.github.com/repos/apt-bundle/apt-bundle/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' || echo "") && \
    if [ -z "$VERSION" ]; then echo "Error: Failed to fetch version" && exit 1; fi && \
    VERSION_NO_V=${VERSION#v} && \
    curl -LO https://github.com/apt-bundle/apt-bundle/releases/download/${VERSION}/apt-bundle_${VERSION_NO_V}_linux_${ARCH}.deb && \
    dpkg -i apt-bundle_${VERSION_NO_V}_linux_${ARCH}.deb && \
    apt-get update && \
    apt-get install -f -y && \
    rm apt-bundle_${VERSION_NO_V}_linux_${ARCH}.deb && \
    apt-bundle --version

# Copy Aptfile and install dependencies
COPY Aptfile /app/Aptfile
WORKDIR /app
RUN apt-bundle install --file /app/Aptfile

# Your application code
COPY . /app
CMD ["/bin/bash"]
```

Each example includes:
- **Dockerfile** - Demonstrates apt-bundle usage
- **Aptfile** - Declares system dependencies
- **Makefile** - Convenient build and run commands
- **README.md** - Example-specific documentation

## Examples

### 1. Basic Development Environment (`1-basic-dev/`)
A simple web application setup with common development tools like git, curl, vim, and network utilities.

**Quick Start:**
```bash
cd 1-basic-dev
make build
make run
```

### 2. Multi-Stage Build (`2-multi-stage/`)
Demonstrates separating build dependencies from runtime dependencies to create smaller production images.

**Quick Start:**
```bash
cd 2-multi-stage
make build
make run
```

### 3. Python Runtime (`3-python-runtime/`)
Python application requiring system libraries for image processing, database connectivity, and SSL.

**Quick Start:**
```bash
cd 3-python-runtime
make build  # Creates requirements.txt if missing
make run
```

### 4. CI/CD Build Image (`4-ci-cd/`)
Docker image for CI/CD pipelines with build tools, testing frameworks, Docker CLI, and cloud tools.

**Quick Start:**
```bash
cd 4-ci-cd
make build
make run  # Includes Docker socket mount
```

### 5. Database Clients (`5-database-clients/`)
Image with database clients (PostgreSQL, MySQL, Redis, MongoDB) and monitoring/debugging tools.

**Quick Start:**
```bash
cd 5-database-clients
make build
make run
# Or connect to databases: make psql DB_HOST=host DB_NAME=db
```

## Common Patterns

All examples follow these patterns:

1. **Install apt-bundle**: Download and install the `.deb` package
2. **Copy Aptfile**: Place the Aptfile in the container
3. **Run apt-bundle**: Install dependencies from the Aptfile
4. **Clean up**: Remove apt cache to reduce image size (where applicable)

## Notes

- All examples use `ubuntu:22.04` as the base image, but can be adapted for other Debian-based distributions
- The `DEBIAN_FRONTEND=noninteractive` environment variable prevents interactive prompts during package installation
- The apt-bundle installation uses the latest release from GitHub
- Each example can be customized to fit your specific needs

