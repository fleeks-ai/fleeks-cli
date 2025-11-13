#!/bin/bash
set -e

VERSION="${1:-dev}"
COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS="-s -w -X 'main.Version=$VERSION' -X 'main.BuildDate=$BUILD_TIME'"

echo "Building Fleeks CLI v$VERSION..."
echo "Commit: $COMMIT"
echo "Build Time: $BUILD_TIME"
echo ""

# Create bin directory
mkdir -p bin

# Windows AMD64
echo "Building Windows AMD64..."
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$LDFLAGS" -o bin/fleeks-windows-amd64.exe .

# macOS Intel
echo "Building macOS Intel..."
GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$LDFLAGS" -o bin/fleeks-darwin-amd64 .

# macOS Apple Silicon
echo "Building macOS Apple Silicon..."
GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="$LDFLAGS" -o bin/fleeks-darwin-arm64 .

# Linux AMD64
echo "Building Linux AMD64..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="$LDFLAGS" -o bin/fleeks-linux-amd64 .

# Linux ARM64 (for Raspberry Pi, ARM servers)
echo "Building Linux ARM64..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="$LDFLAGS" -o bin/fleeks-linux-arm64 .

echo ""
echo "✓ Build complete! Binaries in bin/"
echo ""
ls -lh bin/
