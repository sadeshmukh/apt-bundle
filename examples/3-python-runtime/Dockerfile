FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive
ENV PYTHONUNBUFFERED=1

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

# Install system dependencies via Aptfile
COPY Aptfile /app/Aptfile
WORKDIR /app
RUN apt-bundle install --file /app/Aptfile

# Install Python dependencies
COPY requirements.txt /app/requirements.txt
RUN pip3 install --no-cache-dir -r requirements.txt

# Copy application code
COPY . /app

CMD ["python3", "app.py"]

