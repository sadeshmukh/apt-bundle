# Python Runtime with System Dependencies

Python application setup requiring system libraries for image processing, database connectivity, and cryptography.

## What's Included

- Python runtime: `python3`, `python3-pip`, `python3-venv`
- Image processing: `libjpeg-dev`, `libpng-dev`, `libtiff-dev`
- Database clients: `libpq-dev` (PostgreSQL), `libmysqlclient-dev` (MySQL)
- Cryptography: `libssl-dev`, `libffi-dev`

## Usage

```bash
# Create a requirements.txt file first
echo "flask==2.3.0" > requirements.txt

# Build the image
make build

# Run the application
make run
```

## Note

This example expects a `requirements.txt` file. Create one or modify the Dockerfile to skip pip installation if not needed.

## Comparison

See `Dockerfile.traditional` for the equivalent Dockerfile without apt-bundle. The traditional approach requires:
- Long RUN command with all system libraries listed inline
- Difficult to understand which libraries are for which purpose (image processing, databases, SSL)
- Hard to add or remove dependencies
- No easy way to share system dependencies across Python projects

Using apt-bundle separates system dependencies (in `Aptfile`) from Python dependencies (in `requirements.txt`), making it clear what each provides and allowing easy sharing of system dependency lists.

