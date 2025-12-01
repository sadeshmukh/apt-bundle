# Installing apt-bundle via APT Repository

A production-ready example demonstrating how to install `apt-bundle` from the official APT repository.

## Installation Method

This example uses the official APT repository, which provides better control and version management:

```dockerfile
RUN echo "deb [arch=amd64,arm64,armhf,i386 trusted=yes] https://apt-bundle.org/repo/ stable main" | tee /etc/apt/sources.list.d/apt-bundle.list && \
    apt-get update && \
    apt-get install -y apt-bundle
```

## What's Included

- Utilities: `curl`, `git`

## Usage

```bash
# Build the image
make build

# Run interactively
make run

# Test the installed packages
make test
```

## When to Use This Method

- **Production environments** - More reliable and version-controlled
- **Docker images** - Better layer caching and reproducibility
- **CI/CD pipelines** - Consistent versions across builds
- **When you need specific versions** - Can pin to specific apt-bundle versions
- **Multi-architecture support** - Supports amd64, arm64, armhf, i386

## Benefits Over install.sh

- ✅ Version control - Can pin to specific versions
- ✅ Better caching - APT caches packages efficiently
- ✅ Multi-architecture - Supports all architectures
- ✅ Production-ready - Standard APT package management
- ✅ Reproducible - Same version every time

## Note

The `trusted=yes` flag is used because the repository is currently unsigned. For production use, consider setting up GPG signing for the repository.
