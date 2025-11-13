#!/bin/bash
# Fleeks CLI Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/fleeks-ai/fleeks-cli/main/install.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

echo -e "${GREEN}Fleeks CLI Installer${NC}"
echo ""

# Determine binary name
case "$OS" in
    Linux*)
        case "$ARCH" in
            x86_64) BINARY="fleeks-linux-amd64" ;;
            aarch64|arm64) BINARY="fleeks-linux-arm64" ;;
            *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
        esac
        ;;
    Darwin*)
        case "$ARCH" in
            x86_64) BINARY="fleeks-darwin-amd64" ;;
            arm64) BINARY="fleeks-darwin-arm64" ;;
            *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
        esac
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

# Get latest version
echo "Fetching latest version..."
VERSION=$(curl -s https://api.github.com/repos/fleeks-ai/fleeks-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$VERSION" ]; then
    echo -e "${RED}Failed to fetch latest version${NC}"
    exit 1
fi

echo -e "Latest version: ${GREEN}$VERSION${NC}"
echo ""

# Download binary
DOWNLOAD_URL="https://github.com/fleeks-ai/fleeks-cli/releases/download/$VERSION/$BINARY"
TMP_FILE="/tmp/fleeks"

echo "Downloading Fleeks CLI..."
if command -v curl &> /dev/null; then
    curl -L -o "$TMP_FILE" "$DOWNLOAD_URL"
elif command -v wget &> /dev/null; then
    wget -O "$TMP_FILE" "$DOWNLOAD_URL"
else
    echo -e "${RED}Error: curl or wget is required${NC}"
    exit 1
fi

# Make executable
chmod +x "$TMP_FILE"

# Install
INSTALL_DIR="/usr/local/bin"
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_FILE" "$INSTALL_DIR/fleeks"
    echo -e "${GREEN}✓ Installed to $INSTALL_DIR/fleeks${NC}"
else
    echo -e "${YELLOW}Installing to $INSTALL_DIR requires sudo...${NC}"
    sudo mv "$TMP_FILE" "$INSTALL_DIR/fleeks"
    echo -e "${GREEN}✓ Installed to $INSTALL_DIR/fleeks${NC}"
fi

# Verify installation
echo ""
if command -v fleeks &> /dev/null; then
    fleeks --version
    echo ""
    echo -e "${GREEN}✓ Fleeks CLI installed successfully!${NC}"
    echo ""
    echo "Get started with:"
    echo "  fleeks auth login"
else
    echo -e "${RED}Installation failed${NC}"
    exit 1
fi
