package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/project"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ---- Wrapper ----

// Project is the wrapper for an Aruba Cloud project.
// Construct with aruba.Project(); pass it to ProjectClient.Create / .Update,
// or receive it from .Get / .List.
type Project struct {
	errMixin
	metadataMixin
	projectScopedMixin    // self-referential: ProjectID() returns this project's own ID after Get/Create
	responseMetadataMixin // promotes ID(), RespURI(), CreatedAt(), …
	httpEnvelopeMixin

	description *string
	defaultProj bool // ProjectPropertiesRequest.Default is plain bool — no tri-state needed
}

// Setters — chainable, general → specific

// Named sets the project name.
func (p *Project) Named(n string) *Project { p.named(n); return p }

// AddTag adds a tag (deduped).
func (p *Project) AddTag(t string) *Project { p.addTag(t); return p }

// RemoveTag removes a tag.
func (p *Project) RemoveTag(t string) *Project { p.removeTag(t); return p }

// ReplaceTags overwrites the tag list.
func (p *Project) ReplaceTags(ts ...string) *Project { p.replaceTags(ts...); return p }

// WithDescription sets the project description.
func (p *Project) WithDescription(d string) *Project { p.description = &d; return p }

// AsDefault marks the project as the account default.
func (p *Project) AsDefault() *Project { p.defaultProj = true; return p }

// NotDefault marks the project as not the account default.
func (p *Project) NotDefault() *Project { p.defaultProj = false; return p }

// Getters — general → specific

// URI satisfies Ref; returns the server-assigned URI, or "" before the first reply.
// ID() is promoted from responseMetadataMixin and needs no override.
func (p *Project) URI() string { return p.RespURI() }

// Description returns the description set via WithDescription, or "" if unset.
func (p *Project) Description() string {
	if p.description == nil {
		return ""
	}
	return *p.description
}

// IsDefault returns true if this project is marked as the account default.
func (p *Project) IsDefault() bool { return p.defaultProj }

// Wire converters

// toRequest assembles the Create/Update body from current setter state. Defaults are applied at the wire boundary.
func (p *Project) toRequest() types.ProjectRequest {
	return types.ProjectRequest{
		Metadata: p.toMetadata(),
		Properties: types.ProjectPropertiesRequest{
			Description: p.description,
			Default:     p.defaultProj,
		},
	}
}

// fromResponse hydrates the wrapper from a server reply. Nil-safe.
func (p *Project) fromResponse(resp *types.ProjectResponse) {
	if resp == nil {
		return
	}
	p.setMeta(&resp.Metadata)
	p.named(projectDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		p.replaceTags(resp.Metadata.Tags...)
	}
	p.description = resp.Properties.Description
	p.defaultProj = resp.Properties.Default
	// Seed our own projectID so that ProjectID() works immediately after Create/Get,
	// enabling child-resource setters that assert withProjectID.
	if resp.Metadata.ID != nil {
		p.projectID = *resp.Metadata.ID
	}
}

func projectDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ---- Low-level client interface ----

// projectLowLevelClient is the package-internal seam the adapter consumes.
// Satisfied by *project.projectsClientImpl. Defined here so tests can substitute
// a fake without depending on internal/clients/project test code.
type projectLowLevelClient interface {
	List(ctx context.Context, params *types.RequestParameters) (*types.Response[types.ProjectList], error)
	Get(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	Create(ctx context.Context, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	Update(ctx context.Context, projectID string, body types.ProjectRequest, params *types.RequestParameters) (*types.Response[types.ProjectResponse], error)
	Delete(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[any], error)
}

// ---- Adapter ----

// projectClientAdapter bridges the wrapper API (chainable, error-accumulating,
// wire-shape-hidden) to the low-level client (parameter-explicit, returning
// typed wire structs). Translates Project ↔ types.ProjectRequest/Response and
// surfaces HTTP errors as *aruba.HTTPError.
type projectClientAdapter struct {
	low projectLowLevelClient
}

var _ ProjectClient = (*projectClientAdapter)(nil)

func newProjectClientAdapter(rest *restclient.Client) *projectClientAdapter {
	return &projectClientAdapter{low: project.NewProjectsClientImpl(rest)}
}

// Create posts a new Project to the API and hydrates the wrapper from the response.
func (a *projectClientAdapter) Create(ctx context.Context, p *Project, opts ...CallOption) (*Project, error) {
	if err := p.Err(); err != nil {
		return p, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, p.toRequest(), rp)
	populateHTTPEnvelope(&p.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		p.fromResponse(resp.Data)
	}
	if err != nil {
		// low-level Create wraps *MetadataValidationError via fmt.Errorf("…: %w", err);
		// return the partial *Project so callers can inspect RawHTTP / RawError alongside
		// the typed error (contract preservation from internal/clients/project).
		return p, err
	}
	if resp != nil && !resp.IsSuccess() {
		return p, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return p, nil
}

// Get fetches a Project by Ref and returns a freshly hydrated wrapper.
func (a *projectClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Project, error) {
	id, err := projectIDFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, id, rp)
	out := &Project{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

// Update sends a PUT for the current wrapper state. Requires ID and parent.
func (a *projectClientAdapter) Update(ctx context.Context, p *Project, opts ...CallOption) (*Project, error) {
	if err := p.Err(); err != nil {
		return p, err
	}
	if p.ID() == "" {
		return p, fmt.Errorf("Update: project has no ID — call Get first or seed from Raw metadata")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, p.ID(), p.toRequest(), rp)
	populateHTTPEnvelope(&p.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		p.fromResponse(resp.Data)
	}
	if err != nil {
		return p, err
	}
	if resp != nil && !resp.IsSuccess() {
		return p, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return p, nil
}

// Delete removes the Project identified by Ref.
func (a *projectClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	id, err := projectIDFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, id, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

// List returns a paginated list of all Projects accessible to the caller.
func (a *projectClientAdapter) List(ctx context.Context, opts ...CallOption) (*List[*Project], error) {
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*Project
	if resp != nil && resp.Data != nil {
		items = make([]*Project, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			p := &Project{}
			p.fromResponse(&resp.Data.Values[i])
			items = append(items, p)
		}
	}
	// Pagination by raw URL is not yet wired into the low-level client.
	// The refetch stub matches the deferral pattern used by WaitUntilActive in #181.
	refetch := func(_ context.Context, _ string) (*List[*Project], error) {
		return nil, fmt.Errorf("List pagination by URL not yet wired; re-call List with adjusted CallOptions")
	}
	var total int64
	var self, prev, next, first, last string
	if resp != nil && resp.Data != nil {
		total = resp.Data.Total
		self = resp.Data.Self
		prev = resp.Data.Prev
		next = resp.Data.Next
		first = resp.Data.First
		last = resp.Data.Last
	}
	return newList(items, total, self, prev, next, first, last, resp, opts, refetch), nil
}
