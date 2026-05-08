package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ElasticIP is the wrapper for an Aruba Cloud Elastic IP (a direct child of a Project).
// Construct with aruba.NewElasticIP() and bind it to a parent project via IntoProject(project).
type ElasticIP struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	billingPeriod *string                  // Properties.BillingPlan.BillingPeriod
	address       *string                  // Properties.Address (read-only from response)
	response      *types.ElasticIPResponse // backs Raw()
}

func (e *ElasticIP) IntoProject(p Ref) *ElasticIP          { e.intoProject(p); return e }
func (e *ElasticIP) WithName(n string) *ElasticIP          { e.withName(n); return e }
func (e *ElasticIP) AddTag(t string) *ElasticIP            { e.addTag(t); return e }
func (e *ElasticIP) RemoveTag(t string) *ElasticIP         { e.removeTag(t); return e }
func (e *ElasticIP) ReplaceTags(ts ...string) *ElasticIP   { e.replaceTags(ts...); return e }
func (e *ElasticIP) WithLocation(loc string) *ElasticIP    { e.withLocation(loc); return e }
func (e *ElasticIP) InRegion(region string) *ElasticIP     { e.withLocation(region); return e }
func (e *ElasticIP) WithBillingPeriod(p string) *ElasticIP { e.billingPeriod = &p; return e }

// URI satisfies Ref.
func (e *ElasticIP) URI() string { return e.RespURI() }

// ElasticIPID satisfies withElasticIPID so adapters can extract this ID typed.
func (e *ElasticIP) ElasticIPID() string { return e.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full ElasticIP response.
func (e *ElasticIP) Raw() *types.ElasticIPResponse { return e.response }

// RawRequest returns what toRequest() would emit right now.
func (e *ElasticIP) RawRequest() types.ElasticIPRequest { return e.toRequest() }

func (e *ElasticIP) BillingPeriod() string {
	if e.billingPeriod == nil {
		return ""
	}
	return *e.billingPeriod
}

func (e *ElasticIP) Address() string {
	if e.address == nil {
		return ""
	}
	return *e.address
}

func (e *ElasticIP) toRequest() types.ElasticIPRequest {
	var bp string
	if e.billingPeriod != nil {
		bp = *e.billingPeriod
	}
	return types.ElasticIPRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: e.toMetadata(),
			Location:                e.toLocation(),
		},
		Properties: types.ElasticIPPropertiesRequest{
			BillingPlan: types.BillingPeriodResource{BillingPeriod: bp},
		},
	}
}

func (e *ElasticIP) fromResponse(resp *types.ElasticIPResponse) {
	if resp == nil {
		return
	}
	e.response = resp
	e.setMeta(&resp.Metadata)
	e.withName(elasticIPDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		e.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		e.withLocation(resp.Metadata.LocationResponse.Value)
	}
	e.setStatus(&resp.Status)
	e.setTerminalStates(elasticIPTerminalStates)
	e.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.BillingPlan.BillingPeriod != "" {
		bp := resp.Properties.BillingPlan.BillingPeriod
		e.billingPeriod = &bp
	}
	if resp.Properties.Address != nil && *resp.Properties.Address != "" {
		addr := *resp.Properties.Address
		e.address = &addr
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		e.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if e.projectID == "" && e.RespURI() != "" {
		if pid := parseURIIDs(e.RespURI())["projects"]; pid != "" {
			e.projectID = pid
		}
	}
}

func elasticIPDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var elasticIPTerminalStates = map[string]bool{
	"NotUsed": true,
	"InUse":   true,
	"Error":   false,
	"Failed": false,
}

// WaitUntilNotUsed blocks until the ElasticIP reaches the "NotUsed" state —
// the steady terminal state for an unattached EIP. Call this after Create and
// before passing the EIP to a CloudServer, ContainerRegistry, or LoadBalancer.
func (e *ElasticIP) WaitUntilNotUsed(ctx context.Context, opts ...WaitOption) error {
	return e.WaitUntilStates(ctx, []string{"NotUsed"}, opts...)
}

// WaitUntilUsed blocks until the ElasticIP reaches the "InUse" or "Used"
// state — both signal that the EIP has been bound to a consumer resource. The
// platform may emit either value; this method succeeds on whichever arrives.
func (e *ElasticIP) WaitUntilUsed(ctx context.Context, opts ...WaitOption) error {
	return e.WaitUntilStates(ctx, []string{"InUse", "Used"}, opts...)
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type elasticIPLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ElasticList], error)
	Get(ctx context.Context, projectID, elasticIPID string, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Create(ctx context.Context, projectID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Update(ctx context.Context, projectID, elasticIPID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Delete(ctx context.Context, projectID, elasticIPID string, params *types.RequestParameters) (*types.Response[any], error)
}

type elasticIPsClientAdapter struct{ low elasticIPLowLevelClient }

func newElasticIPsClientAdapter(rest *restclient.Client) *elasticIPsClientAdapter {
	if rest == nil {
		return &elasticIPsClientAdapter{}
	}
	return &elasticIPsClientAdapter{low: network.NewElasticIPsClientImpl(rest)}
}

func (a *elasticIPsClientAdapter) Create(ctx context.Context, e *ElasticIP, opts ...CallOption) (*ElasticIP, error) {
	if err := e.Err(); err != nil {
		return e, err
	}
	if e.ProjectID() == "" {
		return e, fmt.Errorf("Create: elastic IP has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, e.ProjectID(), e.toRequest(), rp)
	populateHTTPEnvelope(&e.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		e.fromResponse(resp.Data)
		e.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, e)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				e.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return e, err
	}
	if resp != nil && !resp.IsSuccess() {
		return e, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return e, nil
}

func (a *elasticIPsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*ElasticIP, error) {
	projectID, elasticIPID, err := elasticIPIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, elasticIPID, rp)
	out := &ElasticIP{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
		out.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, out)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				out.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	out.projectID = projectID
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *elasticIPsClientAdapter) Update(ctx context.Context, e *ElasticIP, opts ...CallOption) (*ElasticIP, error) {
	if err := e.Err(); err != nil {
		return e, err
	}
	if e.ID() == "" {
		return e, fmt.Errorf("Update: elastic IP has no ID — call Get first or seed from response metadata")
	}
	if e.ProjectID() == "" {
		return e, fmt.Errorf("Update: elastic IP has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, e.ProjectID(), e.ID(), e.toRequest(), rp)
	populateHTTPEnvelope(&e.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		e.fromResponse(resp.Data)
		e.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, e)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				e.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return e, err
	}
	if resp != nil && !resp.IsSuccess() {
		return e, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return e, nil
}

func (a *elasticIPsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, elasticIPID, err := elasticIPIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, elasticIPID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *elasticIPsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*ElasticIP], error) {
	projectID, err := projectIDFromRef(project)
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
	var items []*ElasticIP
	if resp != nil && resp.Data != nil {
		items = make([]*ElasticIP, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			e := &ElasticIP{}
			e.fromResponse(&resp.Data.Values[i])
			e.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, e)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					e.fromResponse(fresh.Raw())
				}
				return nil
			})
			if e.projectID == "" {
				e.projectID = projectID
			}
			items = append(items, e)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*ElasticIP], error) {
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

// elasticIPIDsFromRef extracts (projectID, elasticIPID) from a Ref. Tries typed
// assertions first, then falls back to URI path parsing.
func elasticIPIDsFromRef(ref Ref) (projectID, elasticIPID string, err error) {
	eid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withElasticIPID); ok {
			return w.ElasticIPID(), true
		}
		return "", false
	}, "elasticIps")
	if !ok || eid == "" {
		return "", "", fmt.Errorf("cannot determine elastic IP ID from Ref %q", ref.URI())
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
	return pid, eid, nil
}
