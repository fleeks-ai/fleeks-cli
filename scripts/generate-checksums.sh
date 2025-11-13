#!/bin/bash
set -e

cd bin

echo "Generating SHA256 checksums..."

# Generate checksums for all binaries
sha256sum fleeks-* > SHA256SUMS 2>/dev/null || shasum -a 256 fleeks-* > SHA256SUMS

# Display checksums
echo ""
echo "Checksums generated:"
cat SHA256SUMS

# Optionally sign with GPG
if command -v gpg &> /dev/null; then
    echo ""
    read -p "Sign checksums with GPG? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        gpg --detach-sign --armor SHA256SUMS
        echo "✓ Created SHA256SUMS.asc"
    fi
fi

echo ""
echo "✓ Checksums complete!"
