package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ElasticIPRef returns a Ref that points to the ElasticIP with the given IDs.
func ElasticIPRef(projectID, eipID string) Ref {
	return URI(fmt.Sprintf("/projects/%s/providers/Aruba.Network/elasticIps/%s", projectID, eipID))
}

// ---- Wrapper ----

// ElasticIP is the wrapper for an Aruba Cloud Elastic IP (a child of a Project).
// Construct with aruba.NewElasticIP() and bind it via IntoProject(project).
//
// Wraps types.ElasticIPResponse / types.ElasticIPRequest. The wrapper carries
// pointer-typed private fields so unset values round-trip through
// the JSON layer correctly.
type ElasticIP struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	billingPeriod *BillingPeriod           // Properties.BillingPlan.BillingPeriod
	address       *string                  // Properties.Address (read-only from response)
	response      *types.ElasticIPResponse // backs Raw()
}

// Setters — chainable, general → specific

// IntoProject binds this ElasticIP to its parent project. Required before Create.
func (e *ElasticIP) IntoProject(p Ref) *ElasticIP { e.intoProject(p); return e }

// Named sets the resource name. Required by the API.
func (e *ElasticIP) Named(n string) *ElasticIP { e.named(n); return e }

// AddTag appends a tag for filtering and accounting.
func (e *ElasticIP) AddTag(t string) *ElasticIP { e.addTag(t); return e }

// RemoveTag removes a previously-added tag. No-op if absent.
func (e *ElasticIP) RemoveTag(t string) *ElasticIP { e.removeTag(t); return e }

// ReplaceTags replaces the entire tag set with the given values.
func (e *ElasticIP) ReplaceTags(ts ...string) *ElasticIP { e.replaceTags(ts...); return e }

// InRegion sets the region for this resource.
func (e *ElasticIP) InRegion(region Region) *ElasticIP { e.inRegion(region); return e }

// WithBillingPeriod sets the billing period. Defaults to hourly when unset.
func (e *ElasticIP) WithBillingPeriod(p BillingPeriod) *ElasticIP { e.billingPeriod = &p; return e }

// Getters — general → specific

// URI satisfies Ref.
func (e *ElasticIP) URI() string { return e.RespURI() }

// ElasticIPID satisfies withElasticIPID so adapters can extract this ID typed.
func (e *ElasticIP) ElasticIPID() string { return e.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full ElasticIP response.
func (e *ElasticIP) Raw() *types.ElasticIPResponse { return e.response }

// RawRequest returns what toRequest() would emit right now.
func (e *ElasticIP) RawRequest() types.ElasticIPRequest { return e.toRequest() }

// BillingPeriod returns the configured billing period ("" if unset).
func (e *ElasticIP) BillingPeriod() BillingPeriod {
	if e.billingPeriod == nil {
		return ""
	}
	return *e.billingPeriod
}

// Address returns the server-assigned public IP address ("" if unassigned).
func (e *ElasticIP) Address() string {
	if e.address == nil {
		return ""
	}
	return *e.address
}

// Wire converters

// toRequest assembles the Create/Update body from current setter state. Defaults are applied at the wire boundary.
func (e *ElasticIP) toRequest() types.ElasticIPRequest {
	return types.ElasticIPRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: e.toMetadata(),
			Location:                e.toLocation(),
		},
		Properties: types.ElasticIPPropertiesRequest{
			BillingPeriod: elasticIPBillingPeriodWire().Out(defaultBillingPeriod(e.billingPeriod)),
		},
	}
}

// fromResponse hydrates the wrapper from a server reply. Nil-safe.
func (e *ElasticIP) fromResponse(resp *types.ElasticIPResponse) {
	if resp == nil {
		return
	}
	e.response = resp
	e.setMeta(&resp.Metadata)
	e.named(elasticIPDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		e.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		e.inRegion(resp.Metadata.LocationResponse.Value)
	}
	e.setStatus(&resp.Status)
	e.setTerminalStates(elasticIPTerminalStates)
	e.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.BillingPeriod != nil {
		e.billingPeriod = elasticIPBillingPeriodWire().In(resp.Properties.BillingPeriod)
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

// elasticIPBillingPeriodWire returns a translator for ElasticIP's lowercase
// wire billing-period values (e.g. "hourly") vs. the standard SDK constants.
func elasticIPBillingPeriodWire() *billingPeriodTranslator {
	return newBillingPeriodTranslator(map[BillingPeriod]string{
		BillingPeriodHour:  "hourly",
		BillingPeriodMonth: "monthly",
		BillingPeriodYear:  "yearly",
	})
}

var elasticIPTerminalStates = map[string]bool{
	"NotUsed": true,
	"InUse":   true,
	"Error":   false,
	"Failed":  false,
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

// ---- Low-level client interface ----

// elasticIPLowLevelClient is the contract the wrapper depends on. Returning
// *types.Response[T] preserves HTTP envelope details (status code, headers,
// raw body) for the wrapper's diagnostics.
type elasticIPLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ElasticList], error)
	Get(ctx context.Context, projectID, elasticIPID string, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Create(ctx context.Context, projectID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Update(ctx context.Context, projectID, elasticIPID string, body types.ElasticIPRequest, params *types.RequestParameters) (*types.Response[types.ElasticIPResponse], error)
	Delete(ctx context.Context, projectID, elasticIPID string, params *types.RequestParameters) (*types.Response[any], error)
}

// ---- Adapter ----

// elasticIPsClientAdapter bridges the wrapper API (chainable, error-accumulating,
// wire-shape-hidden) to the low-level client (parameter-explicit, returning
// typed wire structs). Translates ElasticIP ↔ types.ElasticIPRequest/Response and
// surfaces HTTP errors as *aruba.HTTPError.
type elasticIPsClientAdapter struct {
	low  elasticIPLowLevelClient
	rest *restclient.Client
}

var _ ElasticIPsClient = (*elasticIPsClientAdapter)(nil)

func newElasticIPsClientAdapter(rest *restclient.Client) *elasticIPsClientAdapter {
	if rest == nil {
		return &elasticIPsClientAdapter{}
	}
	return &elasticIPsClientAdapter{low: network.NewElasticIPsClientImpl(rest), rest: rest}
}

// Create posts a new ElasticIP to the API and hydrates the wrapper from the response.
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

// Get fetches an ElasticIP by Ref and returns a freshly hydrated wrapper.
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

// Update sends a PUT for the current wrapper state. Requires ID and parent.
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

// Delete removes the ElasticIP identified by Ref.
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

// List returns a paginated list of ElasticIP in the given parent scope.
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
	var refetch func(ctx context.Context, pageURL string) (*List[*ElasticIP], error)
	refetch = func(ctx context.Context, pageURL string) (*List[*ElasticIP], error) {
		fetch := listPageFetch[types.ElasticList](a.rest, opts)
		pageResp, fetchErr := fetch(ctx, pageURL)
		if fetchErr != nil {
			return nil, fetchErr
		}
		if pageResp != nil && !pageResp.IsSuccess() {
			return nil, &HTTPError{StatusCode: pageResp.StatusCode, Body: pageResp.RawBody, ErrResp: pageResp.Error}
		}
		var pageItems []*ElasticIP
		if pageResp != nil && pageResp.Data != nil {
			pageItems = make([]*ElasticIP, 0, len(pageResp.Data.Values))
			for i := range pageResp.Data.Values {
				e := &ElasticIP{}
				e.fromResponse(&pageResp.Data.Values[i])
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
				pageItems = append(pageItems, e)
			}
		}
		var total2 int64
		var self2, prev2, next2, first2, last2 string
		if pageResp != nil && pageResp.Data != nil {
			total2 = pageResp.Data.Total
			self2 = pageResp.Data.Self
			prev2 = pageResp.Data.Prev
			next2 = pageResp.Data.Next
			first2 = pageResp.Data.First
			last2 = pageResp.Data.Last
		}
		return newList(pageItems, total2, self2, prev2, next2, first2, last2, pageResp, opts, refetch), nil
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
