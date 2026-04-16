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

The CI pipeline uses golangci-lint v2.11.4 (timeout 5 min). Active linters: `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`, `misspell`, `unparam`, `unconvert`, `goconst`, `gocyclo`. Generated code in `pkg/generated` and examples in `cmd/example` are excluded.

## Reference documentation

- `docs/website/docs/options.md` — Client configuration reference
- `docs/website/docs/filters.md` — Filter syntax and examples
- `docs/website/docs/response-handling.md` — Error handling patterns
- `docs/website/docs/resources.md` — Full API groups and resource listing
- `docs/website/docs/types.md` — Data model guide
