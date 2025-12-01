# Global ARG available to all stages
ARG APT_BUNDLE_REPO_ENTRY="deb [arch=amd64,arm64,armhf,i386 trusted=yes] https://apt-bundle.org/repo/ stable main"

# Build stage
FROM ubuntu:22.04 AS builder

ARG APT_BUNDLE_REPO_ENTRY

ENV DEBIAN_FRONTEND=noninteractive

# Install prerequisites and apt-bundle from APT repository
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    ca-certificates && \
    echo "$APT_BUNDLE_REPO_ENTRY" | tee /etc/apt/sources.list.d/apt-bundle.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends apt-bundle && \
    rm -rf /var/lib/apt/lists/*

# Install build dependencies
COPY Aptfile.build /app/Aptfile.build
WORKDIR /app
RUN apt-bundle install --file /app/Aptfile.build

# Build your application
COPY . /app
RUN make build-app

# Runtime stage
FROM ubuntu:22.04

ARG APT_BUNDLE_REPO_ENTRY

ENV DEBIAN_FRONTEND=noninteractive

# Install prerequisites and apt-bundle from APT repository
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    curl \
    ca-certificates && \
    echo "$APT_BUNDLE_REPO_ENTRY" | tee /etc/apt/sources.list.d/apt-bundle.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends apt-bundle && \
    rm -rf /var/lib/apt/lists/*

# Install runtime dependencies only
COPY Aptfile.runtime /app/Aptfile.runtime
WORKDIR /app
RUN apt-bundle install --file /app/Aptfile.runtime

# Copy built binary from builder
COPY --from=builder /app/bin/myapp /app/myapp

CMD ["/app/myapp"]

