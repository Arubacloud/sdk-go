package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// SecurityGroup is the wrapper for an Aruba Cloud Security Group (a direct child of a VPC).
// Construct with aruba.NewSecurityGroup() and bind it to a parent VPC via IntoVPC(vpc).
type SecurityGroup struct {
	errMixin
	metadataMixin
	vpcScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	defaultSG *bool                        // Properties.Default (request: *bool for omitempty; response: plain bool)
	response  *types.SecurityGroupResponse // backs Raw()
}

// Setters (chainable; promoted methods re-exposed at *SecurityGroup level).
func (sg *SecurityGroup) IntoVPC(v Ref) *SecurityGroup            { sg.intoVPC(v); return sg }
func (sg *SecurityGroup) WithName(n string) *SecurityGroup        { sg.withName(n); return sg }
func (sg *SecurityGroup) AddTag(t string) *SecurityGroup          { sg.addTag(t); return sg }
func (sg *SecurityGroup) RemoveTag(t string) *SecurityGroup       { sg.removeTag(t); return sg }
func (sg *SecurityGroup) ReplaceTags(ts ...string) *SecurityGroup { sg.replaceTags(ts...); return sg }
func (sg *SecurityGroup) WithDefault(b bool) *SecurityGroup       { sg.defaultSG = &b; return sg }

// URI satisfies Ref.
func (sg *SecurityGroup) URI() string { return sg.RespURI() }

// SecurityGroupID satisfies withSecurityGroupID so child wrappers (SecurityGroupRule)
// can extract this ID via typed assertion.
func (sg *SecurityGroup) SecurityGroupID() string { return sg.ID() }

// Raw shadows responseMetadataMixin.Raw() with the full SecurityGroup response.
func (sg *SecurityGroup) Raw() *types.SecurityGroupResponse { return sg.response }

// RawRequest returns what toRequest() would emit right now.
func (sg *SecurityGroup) RawRequest() types.SecurityGroupRequest { return sg.toRequest() }

// Default returns the security group's default flag, or false if unset.
func (sg *SecurityGroup) Default() bool {
	if sg.defaultSG == nil {
		return false
	}
	return *sg.defaultSG
}

func (sg *SecurityGroup) toRequest() types.SecurityGroupRequest {
	return types.SecurityGroupRequest{
		Metadata: sg.toMetadata(),
		Properties: types.SecurityGroupPropertiesRequest{
			Default: sg.defaultSG,
		},
	}
}

func (sg *SecurityGroup) fromResponse(resp *types.SecurityGroupResponse) {
	if resp == nil {
		return
	}
	sg.response = resp
	sg.setMeta(&resp.Metadata)
	sg.withName(securityGroupDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		sg.replaceTags(resp.Metadata.Tags...)
	}
	sg.setStatus(&resp.Status)
	sg.setLinked(resp.Properties.LinkedResources)

	// Properties.Default is plain bool on the response — backfill into our *bool field.
	d := resp.Properties.Default
	sg.defaultSG = &d

	// Backfill ancestor IDs: prefer ProjectResponseMetadata, then URI parse.
	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		sg.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if (sg.vpcID == "" || sg.projectID == "") && sg.RespURI() != "" {
		ids := parseURIIDs(sg.RespURI())
		if sg.vpcID == "" {
			sg.vpcID = ids["vpcs"]
		}
		if sg.projectID == "" {
			sg.projectID = ids["projects"]
		}
	}
}

func securityGroupDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type securityGroupLowLevelClient interface {
	List(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.SecurityGroupList], error)
	Get(ctx context.Context, projectID, vpcID, securityGroupID string, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Create(ctx context.Context, projectID, vpcID string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Update(ctx context.Context, projectID, vpcID, securityGroupID string, body types.SecurityGroupRequest, params *types.RequestParameters) (*types.Response[types.SecurityGroupResponse], error)
	Delete(ctx context.Context, projectID, vpcID, securityGroupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type securityGroupsClientAdapter struct{ low securityGroupLowLevelClient }

func newSecurityGroupsClientAdapter(rest *restclient.Client) *securityGroupsClientAdapter {
	if rest == nil {
		return &securityGroupsClientAdapter{}
	}
	return &securityGroupsClientAdapter{
		low: network.NewSecurityGroupsClientImpl(rest, network.NewVPCsClientImpl(rest)),
	}
}

func (a *securityGroupsClientAdapter) Create(ctx context.Context, sg *SecurityGroup, opts ...CallOption) (*SecurityGroup, error) {
	if err := sg.Err(); err != nil {
		return sg, err
	}
	if sg.VPCID() == "" || sg.ProjectID() == "" {
		return sg, fmt.Errorf("Create: security group has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, sg.ProjectID(), sg.VPCID(), sg.toRequest(), rp)
	populateHTTPEnvelope(&sg.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		sg.fromResponse(resp.Data)
	}
	if err != nil {
		return sg, err
	}
	if resp != nil && !resp.IsSuccess() {
		return sg, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return sg, nil
}

func (a *securityGroupsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*SecurityGroup, error) {
	projectID, vpcID, securityGroupID, err := securityGroupIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, securityGroupID, rp)
	out := &SecurityGroup{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
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

func (a *securityGroupsClientAdapter) Update(ctx context.Context, sg *SecurityGroup, opts ...CallOption) (*SecurityGroup, error) {
	if err := sg.Err(); err != nil {
		return sg, err
	}
	if sg.ID() == "" {
		return sg, fmt.Errorf("Update: security group has no ID — call Get first or seed from response metadata")
	}
	if sg.VPCID() == "" || sg.ProjectID() == "" {
		return sg, fmt.Errorf("Update: security group has no VPC — call IntoVPC first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, sg.ProjectID(), sg.VPCID(), sg.ID(), sg.toRequest(), rp)
	populateHTTPEnvelope(&sg.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		sg.fromResponse(resp.Data)
	}
	if err != nil {
		return sg, err
	}
	if resp != nil && !resp.IsSuccess() {
		return sg, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return sg, nil
}

func (a *securityGroupsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, securityGroupID, err := securityGroupIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, securityGroupID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *securityGroupsClientAdapter) List(ctx context.Context, vpc Ref, opts ...CallOption) (*List[*SecurityGroup], error) {
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
	var items []*SecurityGroup
	if resp != nil && resp.Data != nil {
		items = make([]*SecurityGroup, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			sg := &SecurityGroup{}
			sg.fromResponse(&resp.Data.Values[i])
			if sg.vpcID == "" {
				sg.vpcID = vpcID
			}
			if sg.projectID == "" {
				sg.projectID = projectID
			}
			items = append(items, sg)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*SecurityGroup], error) {
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

// securityGroupIDsFromRef extracts (projectID, vpcID, securityGroupID) from a Ref.
// Tries typed assertions first, then falls back to URI path parsing.
func securityGroupIDsFromRef(ref Ref) (projectID, vpcID, securityGroupID string, err error) {
	sgid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withSecurityGroupID); ok {
			return w.SecurityGroupID(), true
		}
		return "", false
	}, "security-groups")
	if !ok || sgid == "" {
		return "", "", "", fmt.Errorf("cannot determine security group ID from Ref %q", ref.URI())
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
	return pid, vid, sgid, nil
}
