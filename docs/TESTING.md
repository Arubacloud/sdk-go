# Testing Documentation

This guide explains how to test the documentation locally and verify versioning works.

## Prerequisites

1. Node.js 18+ installed
2. Dependencies installed: `make docs-install` or `cd docs && npm install`

## Local Development Testing

### 1. Start Development Server

```bash
# Using Makefile
make docs-serve

# Or directly
cd docs && npm start
```

The server will start at `http://localhost:3000`. You should see:
- Documentation at the root URL (no `/docs` prefix)
- All documentation pages accessible
- Hot reload when editing markdown files

### 2. Test Build

```bash
# Using Makefile
make docs-build

# Or directly
cd docs && npm run build
```

This will:
- Build the documentation for production
- Create a `build/` directory
- Validate all markdown links
- Check for broken references

### 3. Test Production Build Locally

```bash
# Using Makefile
make docs-serve-build

# Or directly
cd docs && npm run build && npm run serve
```

This builds and serves the production version, simulating what GitHub Pages will serve.

## Versioning Testing

### 1. Clean Previous Versions (if any)

```bash
cd docs
rm -rf versioned_docs versioned_sidebars versions.json
```

### 2. Create a Test Version

```bash
# Using Makefile
make docs-version VERSION=0.1.5-test

# Or directly
cd docs && npm run docs:version 0.1.5-test
```

**Expected behavior:**
- Creates `versioned_docs/version-0.1.5-test/` with all markdown files
- Creates `versioned_sidebars/version-0.1.5-test-sidebars.json`
- Creates/updates `versions.json` with the new version

**If you see an error about copying to a subdirectory:**
- The exclude patterns in `docusaurus.config.js` should prevent this
- Make sure `versioned_docs` doesn't exist before running versioning
- Check that all non-markdown files are properly excluded

### 3. Verify Versioned Files

```bash
cd docs

# Check versioned docs exist
ls -la versioned_docs/version-0.1.5-test/

# Check sidebar exists
cat versioned_sidebars/version-0.1.5-test-sidebars.json

# Check versions.json
cat versions.json
```

### 4. Test Build with Versions

```bash
cd docs && npm run build
```

The build should include:
- Current version (latest docs)
- Versioned docs in `build/docs/version-0.1.5-test/`

### 5. Enable Version Dropdown (after first version)

Once you have at least one version, uncomment the version dropdown in `docusaurus.config.js`:

```javascript
{
  type: 'docsVersionDropdown',
  position: 'right',
  dropdownActiveClassDisabled: true,
},
```

Then test locally:
```bash
make docs-serve
```

You should see a version dropdown in the navbar.

### 6. Clean Up Test Version

```bash
cd docs
rm -rf versioned_docs/version-0.1.5-test
rm -f versioned_sidebars/version-0.1.5-test-sidebars.json
# Edit versions.json to remove the test version, or delete it if it's the only version
```

## Automated Testing Script

For Linux/Mac/WSL, you can use the test script:

```bash
cd docs
chmod +x test-versioning.sh
./test-versioning.sh
```

This script will:
1. Clean up existing versions
2. Create a test version
3. Verify all files were created correctly
4. Test the build
5. Report success/failure

## Common Issues

### Issue: "Cannot copy to a subdirectory of itself"

**Solution:** This happens when Docusaurus tries to copy the docs folder into itself. The exclude patterns should prevent this. Make sure:
- `versioned_docs` doesn't exist before versioning
- All non-markdown files are in the exclude list
- The config doesn't have `path: '.'` set

### Issue: Version dropdown doesn't appear

**Solution:** 
- Make sure at least one version exists
- Uncomment the `docsVersionDropdown` in `docusaurus.config.js`
- Restart the dev server

### Issue: Build fails with versioning errors

**Solution:**
- Check that `versions.json` is valid JSON
- Verify all versioned docs exist
- Check that sidebar files match the versioned docs

## CI/CD Testing

The GitHub Actions workflows will automatically:
- Test builds on every push
- Create versions when tags are pushed
- Deploy to GitHub Pages

Check `.github/workflows/docs.yml` and `.github/workflows/docs-version.yml` for details.

