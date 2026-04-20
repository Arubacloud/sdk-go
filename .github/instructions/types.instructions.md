---
applyTo: "pkg/types/**"
---

# Types package conventions

You are editing request/response data models in `pkg/types/`.

## Naming
- Request struct: `<PascalResource>Request`
- Single-item response: `<PascalResource>Response`
- Collection response: `<PascalResource>List` — must embed `ListResponse`

## Field rules
- All optional fields use pointer types: `*string`, `*int32`, `*bool`
- JSON tags are required: `json:"fieldName,omitempty"` — always `omitempty` for optional fields
- Required fields that must always be present on the wire omit `omitempty`
- Never use `interface{}` — use `map[string]interface{}` only for extension maps (`Extensions`)

## Validation helpers
- Add a `Validate<Resource>(projectID, resourceID string) error` helper in `pkg/types/utils.go` for any resource that needs pre-request validation
- Validation returns `fmt.Errorf("field cannot be empty")` — no wrapping, no `%w`

## Struct comments
- Every exported struct must have a doc comment starting with the type name
- Every exported field should have a comment if its meaning is not obvious from the name

## Embedding ListResponse
```go
type MyResourceList struct {
    ListResponse
    Items []MyResourceResponse `json:"items,omitempty"`
}
```
Do not redefine `Total`, `Offset`, `Limit` — they come from `ListResponse`.
