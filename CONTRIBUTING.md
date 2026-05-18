# Contributing to sdk-go

## Prerequisites

- Go 1.22 or later
- An Aruba Cloud account with API credentials (`CLIENT_ID`, `CLIENT_SECRET`)
- `make` available in your PATH

## Development workflow

```bash
# Install dependencies
make install

# Format + vet
make lint

# Run unit tests with coverage
make test

# Full check (tidy, lint, build, test)
make all
```

## Running end-to-end (acceptance) tests locally

The acceptance tests provision real Aruba Cloud resources against a live account.
They are intentionally **not run on every CI push** — see the CI section below.

### Step 1 — export credentials

```bash
export CLIENT_ID=<your-client-id>
export CLIENT_SECRET=<your-client-secret>
```

### Step 2 — create all resources

```bash
make e2e-create
```

This runs `examples/all-resources` in `create` mode, streams output to the
terminal **and** writes it to `create.log`.

At the end of a successful run the summary section prints:

```
=== SDK Example Complete ===
Successfully created resources:
- Project ID: <project-id>
...
```

Copy the **Project ID** — you will need it for the delete step.

### Step 3 — delete all resources

```bash
make e2e-delete PROJECT_ID=<project-id>
```

This re-runs the example in `delete` mode, streams output to the terminal
**and** writes it to `delete.log`.

> **Tip:** if the create run was interrupted and you need to clean up a
> partially-provisioned project, you can call `make e2e-delete` with the
> project ID at any time.

---

## CI / acceptance test pipeline

The repository has three GitHub Actions workflows:

| Workflow | File | Trigger |
|---|---|---|
| CI (unit tests, lint, build, security) | `.github/workflows/ci.yml` | push / PR to `main` or `develop` |
| Docs | `.github/workflows/docs.yml` | push to `main` |
| **Acceptance tests (E2E)** | `.github/workflows/e2e.yml` | **manual only** |

### Running the acceptance pipeline

1. Open the repository on GitHub.
2. Go to **Actions → Acceptance Tests (E2E)**.
3. Click **Run workflow** and confirm.

The pipeline requires two repository secrets to be set by a maintainer:

| Secret | Description |
|---|---|
| `E2E_CLIENT_ID` | Aruba Cloud API client ID |
| `E2E_CLIENT_SECRET` | Aruba Cloud API client secret |

The pipeline runs two sequential jobs:

- **E2E Create** — provisions all resources, uploads `create.log` as an artifact.
- **E2E Delete** — destroys the project created above (runs even if create
  partially failed, as long as a Project ID was captured), uploads `delete.log`.

Logs for both runs are retained as GitHub Actions artifacts for 7 days.

---

## Pull request guidelines

- Keep commits focused; one logical change per commit.
- Follow the existing code style — run `make lint` before pushing.
- Add or update unit tests for any behavioural change.
- Update `CHANGELOG.md` under the `Unreleased` section.
