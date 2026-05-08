package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// Subnet is the wrapper for an Aruba Cloud subnet (a child of a VPC).
// Construct with aruba.NewSubnet() and bind it to a parent VPC via IntoVPC(vpc).
type Subnet struct {
	errMixin
	metadataMixin
	regionalMixin
	vpcScopedMixin // direct parent; populates vpcID + projectID
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	subnetType    *string               // Properties.Type ("Basic" / "Advanced")
	defaultSubnet *bool                 // Properties.Default
	cidr          *string               // Properties.Network.Address
	dhcp          *SubnetDHCP           // Properties.DHCP (sub-builder)
	response      *types.SubnetResponse // backs Raw()
}

func (s *Subnet) IntoVPC(v Ref) *Subnet            { s.intoVPC(v); return s }
func (s *Subnet) WithName(n string) *Subnet        { s.withName(n); return s }
func (s *Subnet) AddTag(t string) *Subnet          { s.addTag(t); return s }
func (s *Subnet) RemoveTag(t string) *Subnet       { s.removeTag(t); return s }
func (s *Subnet) ReplaceTags(ts ...string) *Subnet { s.replaceTags(ts...); return s }
func (s *Subnet) WithLocation(loc string) *Subnet  { s.withLocation(loc); return s }
func (s *Subnet) InRegion(region string) *Subnet   { s.withLocation(region); return s }
func (s *Subnet) WithType(t string) *Subnet        { s.subnetType = &t; return s }
func (s *Subnet) WithDefault(b bool) *Subnet       { s.defaultSubnet = &b; return s }
func (s *Subnet) WithCIDR(cidr string) *Subnet     { s.cidr = &cidr; return s }
func (s *Subnet) WithDHCP(d *SubnetDHCP) *Subnet   { s.dhcp = d; return s }

// URI satisfies Ref.
func (s *Subnet) URI() string { return s.RespURI() }

// SubnetID satisfies withSubnetID so future grandchildren can extract this ID.
func (s *Subnet) SubnetID() string { return s.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full subnet response.
func (s *Subnet) Raw() *types.SubnetResponse { return s.response }

// RawRequest returns what toRequest() would emit right now.
func (s *Subnet) RawRequest() types.SubnetRequest { return s.toRequest() }

func (s *Subnet) Type() string {
	if s.subnetType == nil {
		return ""
	}
	return *s.subnetType
}
func (s *Subnet) IsDefault() bool {
	if s.defaultSubnet == nil {
		return false
	}
	return *s.defaultSubnet
}
func (s *Subnet) CIDR() string {
	if s.cidr == nil {
		return ""
	}
	return *s.cidr
}
func (s *Subnet) DHCP() *SubnetDHCP { return s.dhcp }

func (s *Subnet) toRequest() types.SubnetRequest {
	props := types.SubnetPropertiesRequest{}
	if s.subnetType != nil {
		props.Type = types.SubnetType(*s.subnetType)
	}
	if s.defaultSubnet != nil {
		props.Default = *s.defaultSubnet
	}
	if s.cidr != nil {
		props.Network = &types.SubnetNetwork{Address: *s.cidr}
	}
	if s.dhcp != nil {
		props.DHCP = s.dhcp.toType()
	}
	return types.SubnetRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: s.toMetadata(),
			Location:                s.toLocation(),
		},
		Properties: props,
	}
}

func (s *Subnet) fromResponse(resp *types.SubnetResponse) {
	if resp == nil {
		return
	}
	s.response = resp
	s.setMeta(&resp.Metadata)
	s.withName(subnetDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		s.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		s.withLocation(resp.Metadata.LocationResponse.Value)
	}
	s.setStatus(&resp.Status)
	s.setTerminalStates(subnetTerminalStates)
	s.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.Type != "" {
		t := string(resp.Properties.Type)
		s.subnetType = &t
	}
	d := resp.Properties.Default
	s.defaultSubnet = &d
	if resp.Properties.Network != nil && resp.Properties.Network.Address != "" {
		addr := resp.Properties.Network.Address
		s.cidr = &addr
	}
	if resp.Properties.DHCP != nil {
		s.dhcp = dhcpFromType(resp.Properties.DHCP)
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		s.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	// Backfill ancestor IDs from response URI if not already set.
	if (s.vpcID == "" || s.projectID == "") && s.RespURI() != "" {
		ids := parseURIIDs(s.RespURI())
		if s.vpcID == "" {
			s.vpcID = ids["vpcs"]
		}
		if s.projectID == "" {
			s.projectID = ids["projects"]
		}
	}
}

func subnetDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var subnetTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type subnetLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.SubnetList], error)
	Get(ctx context.Context, projectID, vpcID, subnetID string, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Create(ctx context.Context, projectID, vpcID string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Update(ctx context.Context, projectID, vpcID, subnetID string, body types.SubnetRequest, params *types.RequestParameters) (*types.Response[types.SubnetResponse], error)
	Delete(ctx context.Context, projectID, vpcID, subnetID string, params *types.RequestParameters) (*types.Response[any], error)
}

type subnetsClientAdapter struct{ low subnetLowLevelClient }

func newSubnetsClientAdapter(rest *restclient.Client) *subnetsClientAdapter {
	if rest == nil {
		return &subnetsClientAdapter{}
	}
	return &subnetsClientAdapter{
		low: network.NewSubnetsClientImpl(rest, network.NewVPCsClientImpl(rest)),
	}
}

func (a *subnetsClientAdapter) Create(ctx context.Context, s *Subnet, opts ...CallOption) (*Subnet, error) {
	if err := s.Err(); err != nil {
		return s, err
	}
	if s.VPCID() == "" || s.ProjectID() == "" {
		return s, fmt.Errorf("Create: subnet has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, s.ProjectID(), s.VPCID(), s.toRequest(), rp)
	populateHTTPEnvelope(&s.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		s.fromResponse(resp.Data)
		s.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, s)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				s.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return s, err
	}
	if resp != nil && !resp.IsSuccess() {
		return s, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return s, nil
}

func (a *subnetsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Subnet, error) {
	projectID, vpcID, subnetID, err := subnetIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, subnetID, rp)
	out := &Subnet{}
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
	out.vpcID = vpcID
	out.projectID = projectID
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *subnetsClientAdapter) Update(ctx context.Context, s *Subnet, opts ...CallOption) (*Subnet, error) {
	if err := s.Err(); err != nil {
		return s, err
	}
	if s.ID() == "" {
		return s, fmt.Errorf("Update: subnet has no ID — call Get first or seed from response metadata")
	}
	if s.VPCID() == "" || s.ProjectID() == "" {
		return s, fmt.Errorf("Update: subnet has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, s.ProjectID(), s.VPCID(), s.ID(), s.toRequest(), rp)
	populateHTTPEnvelope(&s.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		s.fromResponse(resp.Data)
		s.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, s)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				s.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return s, err
	}
	if resp != nil && !resp.IsSuccess() {
		return s, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return s, nil
}

func (a *subnetsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, subnetID, err := subnetIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, subnetID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *subnetsClientAdapter) List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*Subnet], error) {
	projectID, vpcID, err := vpcIDsFromRef(vpc)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, vpcID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*Subnet
	if resp != nil && resp.Data != nil {
		items = make([]*Subnet, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			s := &Subnet{}
			s.fromResponse(&resp.Data.Values[i])
			s.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, s)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					s.fromResponse(fresh.Raw())
				}
				return nil
			})
			s.vpcID = vpcID
			if s.projectID == "" {
				s.projectID = projectID
			}
			items = append(items, s)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*Subnet], error) {
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

// subnetIDsFromRef extracts (projectID, vpcID, subnetID) from a Ref. Tries typed
// assertions first, then falls back to URI path parsing.
func subnetIDsFromRef(ref Ref) (projectID, vpcID, subnetID string, err error) {
	sid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withSubnetID); ok {
			return w.SubnetID(), true
		}
		return "", false
	}, "subnets")
	if !ok || sid == "" {
		return "", "", "", fmt.Errorf("cannot determine subnet ID from Ref %q", ref.URI())
	}
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, vid, sid, nil
}
