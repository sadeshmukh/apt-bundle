# Multi-Stage Build

Demonstrates separating build dependencies from runtime dependencies to create smaller, production-ready images.

## What's Included

**Build Stage (`Aptfile.build`):**
- `build-essential` - Compiler toolchain
- `git`, `curl` - Source fetching tools
- `pkg-config`, `libssl-dev` - Development libraries

**Runtime Stage (`Aptfile.runtime`):**
- `ca-certificates` - SSL certificates
- `libssl3` - Runtime SSL library

## Usage

```bash
# Build the image
make build

# Run the application
make run

# Inspect image size
make size
```

## Benefits

- Smaller final image (build tools excluded)
- Better security (fewer attack surfaces)
- Faster deployments

## Comparison

See `Dockerfile.traditional` for the equivalent Dockerfile without apt-bundle. The traditional approach requires:
- Duplicating package lists in both build and runtime stages
- Manual management of which packages belong in which stage
- Harder to maintain consistency between stages
- Difficult to share dependency lists

Using apt-bundle with separate `Aptfile.build` and `Aptfile.runtime` files makes it clear which dependencies belong to which stage and allows easy sharing and version control.

