# Quick Testing Guide

## Step 1: Install Dependencies

```bash
# From project root
make docs-install

# Or manually
cd docs
npm install
```

## Step 2: Test Local Development

```bash
# From project root
make docs-serve

# Or manually
cd docs
npm start
```

Then open `http://localhost:3000` in your browser. You should see:
- Documentation at the root (no `/docs` prefix)
- All pages working
- Hot reload when editing files

Press `Ctrl+C` to stop the server.

## Step 3: Test Build

```bash
# From project root
make docs-build

# Or manually
cd docs
npm run build
```

This should complete without errors and create a `build/` directory.

## Step 4: Test Versioning

```bash
# Clean any existing versions first
cd docs
rm -rf versioned_docs versioned_sidebars versions.json

# Create a test version
make docs-version VERSION=0.1.5-test

# Or manually
npm run docs:version 0.1.5-test
```

**Expected result:**
- ✅ Creates `versioned_docs/version-0.1.5-test/` with markdown files
- ✅ Creates `versioned_sidebars/version-0.1.5-test-sidebars.json`
- ✅ Creates/updates `versions.json`

**If you get an error about "Cannot copy to a subdirectory":**
- The exclude patterns should prevent this
- Make sure `versioned_docs` doesn't exist before running
- Check `docusaurus.config.js` exclude patterns

## Step 5: Test Build with Versions

```bash
cd docs
npm run build
```

The build should succeed and include the versioned docs.

## Step 6: Clean Up Test Version

```bash
cd docs
rm -rf versioned_docs/version-0.1.5-test
rm -f versioned_sidebars/version-0.1.5-test-sidebars.json
rm -f versions.json
```

## All Tests Passed? ✅

If all steps complete successfully:
- ✅ Local development works
- ✅ Build works
- ✅ Versioning works
- ✅ Build with versions works

You're ready to use the documentation system!

