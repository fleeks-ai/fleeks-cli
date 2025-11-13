#!/bin/bash
set -e

VERSION="$1"

if [ -z "$VERSION" ]; then
    echo "Usage: ./create-release.sh v1.0.0"
    exit 1
fi

echo "Creating release $VERSION..."
echo ""

# Check if gh CLI is installed
if ! command -v gh &> /dev/null; then
    echo "Error: GitHub CLI (gh) is not installed"
    echo "Install from: https://cli.github.com/"
    exit 1
fi

# Check if logged in
if ! gh auth status &> /dev/null; then
    echo "Error: Not logged in to GitHub CLI"
    echo "Run: gh auth login"
    exit 1
fi

# 1. Build all binaries
echo "Step 1: Building binaries..."
./scripts/build-all.sh "${VERSION#v}"
echo ""

# 2. Generate checksums
echo "Step 2: Generating checksums..."
./scripts/generate-checksums.sh
echo ""

# 3. Create GitHub release
echo "Step 3: Creating GitHub release..."

# Check if release already exists
if gh release view "$VERSION" &> /dev/null; then
    echo "Release $VERSION already exists!"
    read -p "Delete and recreate? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        gh release delete "$VERSION" --yes
    else
        exit 1
    fi
fi

# Create release notes
cat > /tmp/release-notes.md <<EOF
## Fleeks CLI $VERSION

### Installation

#### macOS
\`\`\`bash
# Intel Macs
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$VERSION/fleeks-darwin-amd64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$VERSION/fleeks-darwin-arm64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/
\`\`\`

#### Linux
\`\`\`bash
# AMD64
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$VERSION/fleeks-linux-amd64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/

# ARM64
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$VERSION/fleeks-linux-arm64 -o fleeks
chmod +x fleeks
sudo mv fleeks /usr/local/bin/
\`\`\`

#### Windows (PowerShell)
\`\`\`powershell
Invoke-WebRequest -Uri "https://github.com/fleeks-ai/fleeks-cli/releases/download/$VERSION/fleeks-windows-amd64.exe" -OutFile fleeks.exe
# Add to PATH or move to a directory in your PATH
\`\`\`

### Package Managers

\`\`\`bash
# Homebrew (macOS/Linux)
brew install fleeks-ai/fleeks/fleeks

# NPM (Cross-platform)
npm install -g @fleeks/cli
\`\`\`

### What's Changed

See the [CHANGELOG](https://github.com/fleeks-ai/fleeks-cli/blob/main/CHANGELOG.md) for full details.

### Verification

Verify your download with SHA256 checksums:
\`\`\`bash
# Download checksums
curl -L https://github.com/fleeks-ai/fleeks-cli/releases/download/$VERSION/SHA256SUMS -o SHA256SUMS

# Verify (Linux/macOS)
sha256sum -c SHA256SUMS --ignore-missing

# Verify (Windows)
# Compare hash in SHA256SUMS with: (Get-FileHash fleeks.exe).Hash
\`\`\`
EOF

# Create the release
gh release create "$VERSION" \
    bin/fleeks-windows-amd64.exe \
    bin/fleeks-darwin-amd64 \
    bin/fleeks-darwin-arm64 \
    bin/fleeks-linux-amd64 \
    bin/fleeks-linux-arm64 \
    bin/SHA256SUMS \
    --title "Release $VERSION" \
    --notes-file /tmp/release-notes.md

echo ""
echo "✓ Release $VERSION created successfully!"
echo "View at: https://github.com/fleeks-ai/fleeks-cli/releases/tag/$VERSION"
echo ""
echo "Next steps:"
echo "  1. Update Homebrew formula: ./scripts/update-homebrew.sh $VERSION"
echo "  2. Update NPM package: ./scripts/publish-npm.sh $VERSION"
echo "  3. Announce on social media and documentation"
