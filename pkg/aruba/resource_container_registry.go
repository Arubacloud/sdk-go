package aruba

import (
	"fmt"
	"strconv"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// ContainerRegistry is the wrapper for an Aruba Cloud Container Registry
// (a direct child of a Project). Construct with aruba.NewContainerRegistry()
// and bind it via IntoProject(project), WithVPC(vpc), etc.
//
// Family A: regional, Metadata/Properties envelope, location-aware.
// Supports full CRUD including Update (same request shape as Create).
//
// Path: /projects/{projectID}/providers/Aruba.Container/registries[/{registryID}]
type ContainerRegistry struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	// Body-refs (single).
	publicIPRef      *string
	vpcRef           *string
	subnetRef        *string
	securityGroupRef *string
	blockStorageRef  *string

	// Registry-specific scalars.
	adminUsername   *string
	concurrentUsers *string // wire "size" — we accept int and itoa; stored as *string
	billingPeriod   *string

	response *types.ContainerRegistryResponse
}

// Standard setters.

func (r *ContainerRegistry) IntoProject(p Ref) *ContainerRegistry  { r.intoProject(p); return r }
func (r *ContainerRegistry) WithName(n string) *ContainerRegistry  { r.withName(n); return r }
func (r *ContainerRegistry) AddTag(t string) *ContainerRegistry    { r.addTag(t); return r }
func (r *ContainerRegistry) RemoveTag(t string) *ContainerRegistry { r.removeTag(t); return r }
func (r *ContainerRegistry) ReplaceTags(ts ...string) *ContainerRegistry {
	r.replaceTags(ts...)
	return r
}
func (r *ContainerRegistry) WithLocation(loc string) *ContainerRegistry {
	r.withLocation(loc)
	return r
}
func (r *ContainerRegistry) InRegion(region string) *ContainerRegistry {
	r.withLocation(region)
	return r
}

// Body-ref setters. Empty URIs are recorded on the error sink and the field
// remains unset.

func (r *ContainerRegistry) WithPublicIP(eip Ref) *ContainerRegistry {
	return r.setSingleRef("WithPublicIP", eip, &r.publicIPRef)
}
func (r *ContainerRegistry) WithVPC(v Ref) *ContainerRegistry {
	return r.setSingleRef("WithVPC", v, &r.vpcRef)
}
func (r *ContainerRegistry) WithSubnet(s Ref) *ContainerRegistry {
	return r.setSingleRef("WithSubnet", s, &r.subnetRef)
}
func (r *ContainerRegistry) WithSecurityGroup(sg Ref) *ContainerRegistry {
	return r.setSingleRef("WithSecurityGroup", sg, &r.securityGroupRef)
}
func (r *ContainerRegistry) WithBlockStorage(vol Ref) *ContainerRegistry {
	return r.setSingleRef("WithBlockStorage", vol, &r.blockStorageRef)
}

func (r *ContainerRegistry) setSingleRef(label string, ref Ref, dst **string) *ContainerRegistry {
	uri := ref.URI()
	if uri == "" {
		r.addErr(fmt.Errorf("%s: empty URI", label))
		return r
	}
	*dst = &uri
	return r
}

// Registry-specific scalar setters.

func (r *ContainerRegistry) WithAdminUsername(u string) *ContainerRegistry {
	r.adminUsername = &u
	return r
}

// WithSize sets the concurrent-users limit. The value is converted to the wire
// string representation ("size" JSON field).
func (r *ContainerRegistry) WithSize(concurrentUsers int) *ContainerRegistry {
	s := strconv.Itoa(concurrentUsers)
	r.concurrentUsers = &s
	return r
}

func (r *ContainerRegistry) WithBillingPeriod(p string) *ContainerRegistry {
	r.billingPeriod = &p
	return r
}

// Ref + ID accessors.

func (r *ContainerRegistry) URI() string                 { return r.RespURI() }
func (r *ContainerRegistry) ContainerRegistryID() string { return r.ID() }

// Raw accessors.

func (r *ContainerRegistry) Raw() *types.ContainerRegistryResponse      { return r.response }
func (r *ContainerRegistry) RawRequest() types.ContainerRegistryRequest { return r.toRequest() }

// Response-preferring accessors (fall back to request-side field when not hydrated).

func (r *ContainerRegistry) PublicIP() string {
	return r.responseURIField(func() string { return r.response.Properties.PublicIp.URI }, r.publicIPRef)
}
func (r *ContainerRegistry) VPC() string {
	return r.responseURIField(func() string { return r.response.Properties.VPC.URI }, r.vpcRef)
}
func (r *ContainerRegistry) Subnet() string {
	return r.responseURIField(func() string { return r.response.Properties.Subnet.URI }, r.subnetRef)
}
func (r *ContainerRegistry) SecurityGroup() string {
	return r.responseURIField(func() string { return r.response.Properties.SecurityGroup.URI }, r.securityGroupRef)
}
func (r *ContainerRegistry) BlockStorage() string {
	return r.responseURIField(func() string { return r.response.Properties.BlockStorage.URI }, r.blockStorageRef)
}
func (r *ContainerRegistry) responseURIField(fromResp func() string, fallback *string) string {
	if r.response != nil {
		if u := fromResp(); u != "" {
			return u
		}
	}
	return containerRegistryDeref(fallback)
}

func (r *ContainerRegistry) AdminUsername() string {
	if r.response != nil && r.response.Properties.AdminUser != nil {
		return r.response.Properties.AdminUser.Username
	}
	return containerRegistryDeref(r.adminUsername)
}

// Size returns the concurrent-users limit. Returns 0 when not set or when the
// wire string cannot be parsed as an integer.
func (r *ContainerRegistry) Size() int {
	if r.response != nil && r.response.Properties.ConcurrentUsers != nil {
		if n, err := strconv.Atoi(*r.response.Properties.ConcurrentUsers); err == nil {
			return n
		}
	}
	if r.concurrentUsers != nil {
		if n, err := strconv.Atoi(*r.concurrentUsers); err == nil {
			return n
		}
	}
	return 0
}

func (r *ContainerRegistry) BillingPeriod() string {
	if r.response != nil && r.response.Properties.BillingPlan != nil {
		return r.response.Properties.BillingPlan.BillingPeriod
	}
	return containerRegistryDeref(r.billingPeriod)
}

func (r *ContainerRegistry) toRequest() types.ContainerRegistryRequest {
	props := types.ContainerRegistryPropertiesRequest{}
	if r.publicIPRef != nil {
		props.PublicIp = types.ReferenceResource{URI: *r.publicIPRef}
	}
	if r.vpcRef != nil {
		props.VPC = types.ReferenceResource{URI: *r.vpcRef}
	}
	if r.subnetRef != nil {
		props.Subnet = types.ReferenceResource{URI: *r.subnetRef}
	}
	if r.securityGroupRef != nil {
		props.SecurityGroup = types.ReferenceResource{URI: *r.securityGroupRef}
	}
	if r.blockStorageRef != nil {
		props.BlockStorage = types.ReferenceResource{URI: *r.blockStorageRef}
	}
	if r.adminUsername != nil {
		props.AdminUser = &types.UserCredential{Username: *r.adminUsername}
	}
	if r.concurrentUsers != nil {
		props.ConcurrentUsers = r.concurrentUsers
	}
	if r.billingPeriod != nil {
		props.BillingPlan = &types.BillingPeriodResource{BillingPeriod: *r.billingPeriod}
	}
	return types.ContainerRegistryRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: r.toMetadata(),
			Location:                r.toLocation(),
		},
		Properties: props,
	}
}

func (r *ContainerRegistry) fromResponse(resp *types.ContainerRegistryResponse) {
	if resp == nil {
		return
	}
	r.response = resp
	r.setMeta(&resp.Metadata)
	r.withName(containerRegistryDeref(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		r.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		r.withLocation(resp.Metadata.LocationResponse.Value)
	}
	r.setStatus(&resp.Status)

	if resp.Properties.PublicIp.URI != "" {
		v := resp.Properties.PublicIp.URI
		r.publicIPRef = &v
	}
	if resp.Properties.VPC.URI != "" {
		v := resp.Properties.VPC.URI
		r.vpcRef = &v
	}
	if resp.Properties.Subnet.URI != "" {
		v := resp.Properties.Subnet.URI
		r.subnetRef = &v
	}
	if resp.Properties.SecurityGroup.URI != "" {
		v := resp.Properties.SecurityGroup.URI
		r.securityGroupRef = &v
	}
	if resp.Properties.BlockStorage.URI != "" {
		v := resp.Properties.BlockStorage.URI
		r.blockStorageRef = &v
	}
	if resp.Properties.AdminUser != nil && resp.Properties.AdminUser.Username != "" {
		v := resp.Properties.AdminUser.Username
		r.adminUsername = &v
	}
	if resp.Properties.ConcurrentUsers != nil && *resp.Properties.ConcurrentUsers != "" {
		v := *resp.Properties.ConcurrentUsers
		r.concurrentUsers = &v
	}
	if resp.Properties.BillingPlan != nil && resp.Properties.BillingPlan.BillingPeriod != "" {
		v := resp.Properties.BillingPlan.BillingPeriod
		r.billingPeriod = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		r.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if r.projectID == "" && r.RespURI() != "" {
		if pid := parseURIIDs(r.RespURI())["projects"]; pid != "" {
			r.projectID = pid
		}
	}
}

func containerRegistryDeref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
