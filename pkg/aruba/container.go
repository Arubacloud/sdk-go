package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/container"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

type ContainerClient interface {
	KaaS() KaaSClient
	ContainerRegistry() ContainerRegistryClient
}

type containerClientImpl struct {
	kaasClient              KaaSClient
	containerRegistryClient ContainerRegistryClient
}

// ContainerRegistry implements ContainerClient.
func (c *containerClientImpl) ContainerRegistry() ContainerRegistryClient {
	return c.containerRegistryClient
}

var _ ContainerClient = (*containerClientImpl)(nil)

func (c *containerClientImpl) KaaS() KaaSClient {
	return c.kaasClient
}

type KaaSClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*KaaS], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*KaaS, error)
	Create(ctx context.Context, k *KaaS, opts ...CallOption) (*KaaS, error)
	Update(ctx context.Context, k *KaaS, opts ...CallOption) (*KaaS, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type ContainerRegistryClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*ContainerRegistry], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*ContainerRegistry, error)
	Create(ctx context.Context, r *ContainerRegistry, opts ...CallOption) (*ContainerRegistry, error)
	Update(ctx context.Context, r *ContainerRegistry, opts ...CallOption) (*ContainerRegistry, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

// containerRegistryIDsFromRef extracts (projectID, registryID) from a Ref.
// Uses URI segment fallback on "registries" — no typed ancestor interface needed
// since ContainerRegistry has no descendant resource types.
func containerRegistryIDsFromRef(ref Ref) (projectID, registryID string, err error) {
	rid, ok := extractID(ref, func(r Ref) (string, bool) {
		return "", false // no typed interface — URI-only path
	}, "registries")
	if !ok || rid == "" {
		return "", "", fmt.Errorf("cannot determine registry ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, rid, nil
}

type containerRegistriesLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryList], error)
	Get(ctx context.Context, projectID, registryID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	Create(ctx context.Context, projectID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	Update(ctx context.Context, projectID, registryID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	Delete(ctx context.Context, projectID, registryID string, params *types.RequestParameters) (*types.Response[any], error)
}

type containerRegistriesClientAdapter struct {
	low containerRegistriesLowLevelClient
}

func newContainerRegistriesClientAdapter(rest *restclient.Client) *containerRegistriesClientAdapter {
	if rest == nil {
		return &containerRegistriesClientAdapter{}
	}
	return &containerRegistriesClientAdapter{low: container.NewContainerRegistryClientImpl(rest)}
}

func (a *containerRegistriesClientAdapter) Create(ctx context.Context, r *ContainerRegistry, opts ...CallOption) (*ContainerRegistry, error) {
	if err := r.Err(); err != nil {
		return r, err
	}
	if r.ProjectID() == "" {
		return r, fmt.Errorf("Create: ContainerRegistry has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, r.ProjectID(), r.toRequest(), rp)
	populateHTTPEnvelope(&r.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		r.fromResponse(resp.Data)
	}
	if err != nil {
		return r, err
	}
	if resp != nil && !resp.IsSuccess() {
		return r, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return r, nil
}

func (a *containerRegistriesClientAdapter) Update(ctx context.Context, r *ContainerRegistry, opts ...CallOption) (*ContainerRegistry, error) {
	if err := r.Err(); err != nil {
		return r, err
	}
	if r.ContainerRegistryID() == "" {
		return r, fmt.Errorf("Update: ContainerRegistry has no ID")
	}
	if r.ProjectID() == "" {
		return r, fmt.Errorf("Update: ContainerRegistry has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, r.ProjectID(), r.ContainerRegistryID(), r.toRequest(), rp)
	populateHTTPEnvelope(&r.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		r.fromResponse(resp.Data)
	}
	if err != nil {
		return r, err
	}
	if resp != nil && !resp.IsSuccess() {
		return r, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return r, nil
}

func (a *containerRegistriesClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*ContainerRegistry, error) {
	projectID, registryID, err := containerRegistryIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, registryID, rp)
	out := &ContainerRegistry{}
	out.projectID = projectID
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if out.projectID == "" {
		out.projectID = projectID
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *containerRegistriesClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, registryID, err := containerRegistryIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, registryID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *containerRegistriesClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*ContainerRegistry], error) {
	projectID, err := projectIDFromRef(parent)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*ContainerRegistry
	if resp != nil && resp.Data != nil {
		items = make([]*ContainerRegistry, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			cr := &ContainerRegistry{}
			cr.projectID = projectID
			cr.fromResponse(&resp.Data.Values[i])
			if cr.projectID == "" {
				cr.projectID = projectID
			}
			items = append(items, cr)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*ContainerRegistry], error) {
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
