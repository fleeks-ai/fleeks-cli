#!/bin/bash
# Fleeks CLI Installer - CDN Distribution
# Usage: curl -fsSL http://downloads.fleeks.ai/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CDN_BASE_URL="http://downloads.fleeks.ai"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="fleeks"

# Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Fleeks CLI Installation Script       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Determine binary name
case "$OS" in
    Linux*)
        case "$ARCH" in
            x86_64) BINARY="fleeks-linux-amd64" ;;
            aarch64|arm64) BINARY="fleeks-linux-arm64" ;;
            *) echo -e "${RED}[ERROR] Unsupported architecture: $ARCH${NC}"; exit 1 ;;
        esac
        ;;
    Darwin*)
        case "$ARCH" in
            x86_64) BINARY="fleeks-darwin-amd64" ;;
            arm64) BINARY="fleeks-darwin-arm64" ;;
            *) echo -e "${RED}[ERROR] Unsupported architecture: $ARCH${NC}"; exit 1 ;;
        esac
        ;;
    *)
        echo -e "${RED}[ERROR] Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}[INFO]${NC} Detected platform: $OS ($ARCH)"

# Get latest version from GitHub API
echo -e "${GREEN}[INFO]${NC} Fetching latest version..."
VERSION=$(curl -s https://api.github.com/repos/fleeks-ai/fleeks-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"v?([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${YELLOW}[WARN]${NC} Could not determine latest version, using 'latest' folder"
    VERSION="latest"
    DOWNLOAD_URL="${CDN_BASE_URL}/latest/${BINARY}"
else
    echo -e "${GREEN}[INFO]${NC} Latest version: v${VERSION}"
    DOWNLOAD_URL="${CDN_BASE_URL}/v${VERSION}/${BINARY}"
fi

echo ""

# Download binary
TMP_FILE="/tmp/fleeks"

echo -e "${GREEN}[INFO]${NC} Downloading from: ${DOWNLOAD_URL}"
if command -v curl &> /dev/null; then
    if ! curl -fL -o "$TMP_FILE" "$DOWNLOAD_URL"; then
        echo -e "${RED}[ERROR] Failed to download binary${NC}"
        exit 1
    fi
elif command -v wget &> /dev/null; then
    if ! wget -q -O "$TMP_FILE" "$DOWNLOAD_URL"; then
        echo -e "${RED}[ERROR] Failed to download binary${NC}"
        exit 1
    fi
else
    echo -e "${RED}[ERROR] curl or wget is required${NC}"
    exit 1
fi

echo -e "${GREEN}[SUCCESS]${NC} Downloaded successfully"

# Verify checksum (optional)
if [ "$VERSION" != "latest" ]; then
    echo -e "${GREEN}[INFO]${NC} Verifying checksum..."
    CHECKSUM_URL="${CDN_BASE_URL}/v${VERSION}/SHA256SUMS.txt"
    if curl -fsSL "$CHECKSUM_URL" 2>/dev/null | grep "$BINARY" | (cd /tmp && sha256sum -c --status 2>/dev/null); then
        echo -e "${GREEN}[SUCCESS]${NC} Checksum verified"
    else
        echo -e "${YELLOW}[WARN]${NC} Checksum verification failed or not available"
    fi
fi

# Make executable
chmod +x "$TMP_FILE"

# Install
echo -e "${GREEN}[INFO]${NC} Installing to ${INSTALL_DIR}..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$INSTALL_DIR/fleeks"
    echo -e "${GREEN}[SUCCESS]${NC} Installed to $INSTALL_DIR/fleeks"
else
    echo -e "${YELLOW}[WARN]${NC} Installing to $INSTALL_DIR requires sudo..."
    sudo mv "$TMP_FILE" "$INSTALL_DIR/fleeks"
    echo -e "${GREEN}[SUCCESS]${NC} Installed to $INSTALL_DIR/fleeks"
fi

# Verify installation
echo ""
echo -e "${GREEN}[INFO]${NC} Verifying installation..."
if command -v fleeks &> /dev/null; then
    INSTALLED_VERSION=$(fleeks --version 2>&1 | head -n1)
    echo -e "${GREEN}[SUCCESS]${NC} ${INSTALLED_VERSION}"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}Installation complete!${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "Get started with:"
    echo -e "  ${YELLOW}fleeks auth login${NC}"
    echo -e "  ${YELLOW}fleeks --help${NC}"
    echo ""
else
    echo -e "${RED}[ERROR] Installation failed${NC}"
    exit 1
fi
