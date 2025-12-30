#!/usr/bin/env node

/**
 * Custom script to create Docusaurus documentation versions
 * This script manually copies only markdown files to avoid circular copy issues
 * when the config file is in the docs folder.
 */

const fs = require('fs');
const path = require('path');

const version = process.argv[2];

if (!version) {
  console.error('Error: Version number is required');
  console.error('Usage: node version-docs.js <version>');
  process.exit(1);
}

const docsDir = __dirname;
const versionedDocsDir = path.join(docsDir, 'versioned_docs');
const versionDir = path.join(versionedDocsDir, `version-${version}`);
const versionedSidebarsDir = path.join(docsDir, 'versioned_sidebars');
const versionsJsonPath = path.join(docsDir, 'versions.json');

// Find all markdown files in docs directory (exclude README and test files)
const allFiles = fs.readdirSync(docsDir);
const markdownFiles = allFiles.filter(file => {
  return file.endsWith('.md') && 
         file !== 'README.md' && 
         !file.includes('TEST') && 
         !file.includes('test');
});

console.log(`Creating version ${version}...`);

// Create directories
if (!fs.existsSync(versionedDocsDir)) {
  fs.mkdirSync(versionedDocsDir, { recursive: true });
}
if (!fs.existsSync(versionedSidebarsDir)) {
  fs.mkdirSync(versionedSidebarsDir, { recursive: true });
}

// Create version directory
if (!fs.existsSync(versionDir)) {
  fs.mkdirSync(versionDir, { recursive: true });
} else {
  console.error(`Error: Version ${version} already exists`);
  process.exit(1);
}

// Copy only markdown files
let copiedCount = 0;
markdownFiles.forEach(file => {
  const srcPath = path.join(docsDir, file);
  if (fs.existsSync(srcPath)) {
    const destPath = path.join(versionDir, file);
    fs.copyFileSync(srcPath, destPath);
    console.log(`  ✓ Copied ${file}`);
    copiedCount++;
  } else {
    console.warn(`  ⚠ File ${file} not found, skipping`);
  }
});

if (copiedCount === 0) {
  console.error('Error: No markdown files were copied');
  process.exit(1);
}

// Create versioned sidebar
const sidebarPath = path.join(docsDir, 'sidebars.js');
const versionedSidebarPath = path.join(versionedSidebarsDir, `version-${version}-sidebars.json`);

// Read the current sidebar structure
const sidebarModule = require(sidebarPath);
const versionedSidebar = {
  tutorialSidebar: sidebarModule.tutorialSidebar.map(item => {
    // For versioned sidebars, we need to ensure the IDs match the copied files
    if (item.type === 'doc') {
      return {
        type: 'doc',
        id: item.id,
        label: item.label,
      };
    }
    return item;
  }),
};

fs.writeFileSync(versionedSidebarPath, JSON.stringify(versionedSidebar, null, 2));
console.log(`  ✓ Created sidebar configuration`);

// Update versions.json
let versions = [];
if (fs.existsSync(versionsJsonPath)) {
  try {
    versions = JSON.parse(fs.readFileSync(versionsJsonPath, 'utf8'));
  } catch (e) {
    console.warn('  ⚠ Could not parse existing versions.json, creating new one');
  }
}

if (versions.includes(version)) {
  console.warn(`  ⚠ Version ${version} already exists in versions.json`);
} else {
  versions.unshift(version); // Add to beginning (newest first)
  fs.writeFileSync(versionsJsonPath, JSON.stringify(versions, null, 2));
  console.log(`  ✓ Updated versions.json`);
}

console.log(`\n✅ Version ${version} created successfully!`);
console.log(`\nNext steps:`);
console.log(`  1. Review the versioned files in ${path.relative(process.cwd(), versionDir)}`);
console.log(`  2. Enable versioning in docusaurus.config.js by uncommenting docsVersionDropdown`);
console.log(`  3. Test locally with: make docs-serve`);

