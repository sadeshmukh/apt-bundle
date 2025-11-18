# Basic Development Environment

A simple Docker setup demonstrating how to use `apt-bundle` to install common development tools.

## What's Included

- Build tools: `build-essential`
- Version control: `git`
- Text editors: `vim`
- Utilities: `curl`, `wget`, `jq`, `htop`
- Network tools: `net-tools`

## Usage

```bash
# Build the image
make build

# Run interactively
make run

# Or use Docker directly
docker build -t basic-dev .
docker run -it basic-dev
```

## Comparison

See `Dockerfile.traditional` for the equivalent Dockerfile without apt-bundle. The traditional approach requires:
- Long, hard-to-maintain RUN commands with all packages listed inline
- Difficult to add or remove packages (requires editing Dockerfile)
- No easy way to share dependency lists across projects
- Harder to version control dependency changes

Using apt-bundle provides a cleaner, more maintainable solution where dependencies are declared in `Aptfile` and can be easily shared and version controlled.

