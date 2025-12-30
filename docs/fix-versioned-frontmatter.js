#!/usr/bin/env node

/**
 * Script to fix front matter in existing versioned docs
 * Removes slug from front matter (slug is only valid for current version)
 */

const fs = require('fs');
const path = require('path');

const docsDir = __dirname;
const versionedDocsDir = path.join(docsDir, 'versioned_docs');

if (!fs.existsSync(versionedDocsDir)) {
  console.log('No versioned_docs directory found, nothing to fix');
  process.exit(0);
}

// Helper function to remove all front matter from versioned docs
// Versioned docs should not have front matter to avoid Docusaurus warnings
function removeFrontMatter(content) {
  // Match front matter (YAML between --- markers)
  const frontMatterRegex = /^---\s*\n([\s\S]*?)\n---\s*\n/;
  const match = content.match(frontMatterRegex);
  
  if (!match) {
    return content; // No front matter, return as-is
  }
  
  // Remove the entire front matter block
  const restOfContent = content.substring(match[0].length);
  return restOfContent;
}

// Find all markdown files in versioned_docs
function findMarkdownFiles(dir) {
  const files = [];
  const entries = fs.readdirSync(dir, { withFileTypes: true });
  
  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...findMarkdownFiles(fullPath));
    } else if (entry.isFile() && entry.name.endsWith('.md')) {
      files.push(fullPath);
    }
  }
  
  return files;
}

const markdownFiles = findMarkdownFiles(versionedDocsDir);

if (markdownFiles.length === 0) {
  console.log('No markdown files found in versioned_docs');
  process.exit(0);
}

console.log(`Found ${markdownFiles.length} markdown file(s) in versioned_docs`);
let fixedCount = 0;

markdownFiles.forEach(filePath => {
  const content = fs.readFileSync(filePath, 'utf8');
  const originalContent = content;
  const fixedContent = removeFrontMatter(content);
  
  if (originalContent !== fixedContent) {
    fs.writeFileSync(filePath, fixedContent, 'utf8');
    console.log(`  ✓ Fixed ${path.relative(versionedDocsDir, filePath)}`);
    fixedCount++;
  }
});

if (fixedCount > 0) {
  console.log(`\n✅ Fixed front matter in ${fixedCount} file(s)`);
} else {
  console.log('\n✅ No files needed fixing');
}

