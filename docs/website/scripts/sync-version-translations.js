const fs = require('fs');
const path = require('path');

// Get the docs/website directory (parent of scripts)
// __dirname will be scripts/ when running from package.json
const docsWebsiteDir = path.resolve(__dirname, '..');
const versionsFile = path.resolve(docsWebsiteDir, 'versions.json');
const translationsSource = path.resolve(docsWebsiteDir, 'i18n', 'it', 'docusaurus-plugin-content-docs', 'current');
const i18nBase = path.resolve(docsWebsiteDir, 'i18n', 'it', 'docusaurus-plugin-content-docs');

// Read versions (support both array and file existence)
let versions = [];
if (fs.existsSync(versionsFile)) {
  versions = JSON.parse(fs.readFileSync(versionsFile, 'utf8'));
}
if (!Array.isArray(versions) || versions.length === 0) {
  console.log('No versions found in versions.json, nothing to sync.');
  process.exit(0);
}

// Get all files from current translations (only .md and current.json)
const sourceFiles = fs.readdirSync(translationsSource, { withFileTypes: true })
  .filter(dirent => dirent.isFile() && (dirent.name.endsWith('.md') || dirent.name === 'current.json'))
  .map(dirent => dirent.name);

console.log(`Syncing translations from 'current' to ${versions.length} versions...`);

// Copy translations to each version
versions.forEach(version => {
  const versionDir = path.join(i18nBase, `version-${version}`);

  // Create version directory if it doesn't exist
  if (!fs.existsSync(versionDir)) {
    fs.mkdirSync(versionDir, { recursive: true });
    console.log(`Created directory: ${versionDir}`);
  }

  // Copy each file
  sourceFiles.forEach(file => {
    const sourceFile = path.join(translationsSource, file);

    // For current.json, rename to version-{version}.json and update content
    if (file === 'current.json') {
      const destFileName = `version-${version}.json`;
      const finalDestFile = path.join(versionDir, destFileName);

      // Read and update the JSON content
      const content = JSON.parse(fs.readFileSync(sourceFile, 'utf8'));
      // Update version.label message if needed
      if (content['version.label']) {
        content['version.label'].message = version;
      }

      fs.writeFileSync(finalDestFile, JSON.stringify(content, null, 2) + '\n');
      console.log(`  Copied and updated ${file} -> version-${version}/${destFileName}`);
    } else {
      // For markdown files, copy as-is
      const finalDestFile = path.join(versionDir, file);
      fs.copyFileSync(sourceFile, finalDestFile);
      console.log(`  Copied ${file} -> version-${version}/${file}`);
    }
  });
});

console.log('Translation sync completed!');
