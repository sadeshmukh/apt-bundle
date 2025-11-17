#!/bin/bash
set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO="apt-bundle/apt-bundle"
GITHUB_API="https://api.github.com/repos/${REPO}"
GITHUB_RELEASES="https://github.com/${REPO}/releases"

# Error handling
error() {
    echo -e "${RED}Error:${NC} $1" >&2
    exit 1
}

info() {
    echo -e "${GREEN}Info:${NC} $1"
}

warn() {
    echo -e "${YELLOW}Warning:${NC} $1"
}

# Check if running on Debian-based system
check_system() {
    if [ ! -f /etc/debian_version ]; then
        error "This installer is for Debian-based systems only"
    fi
    
    if ! command -v dpkg >/dev/null 2>&1; then
        error "dpkg is required but not installed"
    fi
}

# Detect system architecture
detect_arch() {
    local arch=$(dpkg --print-architecture)
    
    case "$arch" in
        amd64)
            echo "amd64"
            ;;
        arm64)
            echo "arm64"
            ;;
        armhf)
            echo "armhf"
            ;;
        i386)
            echo "i386"
            ;;
        *)
            error "Unsupported architecture: $arch"
            ;;
    esac
}

# Check if running as root
check_root() {
    if [ "$EUID" -ne 0 ]; then
        error "This script must be run as root (use sudo)"
    fi
}

# Fetch latest release info
get_latest_release() {
    local token="${GITHUB_TOKEN:-}"
    local auth_header=""
    
    if [ -n "$token" ]; then
        auth_header="-H \"Authorization: token ${token}\""
    fi
    
    local latest_tag=$(curl -s ${auth_header} "${GITHUB_API}/releases/latest" | \
        grep '"tag_name":' | \
        sed -E 's/.*"([^"]+)".*/\1/' || true)
    
    if [ -z "$latest_tag" ]; then
        error "Failed to fetch latest release. Check your internet connection and GitHub access."
    fi
    
    echo "$latest_tag"
}

# Download and install .deb package
install_package() {
    local tag=$1
    local arch=$2
    local package_name="apt-bundle_${tag#v}_linux_${arch}.deb"
    local download_url="${GITHUB_RELEASES}/download/${tag}/${package_name}"
    local temp_file=$(mktemp)
    
    info "Downloading ${package_name}..."
    
    if ! curl -fsSL -o "$temp_file" "$download_url"; then
        rm -f "$temp_file"
        error "Failed to download package. URL: ${download_url}"
    fi
    
    info "Installing package..."
    if ! dpkg -i "$temp_file"; then
        warn "dpkg reported errors. Attempting to fix dependencies..."
        apt-get update -qq || true
        apt-get install -f -y || error "Failed to install dependencies"
        dpkg -i "$temp_file" || error "Failed to install package"
    fi
    
    rm -f "$temp_file"
    info "apt-bundle installed successfully!"
}

# Verify installation
verify_installation() {
    if command -v apt-bundle >/dev/null 2>&1; then
        local version=$(apt-bundle --version 2>/dev/null || echo "unknown")
        info "Installation verified. Version: ${version}"
    else
        warn "apt-bundle command not found in PATH. You may need to log out and back in."
    fi
}

# Main execution
main() {
    info "apt-bundle installer"
    info "Repository: ${REPO}"
    
    check_system
    check_root
    
    local arch=$(detect_arch)
    info "Detected architecture: ${arch}"
    
    info "Fetching latest release..."
    local latest_tag=$(get_latest_release)
    info "Latest release: ${latest_tag}"
    
    install_package "$latest_tag" "$arch"
    verify_installation
    
    echo ""
    info "Installation complete! Run 'apt-bundle --help' to get started."
}

# Run main function
main "$@"

