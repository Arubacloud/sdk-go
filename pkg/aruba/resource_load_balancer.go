package aruba

import "github.com/Arubacloud/sdk-go/pkg/types"

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
