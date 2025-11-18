#!/bin/bash
# Script to export GPG keys for apt-bundle repository
# Run this after generating the key with: gpg --batch --gen-key .github/workflows/scripts/gpg-keygen.conf

set -e

# Find the key ID using machine-readable format for robust parsing
KEY_ID=$(gpg --list-keys --keyid-format LONG --with-colons "maintainers@apt-bundle.org" | awk -F: '/^pub:/ {print $5; exit}')

if [ -z "$KEY_ID" ]; then
    echo "Error: Could not find GPG key for maintainers@apt-bundle.org"
    exit 1
fi

echo "Found key ID: ${KEY_ID}"
echo ""

# Export public key
echo "Exporting public key to public.key..."
gpg --armor --export "${KEY_ID}" > public.key
echo "✓ Public key exported"

# Export private key (for GitHub Secrets)
echo "Exporting private key to private.key..."
gpg --armor --export-secret-keys "${KEY_ID}" > private.key
echo "✓ Private key exported"

# Show key fingerprint
echo ""
echo "Key fingerprint:"
gpg --fingerprint "${KEY_ID}"

KEY_ID_SHORT=$(echo "${KEY_ID}" | tail -c 9)

echo ""
echo "Next steps:"
echo "1. Add private.key content to GitHub Secret: GPG_PRIVATE_KEY"
echo "2. Add key ID (last 8 chars: ${KEY_ID_SHORT}) to GitHub Secret: GPG_KEY_ID"
echo "3. Commit public.key to repository (or host at apt-bundle.org/public.key)"
echo "4. Securely delete private.key after adding to GitHub Secrets"

