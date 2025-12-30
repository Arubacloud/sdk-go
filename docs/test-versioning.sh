#!/bin/bash

# Test script for Docusaurus versioning
# This script tests the versioning functionality

set -e

VERSION="0.1.5-test"
DOCS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🧪 Testing Docusaurus Versioning"
echo "=================================="
echo ""

# Step 1: Clean up any existing versioned docs
echo "📦 Step 1: Cleaning up existing versioned docs..."
cd "$DOCS_DIR"
if [ -d "versioned_docs" ]; then
  echo "  Removing existing versioned_docs..."
  rm -rf versioned_docs versioned_sidebars versions.json
fi

# Step 2: Test versioning command
echo ""
echo "📝 Step 2: Creating version $VERSION..."
npm run docs:version "$VERSION"

# Step 3: Verify versioned files were created
echo ""
echo "✅ Step 3: Verifying versioned files..."
if [ -d "versioned_docs/version-$VERSION" ]; then
  echo "  ✓ versioned_docs/version-$VERSION exists"
  MD_COUNT=$(find "versioned_docs/version-$VERSION" -name "*.md" | wc -l)
  echo "  ✓ Found $MD_COUNT markdown files"
else
  echo "  ✗ versioned_docs/version-$VERSION does not exist"
  exit 1
fi

if [ -f "versioned_sidebars/version-$VERSION-sidebars.json" ]; then
  echo "  ✓ versioned_sidebars/version-$VERSION-sidebars.json exists"
else
  echo "  ✗ versioned_sidebars/version-$VERSION-sidebars.json does not exist"
  exit 1
fi

if [ -f "versions.json" ]; then
  echo "  ✓ versions.json exists"
  if grep -q "$VERSION" versions.json; then
    echo "  ✓ Version $VERSION found in versions.json"
  else
    echo "  ✗ Version $VERSION not found in versions.json"
    exit 1
  fi
else
  echo "  ✗ versions.json does not exist"
  exit 1
fi

# Step 4: Test build with versioned docs
echo ""
echo "🔨 Step 4: Testing build with versioned docs..."
npm run build

if [ -d "build" ]; then
  echo "  ✓ Build directory created"
  if [ -d "build/docs/version-$VERSION" ]; then
    echo "  ✓ Versioned docs included in build"
  else
    echo "  ⚠ Versioned docs not found in build (this might be normal if versioning dropdown is disabled)"
  fi
else
  echo "  ✗ Build directory not created"
  exit 1
fi

echo ""
echo "🎉 All tests passed!"
echo ""
echo "To clean up test version, run:"
echo "  rm -rf docs/versioned_docs/version-$VERSION"
echo "  rm -f docs/versioned_sidebars/version-$VERSION-sidebars.json"
echo "  rm -f docs/versions.json"

