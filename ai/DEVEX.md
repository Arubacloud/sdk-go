# Development Experience

## Build & test commands

```bash
make build        # go build ./...
make test         # run all tests with race detection and coverage
make test-short   # quick tests without coverage
make lint         # go fmt + go vet
make verify       # lint + test (recommended before committing)
make all          # tidy + lint + build + test
```

Run a single test:
```bash
go test -v -run TestName ./pkg/...
go test -v -run TestName ./internal/...
```

## Linting

The CI pipeline uses golangci-lint v2.11.4 (timeout 5 min). Active linters: `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`, `misspell`, `unparam`, `unconvert`, `goconst`, `gocyclo`. Examples in `examples/all-resources` are excluded. The only generated code in the repo is mockgen output (`zz_mock_*_test.go` files); there is no `pkg/generated` directory.

## Reference documentation

- `docs/website/docs/intro.md` — Quick start: install, first client, first resource
- `docs/website/docs/walkthrough.md` — Full CRUD walkthrough with real resources
- `docs/website/docs/resources.md` — Full API groups and resource listing
- `docs/website/docs/filters.md` — Filter syntax and examples
- `docs/website/docs/response-handling.md` — Error handling patterns
- `docs/website/docs/async.md` — Async polling, wait helpers, and advanced WaitFor usage
- `docs/website/docs/multitenancy.md` — Multi-tenant client management
- `docs/website/docs/options.md` — Client configuration reference
- `pkg/aruba/aliases.go` — Typed enum constants for all resources

## Docs site

The Docusaurus site lives in `docs/website/`. Use `make docs-serve` (local dev) and `make docs-build` (production build). See `docs/README.md` for the full versioning and Italian-translation workflow.
