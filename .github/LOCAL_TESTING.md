# Testing GitHub Actions Workflows Locally

There are two ways to test the CI workflow locally before pushing:

## Option 1: Quick Test with Makefile (Recommended)

The easiest way is to use the `ci-test` Makefile target, which mimics the exact CI build step:

```bash
# Test with default architecture (amd64) and version from VERSION file
make ci-test

# Test with specific architecture
make ci-test ARCH=arm64

# Test with specific architecture and version
make ci-test ARCH=amd64 VERSION=0.1.5
```

This will:
- Build the binary for the specified architecture
- Create a .deb package using nfpm (with the same environment variables as CI)
- Rename the package to include architecture in filename
- Copy it to `artifacts/` directory

**Note:** This only tests the build step. It doesn't test the full workflow (version calculation, release creation, etc.).

## Option 2: Full Workflow Testing with `act`

For testing the complete GitHub Actions workflow, you can use [act](https://github.com/nektos/act):

### Installation

```bash
# On macOS
brew install act

# On Linux (using the install script)
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash

# Or download from releases: https://github.com/nektos/act/releases
```

### Basic Usage

```bash
# List all workflows
act -l

# Run the build job (dry-run, won't actually build)
act -j build --dryrun

# Run the build job for a specific architecture
act -j build --matrix goarch:amd64 --matrix debarch:amd64

# Run with environment variables
act -j build -e VERSION=0.1.0
```

### Limitations

- `act` runs workflows in Docker containers, so it's slower than the Makefile approach
- Some GitHub Actions features may not work perfectly (secrets, API calls, etc.)
- The release job requires GitHub API access which may not work locally

### Recommended Workflow

1. **Quick iteration**: Use `make ci-test` to test build changes rapidly
2. **Full workflow**: Use `act` before pushing to test the complete workflow
3. **Final verification**: Push to a test branch and let GitHub Actions run

## Troubleshooting

### nfpm not found
The Makefile will auto-install nfpm if it's missing. For `act`, ensure Docker is running.

### Architecture-specific issues
Test each architecture individually:
```bash
make ci-test ARCH=amd64
make ci-test ARCH=arm64
make ci-test ARCH=armhf
make ci-test ARCH=i386
```

### Version format
The CI workflow expects versions in `MAJOR.MINOR.PATCH` format. The Makefile defaults to `VERSION.0` if only `MAJOR.MINOR` is in the VERSION file.

