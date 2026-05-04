package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/network"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// VPC is the wrapper for an Aruba Cloud VPC. Construct with aruba.NewVPC()
// and bind it to a project via IntoProject(parent). Pass to VPCsClient.Create
// / .Update or receive from .Get / .List.
type VPC struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	defaultVPC *bool
	preset     *bool
	response   *types.VPCResponse
}

func (v *VPC) IntoProject(p Ref) *VPC        { v.intoProject(p); return v }
func (v *VPC) WithName(n string) *VPC        { v.withName(n); return v }
func (v *VPC) AddTag(t string) *VPC          { v.addTag(t); return v }
func (v *VPC) RemoveTag(t string) *VPC       { v.removeTag(t); return v }
func (v *VPC) ReplaceTags(ts ...string) *VPC { v.replaceTags(ts...); return v }
func (v *VPC) WithLocation(loc string) *VPC  { v.withLocation(loc); return v }
func (v *VPC) InRegion(region string) *VPC   { v.withLocation(region); return v }
func (v *VPC) WithDefault(b bool) *VPC       { v.defaultVPC = &b; return v }
func (v *VPC) WithPreset(b bool) *VPC        { v.preset = &b; return v }

// URI satisfies Ref.
func (v *VPC) URI() string { return v.RespURI() }

// VPCID satisfies withVPCID so children's IntoVPC can extract the parent ID.
func (v *VPC) VPCID() string { return v.ID() }

// Raw shadows the promoted responseMetadataMixin.Raw() returning the full response.
func (v *VPC) Raw() *types.VPCResponse { return v.response }

// RawRequest returns the wire-level request that toRequest() would emit.
func (v *VPC) RawRequest() types.VPCRequest { return v.toRequest() }

// IsDefault returns true if this VPC is the account-region default.
func (v *VPC) IsDefault() bool {
	if v.defaultVPC == nil {
		return false
	}
	return *v.defaultVPC
}

// IsPreset returns true if the VPC was created with a preset subnet/SG.
func (v *VPC) IsPreset() bool {
	if v.preset == nil {
		return false
	}
	return *v.preset
}

func (v *VPC) toRequest() types.VPCRequest {
	var props *types.VPCProperties
	if v.defaultVPC != nil || v.preset != nil {
		props = &types.VPCProperties{Default: v.defaultVPC, Preset: v.preset}
	}
	return types.VPCRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: v.toMetadata(),
			Location:                v.toLocation(),
		},
		Properties: types.VPCPropertiesRequest{Properties: props},
	}
}

func (v *VPC) fromResponse(resp *types.VPCResponse) {
	if resp == nil {
		return
	}
	v.response = resp
	v.setMeta(&resp.Metadata)
	v.withName(vpcDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		v.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		v.withLocation(resp.Metadata.LocationResponse.Value)
	}
	v.setStatus(&resp.Status)
	v.setTerminalStates(vpcTerminalStates)
	v.setLinked(resp.Properties.LinkedResources)
	d := resp.Properties.Default
	v.defaultVPC = &d
	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		v.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
}

func vpcDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

var vpcTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type vpcLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.VPCList], error)
	Get(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Create(ctx context.Context, projectID string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Update(ctx context.Context, projectID, vpcID string, body types.VPCRequest, params *types.RequestParameters) (*types.Response[types.VPCResponse], error)
	Delete(ctx context.Context, projectID, vpcID string, params *types.RequestParameters) (*types.Response[any], error)
}

type vpcsClientAdapter struct{ low vpcLowLevelClient }

func newVPCsClientAdapter(rest *restclient.Client) *vpcsClientAdapter {
	if rest == nil {
		return &vpcsClientAdapter{}
	}
	return &vpcsClientAdapter{low: network.NewVPCsClientImpl(rest)}
}

func (a *vpcsClientAdapter) Create(ctx context.Context, v *VPC, opts ...CallOption) (*VPC, error) {
	if err := v.Err(); err != nil {
		return v, err
	}
	if v.ProjectID() == "" {
		return v, fmt.Errorf("Create: VPC has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, v.ProjectID(), v.toRequest(), rp)
	populateHTTPEnvelope(&v.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		v.fromResponse(resp.Data)
		v.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, v)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				v.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return v, err
	}
	if resp != nil && !resp.IsSuccess() {
		return v, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return v, nil
}

func (a *vpcsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*VPC, error) {
	projectID, vpcID, err := vpcIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, vpcID, rp)
	out := &VPC{}
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
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *vpcsClientAdapter) Update(ctx context.Context, v *VPC, opts ...CallOption) (*VPC, error) {
	if err := v.Err(); err != nil {
		return v, err
	}
	if v.ID() == "" {
		return v, fmt.Errorf("Update: VPC has no ID — call Get first or seed from response metadata")
	}
	if v.ProjectID() == "" {
		return v, fmt.Errorf("Update: VPC has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, v.ProjectID(), v.ID(), v.toRequest(), rp)
	populateHTTPEnvelope(&v.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		v.fromResponse(resp.Data)
		v.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, v)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				v.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return v, err
	}
	if resp != nil && !resp.IsSuccess() {
		return v, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return v, nil
}

func (a *vpcsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, vpcID, err := vpcIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, vpcID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *vpcsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*VPC], error) {
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
	var items []*VPC
	if resp != nil && resp.Data != nil {
		items = make([]*VPC, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			v := &VPC{}
			v.fromResponse(&resp.Data.Values[i])
			v.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, v)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					v.fromResponse(fresh.Raw())
				}
				return nil
			})
			items = append(items, v)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*VPC], error) {
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

// vpcIDsFromRef extracts (projectID, vpcID) from a Ref. Tries typed assertions
// first, then falls back to URI path parsing.
func vpcIDsFromRef(ref Ref) (projectID, vpcID string, err error) {
	vid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid == "" {
		return "", "", fmt.Errorf("cannot determine VPC ID from Ref %q", ref.URI())
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
	return pid, vid, nil
}
