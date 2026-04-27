package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/project"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ProjectClient is the wrapper-based public surface for project CRUD.
type ProjectClient interface {
	Create(ctx context.Context, p *Project, opts ...CallOption) (*Project, error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*Project, error)
	Update(ctx context.Context, p *Project, opts ...CallOption) (*Project, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
	List(ctx context.Context, opts ...CallOption) (*List[*Project], error)
}

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

// projectClientAdapter bridges the wrapper interface to the existing raw-types impl.
type projectClientAdapter struct {
	low projectLowLevelClient
}

func newProjectClientAdapter(rest *restclient.Client) *projectClientAdapter {
	return &projectClientAdapter{low: project.NewProjectsClientImpl(rest)}
}

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

// projectIDFromRef extracts a project ID from a Ref, preferring the typed
// withProjectID assertion and falling back to URI path parsing.
func projectIDFromRef(ref Ref) (string, error) {
	id, ok := extractID(ref, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || id == "" {
		return "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return id, nil
}
