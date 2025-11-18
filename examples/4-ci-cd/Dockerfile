FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Install prerequisites and apt-bundle
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        curl \
        ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    curl -fsSL https://raw.githubusercontent.com/apt-bundle/apt-bundle/main/install.sh | bash

# Install all CI/CD dependencies
COPY Aptfile /app/Aptfile
WORKDIR /app
RUN apt-bundle install --file /app/Aptfile

# Clean up apt cache to reduce image size
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /workspace

# Default command for CI
CMD ["/bin/bash"]

