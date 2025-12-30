#!/bin/bash

# Simple script to remove a specific version
# Usage: ./cleanup-version.sh <version>
# Example: ./cleanup-version.sh 0.1.5-test

VERSION=$1

if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 0.1.5-test"
  exit 1
fi

echo "Removing version: $VERSION"

# Remove versioned docs directory
if [ -d "versioned_docs/version-${VERSION}" ]; then
  rm -rf "versioned_docs/version-${VERSION}"
  echo "✓ Removed versioned_docs/version-${VERSION}"
fi

# Remove versioned sidebars file
if [ -f "versioned_sidebars/version-${VERSION}-sidebars.json" ]; then
  rm -f "versioned_sidebars/version-${VERSION}-sidebars.json"
  echo "✓ Removed versioned_sidebars/version-${VERSION}-sidebars.json"
fi

# Update versions.json
if [ -f "versions.json" ]; then
  # Use node to update JSON (more reliable than sed)
  node -e "
    const fs = require('fs');
    const versions = JSON.parse(fs.readFileSync('versions.json', 'utf-8'));
    const updated = versions.filter(v => v !== '$VERSION');
    fs.writeFileSync('versions.json', JSON.stringify(updated, null, 2) + '\n');
    console.log('✓ Updated versions.json');
  "
fi

echo "✓ Cleanup complete for version: $VERSION"

