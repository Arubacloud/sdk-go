package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// LoadBalancer is the wrapper for an Aruba Cloud Load Balancer (a direct child of a Project).
// LoadBalancer is read-only: instances are obtained via Client.FromNetwork().LoadBalancers().Get/List.
// There is no NewLoadBalancer() factory — the resource cannot be created or mutated through the SDK.
type LoadBalancer struct {
	metadataMixin         // Name(), Tags() — populated from response metadata
	regionalMixin         // Region() — populated from response location
	projectScopedMixin    // ProjectID() — back-filled from response/URI; intoProject() is unexported and unused
	responseMetadataMixin // ID(), RespURI(), CreatedAt(), UpdatedAt(), Version()
	statusMixin           // State(), IsDisabled(), FailureReason(), DisableReasons(), PreviousState()
	linkedMixin           // LinkedResources()
	httpEnvelopeMixin     // RawHTTP(), StatusCode(), Headers(), RawError()

	address  *string                     // Properties.Address (read-only from response)
	vpc      *types.ReferenceResource    // Properties.VPC (linked VPC reference)
	response *types.LoadBalancerResponse // backs Raw()
}

// URI satisfies Ref.
func (l *LoadBalancer) URI() string { return l.RespURI() }

// LoadBalancerID satisfies withLoadBalancerID so adapters can extract this ID typed.
func (l *LoadBalancer) LoadBalancerID() string { return l.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full LoadBalancer response.
func (l *LoadBalancer) Raw() *types.LoadBalancerResponse { return l.response }

// Address returns the public IP address assigned to this Load Balancer, or "" if absent.
func (l *LoadBalancer) Address() string {
	if l.address == nil {
		return ""
	}
	return *l.address
}

// VPC returns the linked VPC reference URI, or "" if the Load Balancer is not VPC-attached.
func (l *LoadBalancer) VPC() string {
	if l.vpc == nil {
		return ""
	}
	return l.vpc.URI
}

func (l *LoadBalancer) fromResponse(resp *types.LoadBalancerResponse) {
	if resp == nil {
		return
	}
	l.response = resp
	l.setMeta(&resp.Metadata)
	l.withName(loadBalancerDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		l.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		l.withLocation(resp.Metadata.LocationResponse.Value)
	}
	l.setStatus(&resp.Status)
	l.setTerminalStates(loadBalancerTerminalStates)
	l.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.Address != nil && *resp.Properties.Address != "" {
		addr := *resp.Properties.Address
		l.address = &addr
	}
	if resp.Properties.VPC != nil {
		v := *resp.Properties.VPC
		l.vpc = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		l.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if l.projectID == "" && l.RespURI() != "" {
		if pid := parseURIIDs(l.RespURI())["projects"]; pid != "" {
			l.projectID = pid
		}
	}
}

func loadBalancerDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var loadBalancerTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type loadBalancerLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.LoadBalancerList], error)
	Get(ctx context.Context, projectID, loadBalancerID string, params *types.RequestParameters) (*types.Response[types.LoadBalancerResponse], error)
}

type loadBalancersClientAdapter struct{ low loadBalancerLowLevelClient }

func newLoadBalancersClientAdapter(rest *restclient.Client) *loadBalancersClientAdapter {
	if rest == nil {
		return &loadBalancersClientAdapter{}
	}
	return &loadBalancersClientAdapter{low: network.NewLoadBalancersClientImpl(rest)}
}

func (a *loadBalancersClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*LoadBalancer, error) {
	projectID, loadBalancerID, err := loadBalancerIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, loadBalancerID, rp)
	out := &LoadBalancer{}
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

func (a *loadBalancersClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*LoadBalancer], error) {
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
	var items []*LoadBalancer
	if resp != nil && resp.Data != nil {
		items = make([]*LoadBalancer, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			lb := &LoadBalancer{}
			lb.fromResponse(&resp.Data.Values[i])
			lb.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, lb)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					lb.fromResponse(fresh.Raw())
				}
				return nil
			})
			if lb.projectID == "" {
				lb.projectID = projectID
			}
			items = append(items, lb)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*LoadBalancer], error) {
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

// loadBalancerIDsFromRef extracts (projectID, loadBalancerID) from a Ref. Tries typed
// assertions first, then falls back to URI path parsing.
func loadBalancerIDsFromRef(ref Ref) (projectID, loadBalancerID string, err error) {
	lid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withLoadBalancerID); ok {
			return w.LoadBalancerID(), true
		}
		return "", false
	}, "loadbalancers")
	if !ok || lid == "" {
		return "", "", fmt.Errorf("cannot determine load balancer ID from Ref %q", ref.URI())
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
	return pid, lid, nil
}
