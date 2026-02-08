# Documentation

This directory contains the documentation for the Aruba Cloud SDK for Go.

## Development

### Prerequisites

- Node.js 18+ and npm

### Local Development (Without Versions)

When you first set up the documentation, versioning is disabled by default. This allows you to develop and test locally without needing to create versions first.

**Using Make (recommended):**

1. Install dependencies:
   ```bash
   make docs-install
   ```

2. Start the development server:
   ```bash
   make docs-serve
   # or simply
   make docs
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser.

**Using npm directly:**

1. Install dependencies:
   ```bash
   cd docs
   npm install
   ```

2. Start the development server:
   ```bash
   npm start
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser.

### Testing the Build Locally

To test the production build locally (simulates what CI does):

```bash
# Using Make (builds and serves in one command)
make docs-serve-build

# Or step by step
make docs-build    # Build the documentation
make docs-test     # Test and validate
cd docs && npm run serve  # Serve the built site
```

The built files will be in the `docs/build/` directory. The `docs-serve-build` target builds and serves the production version, which is useful for testing exactly what will be deployed.

### Building

To build the documentation for production:

```bash
make docs-build
# or
cd docs && npm run build
```

The built files will be in the `docs/build/` directory.

## Versioning

The documentation supports versioning, allowing multiple versions to be maintained simultaneously.

### Testing with Versions Enabled

Once you've created your first version and enabled versioning:

1. **Enable versioning in the configuration file**:
   - Remove `disableVersioning: true,` from the docs configuration
   - Uncomment the `docsVersionDropdown` in the navbar items

2. **Test locally**:
   ```bash
   make docs-serve
   ```
   The version dropdown will appear in the navbar, and you can switch between versions.

### Creating a New Version

When you release a new version of the SDK:

1. **First version only**: Enable versioning in the configuration file:
   - Remove `disableVersioning: true,` from the docs configuration
   - Uncomment the `docsVersionDropdown` in the navbar items

2. Create a git tag with the version number (e.g., `v1.0.0`):
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. The GitHub Actions workflow will automatically:
   - Extract the version from the tag
   - Create a versioned copy of the current documentation
   - Commit the versioned docs to the main branch

Alternatively, you can create a version manually:

```bash
npm run docs:version <version>
```

For example:
```bash
npm run docs:version 1.0.0
```

This will:
- Copy the current `docs/` content to `versioned_docs/version-1.0.0/`
- Create a versioned sidebar at `versioned_sidebars/version-1.0.0-sidebars.json`
- Update `versions.json` with the new version

### Version Structure

- **Current/Next**: The latest documentation (from `docs/` folder) is labeled as "Next" and is accessible at the root path
- **Versioned**: Older versions are stored in `versioned_docs/` and accessible at `/version-<version>/`

### Version Dropdown

The navbar includes a version dropdown that allows users to switch between different documentation versions.

### Syncing Italian Translations for All Versions

The site is available in English (default) and Italian. When you add a new version (e.g. after running the release workflow or `npm run docs:version`), Italian content for that version must be created so the Italian locale does not fall back to English.

From the `docs/website` directory, run:

```bash
npm run sync-version-translations
```

This copies the current Italian docs from `i18n/it/docusaurus-plugin-content-docs/current/` to each version listed in `versions.json` (e.g. `version-0.1.16`, `version-0.1.17`, …). Run it after adding a new version so every release has Italian content.

## Deployment

The documentation is automatically deployed to GitHub Pages when changes are pushed to the `main` branch. The deployment is handled by the `.github/workflows/docs.yml` workflow.

## CI/CD

The documentation has two GitHub Actions workflows:

1. **docs.yml**: Builds, tests, and deploys the documentation
   - Runs on pushes to `main` and pull requests
   - Tests for broken links and validates markdown
   - Deploys to GitHub Pages on `main` branch

2. **docs-version.yml**: Creates versioned documentation
   - Triggers on version tags (e.g., `v1.0.0`)
   - Creates a versioned copy of the documentation
   - Commits the versioned docs back to the repository
