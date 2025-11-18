FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive

# Install prerequisites and apt-bundle
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl ca-certificates && \
    rm -rf /var/lib/apt/lists/* && \
    curl -fsSL https://raw.githubusercontent.com/apt-bundle/apt-bundle/main/install.sh | bash

# Install database clients and tools
COPY Aptfile /app/Aptfile
WORKDIR /app
RUN apt-bundle install --file /app/Aptfile

# Clean up
RUN apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy utility scripts
COPY scripts/ /app/scripts/
RUN chmod +x /app/scripts/*.sh

WORKDIR /app

CMD ["/bin/bash"]

