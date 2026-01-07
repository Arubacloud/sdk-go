#!/usr/bin/env node

/**
 * Cleanup script to remove old versioned docs, keeping only the last N releases
 * Default: keeps last 5 releases
 * Usage: node cleanup-versions.js [--dry-run]
 * Environment: KEEP_LAST=5 (default) - number of releases to keep
 */

const fs = require('fs');
const path = require('path');
const https = require('https');

const DRY_RUN = process.argv.includes('--dry-run');

function getGitHubReleases(keepLast = 5) {
  return new Promise((resolve, reject) => {
    // Try to get repo from git config, or use environment variable
    let repo = process.env.GITHUB_REPOSITORY;
    
    if (!repo) {
      // Try to get from git remote
      try {
        const { execSync } = require('child_process');
        const remoteUrl = execSync('git config --get remote.origin.url', { encoding: 'utf-8' }).trim();
        // Extract owner/repo from git URL (handles both https and ssh formats)
        const match = remoteUrl.match(/(?:github\.com[/:]|git@github\.com:)([^/]+)\/([^/]+?)(?:\.git)?$/);
        if (match) {
          repo = `${match[1]}/${match[2]}`;
        }
      } catch (error) {
        // Ignore git errors
      }
    }
    
    if (!repo) {
      console.error('Error: Could not determine repository. Set GITHUB_REPOSITORY environment variable.');
      console.error('Example: GITHUB_REPOSITORY=Arubacloud/sdk-go node cleanup-versions.js');
      process.exit(1);
    }
    
    // Get GitHub token from environment (optional, but recommended for rate limits)
    const token = process.env.GITHUB_TOKEN;
    const authHeader = token ? `token ${token}` : '';
    
    const options = {
      hostname: 'api.github.com',
      path: `/repos/${repo}/releases`,
      method: 'GET',
      headers: {
        'User-Agent': 'Node.js',
        ...(authHeader && { 'Authorization': authHeader })
      }
    };
    
    https.get(options, (res) => {
      let data = '';
      
      res.on('data', (chunk) => {
        data += chunk;
      });
      
      res.on('end', () => {
        if (res.statusCode !== 200) {
          console.error(`Error fetching releases: ${res.statusCode} ${res.statusMessage}`);
          if (res.statusCode === 404) {
            console.error(`Repository not found: ${repo}`);
            console.error('Make sure the repository exists and is accessible.');
          } else if (res.statusCode === 401) {
            console.error('Authentication failed. Set GITHUB_TOKEN environment variable for private repos or higher rate limits.');
          }
          reject(new Error(`HTTP ${res.statusCode}`));
          return;
        }
        
        try {
          const releases = JSON.parse(data);
          
          // Sort by semantic version (newest first)
          // Extract version numbers and sort them properly
          const sortedReleases = releases
            .filter(r => r.published_at) // Only published releases
            .map(r => ({
              tag: r.tag_name,
              version: r.tag_name.replace(/^v/, '') // Remove 'v' prefix
            }))
            .sort((a, b) => {
              // Semantic version comparison
              const aParts = a.version.split('.').map(Number);
              const bParts = b.version.split('.').map(Number);
              
              // Compare major, minor, patch
              for (let i = 0; i < Math.max(aParts.length, bParts.length); i++) {
                const aPart = aParts[i] || 0;
                const bPart = bParts[i] || 0;
                if (aPart !== bPart) {
                  return bPart - aPart; // Descending order (newest first)
                }
              }
              return 0;
            })
            .map(r => r.version); // Extract just the version string
          
          // Keep only the last N releases
          const releasesToKeep = sortedReleases.slice(0, keepLast);
          console.log(`Keeping last ${keepLast} release(s): ${releasesToKeep.join(', ')}`);
          
          resolve({
            all: sortedReleases,
            toKeep: releasesToKeep
          });
        } catch (error) {
          console.error('Error parsing releases data:', error.message);
          reject(error);
        }
      });
    }).on('error', (error) => {
      console.error('Error fetching releases:', error.message);
      console.error('Make sure you have internet connectivity.');
      reject(error);
    });
  });
}

function getLocalVersions() {
  const versionsFile = path.join(__dirname, 'versions.json');
  if (!fs.existsSync(versionsFile)) {
    return [];
  }
  return JSON.parse(fs.readFileSync(versionsFile, 'utf-8'));
}

function removeVersion(version) {
  const versionedDocsDir = path.join(__dirname, 'versioned_docs', `version-${version}`);
  const versionedSidebarsFile = path.join(__dirname, 'versioned_sidebars', `version-${version}-sidebars.json`);
  const versionsFile = path.join(__dirname, 'versions.json');
  
  console.log(`\nRemoving version: ${version}`);
  
  if (DRY_RUN) {
    console.log(`  [DRY RUN] Would remove: ${versionedDocsDir}`);
    console.log(`  [DRY RUN] Would remove: ${versionedSidebarsFile}`);
    return;
  }
  
  // Remove versioned docs directory
  if (fs.existsSync(versionedDocsDir)) {
    fs.rmSync(versionedDocsDir, { recursive: true, force: true });
    console.log(`  ✓ Removed: ${versionedDocsDir}`);
  }
  
  // Remove versioned sidebars file
  if (fs.existsSync(versionedSidebarsFile)) {
    fs.unlinkSync(versionedSidebarsFile);
    console.log(`  ✓ Removed: ${versionedSidebarsFile}`);
  }
  
  // Update versions.json
  const versions = getLocalVersions();
  const updatedVersions = versions.filter(v => v !== version);
  fs.writeFileSync(versionsFile, JSON.stringify(updatedVersions, null, 2) + '\n');
  console.log(`  ✓ Updated: ${versionsFile}`);
}

async function main() {
  const KEEP_LAST = parseInt(process.env.KEEP_LAST || '5', 10);
  
  console.log('Cleaning up versioned docs...\n');
  console.log(`Keeping only the last ${KEEP_LAST} release(s)\n`);
  
  if (DRY_RUN) {
    console.log('🔍 DRY RUN MODE - No files will be deleted\n');
  }
  
  // Get versions from GitHub releases
  let releases;
  try {
    releases = await getGitHubReleases(KEEP_LAST);
  } catch (error) {
    console.error('\nCould not fetch releases from GitHub API.');
    console.error('Make sure you have internet connectivity.');
    console.error('For private repos or higher rate limits, set GITHUB_TOKEN environment variable.');
    console.error('Alternatively, you can manually edit versions.json and remove old versions.');
    process.exit(1);
  }
  
  // Get local versions
  const localVersions = getLocalVersions();
  console.log(`\nFound ${localVersions.length} local version(s): ${localVersions.join(', ')}`);
  
  // Find versions to remove:
  // 1. Versions not in any release
  // 2. Versions older than the last KEEP_LAST releases
  const versionsToRemove = localVersions.filter(v => !releases.toKeep.includes(v));
  
  if (versionsToRemove.length === 0) {
    console.log('\n✓ All local versions are in the last releases. No cleanup needed.');
    return;
  }
  
  console.log(`\n⚠️  Found ${versionsToRemove.length} version(s) to remove: ${versionsToRemove.join(', ')}`);
  if (releases.all.length > KEEP_LAST) {
    const removedReleases = releases.all.slice(KEEP_LAST);
    console.log(`   (Removing older releases: ${removedReleases.join(', ')})`);
  }
  
  // Remove each invalid/old version
  versionsToRemove.forEach(version => {
    removeVersion(version);
  });
  
  // Update versions.json to only include versions to keep
  const updatedVersions = localVersions.filter(v => releases.toKeep.includes(v));
  const versionsFile = path.join(__dirname, 'versions.json');
  if (!DRY_RUN) {
    fs.writeFileSync(versionsFile, JSON.stringify(updatedVersions, null, 2) + '\n');
    console.log(`\n✓ Updated versions.json to keep only: ${updatedVersions.join(', ')}`);
  }
  
  console.log(`\n✓ Cleanup complete! Removed ${versionsToRemove.length} version(s).`);
}

main().catch(error => {
  console.error('Fatal error:', error.message);
  process.exit(1);
});

