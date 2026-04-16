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

The CI pipeline uses golangci-lint v2.11.4 (timeout 5 min). Active linters: `errcheck`, `govet`, `staticcheck`, `unused`, `misspell`, `unparam`, `unconvert`, `goconst`, `gocyclo`. Generated code in `pkg/generated` and examples in `cmd/example` are excluded.

## Reference documentation

- `doc/OPTIONS.md` — Client configuration reference
- `doc/FILTERS.md` — Filter syntax and examples
- `doc/RESPONSE_HANDLING.md` — Error handling patterns
- `doc/RESOURCES.md` — Full API groups and resource listing
- `doc/TYPES.md` — Data model guide
