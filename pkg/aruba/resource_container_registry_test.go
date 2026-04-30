package aruba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// --------------------------------------------------------------------------
// Compile-time interface satisfaction
// --------------------------------------------------------------------------

var (
	_ Ref     = (*ContainerRegistry)(nil)
	_ Wrapper = (*ContainerRegistry)(nil)
)

// --------------------------------------------------------------------------
// Fluent setters
// --------------------------------------------------------------------------

func TestContainerRegistry_FluentSetters(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p-1", "my-proj", "/projects/p-1"))

	vpcURI := URI("/projects/p-1/providers/Aruba.Network/vpcs/vpc-1")
	subnetURI := URI("/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1")
	sgURI := URI("/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1")
	eipURI := URI("/projects/p-1/providers/Aruba.Network/elasticips/eip-1")
	bsURI := URI("/projects/p-1/providers/Aruba.Storage/blockstorages/bs-1")

	cr := NewContainerRegistry().
		IntoProject(proj).
		WithName("my-registry").
		AddTag("env:prod").
		AddTag("registry").
		AddTag("env:prod"). // dedupe
		WithLocation("ITBG-1").
		WithVPC(vpcURI).
		WithSubnet(subnetURI).
		WithSecurityGroup(sgURI).
		WithPublicIP(eipURI).
		WithBlockStorage(bsURI).
		WithAdminUsername("admin").
		WithSize(50).
		WithBillingPeriod("Hour")

	if cr.Name() != "my-registry" {
		t.Errorf("Name() = %q", cr.Name())
	}
	if tags := cr.Tags(); len(tags) != 2 || tags[0] != "env:prod" || tags[1] != "registry" {
		t.Errorf("Tags() = %v", tags)
	}
	if cr.Region() != "ITBG-1" {
		t.Errorf("Region() = %q", cr.Region())
	}
	if cr.VPC() != vpcURI.URI() {
		t.Errorf("VPC() = %q", cr.VPC())
	}
	if cr.Subnet() != subnetURI.URI() {
		t.Errorf("Subnet() = %q", cr.Subnet())
	}
	if cr.SecurityGroup() != sgURI.URI() {
		t.Errorf("SecurityGroup() = %q", cr.SecurityGroup())
	}
	if cr.PublicIP() != eipURI.URI() {
		t.Errorf("PublicIP() = %q", cr.PublicIP())
	}
	if cr.BlockStorage() != bsURI.URI() {
		t.Errorf("BlockStorage() = %q", cr.BlockStorage())
	}
	if cr.AdminUsername() != "admin" {
		t.Errorf("AdminUsername() = %q", cr.AdminUsername())
	}
	if cr.Size() != 50 {
		t.Errorf("Size() = %d", cr.Size())
	}
	if cr.BillingPeriod() != "Hour" {
		t.Errorf("BillingPeriod() = %q", cr.BillingPeriod())
	}
	if cr.ProjectID() != "p-1" {
		t.Errorf("ProjectID() = %q", cr.ProjectID())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}

	cr.RemoveTag("env:prod")
	if tags := cr.Tags(); len(tags) != 1 || tags[0] != "registry" {
		t.Errorf("after RemoveTag Tags() = %v", tags)
	}

	cr.ReplaceTags("x", "y")
	if tags := cr.Tags(); len(tags) != 2 || tags[0] != "x" || tags[1] != "y" {
		t.Errorf("after ReplaceTags Tags() = %v", tags)
	}
}

// --------------------------------------------------------------------------
// IntoProject
// --------------------------------------------------------------------------

func TestContainerRegistry_IntoProject_TypedRef(t *testing.T) {
	proj := &Project{}
	proj.fromResponse(projectTestResponse("p-42", "n", "/projects/p-42"))
	cr := NewContainerRegistry().IntoProject(proj)
	if cr.ProjectID() != "p-42" {
		t.Errorf("ProjectID() = %q", cr.ProjectID())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_IntoProject_URIRef(t *testing.T) {
	cr := NewContainerRegistry().IntoProject(URI("/projects/p-uri"))
	if cr.ProjectID() != "p-uri" {
		t.Errorf("ProjectID() = %q", cr.ProjectID())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_IntoProject_BadRef(t *testing.T) {
	cr := NewContainerRegistry().IntoProject(URI("/garbage"))
	if cr.Err() == nil {
		t.Error("expected Err() != nil for unresolvable Ref")
	}
}

// --------------------------------------------------------------------------
// WithPublicIP body-ref setter
// --------------------------------------------------------------------------

func TestContainerRegistry_WithPublicIP_URIRef(t *testing.T) {
	uri := "/projects/p-1/providers/Aruba.Network/elasticips/eip-1"
	cr := NewContainerRegistry().WithPublicIP(URI(uri))
	if cr.PublicIP() != uri {
		t.Errorf("PublicIP() = %q", cr.PublicIP())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_WithPublicIP_TypedRef(t *testing.T) {
	ref := URI("/projects/p-1/providers/Aruba.Network/elasticips/eip-1")
	cr := NewContainerRegistry().WithPublicIP(ref)
	if cr.PublicIP() != ref.URI() {
		t.Errorf("PublicIP() = %q, want %q", cr.PublicIP(), ref.URI())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_WithPublicIP_EmptyURI(t *testing.T) {
	cr := NewContainerRegistry().WithPublicIP(URI(""))
	if cr.Err() == nil {
		t.Error("expected Err() != nil for empty PublicIP URI")
	}
	if cr.PublicIP() != "" {
		t.Errorf("PublicIP() should remain empty, got %q", cr.PublicIP())
	}
}

// --------------------------------------------------------------------------
// WithVPC body-ref setter
// --------------------------------------------------------------------------

func TestContainerRegistry_WithVPC_URIRef(t *testing.T) {
	uri := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1"
	cr := NewContainerRegistry().WithVPC(URI(uri))
	if cr.VPC() != uri {
		t.Errorf("VPC() = %q", cr.VPC())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_WithVPC_TypedRef(t *testing.T) {
	ref := URI("/projects/p-1/providers/Aruba.Network/vpcs/vpc-1")
	cr := NewContainerRegistry().WithVPC(ref)
	if cr.VPC() != ref.URI() {
		t.Errorf("VPC() = %q, want %q", cr.VPC(), ref.URI())
	}
}

func TestContainerRegistry_WithVPC_EmptyURI(t *testing.T) {
	cr := NewContainerRegistry().WithVPC(URI(""))
	if cr.Err() == nil {
		t.Error("expected Err() != nil for empty VPC URI")
	}
	if cr.VPC() != "" {
		t.Errorf("VPC() should remain empty, got %q", cr.VPC())
	}
}

// --------------------------------------------------------------------------
// WithSubnet body-ref setter
// --------------------------------------------------------------------------

func TestContainerRegistry_WithSubnet_URIRef(t *testing.T) {
	uri := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1"
	cr := NewContainerRegistry().WithSubnet(URI(uri))
	if cr.Subnet() != uri {
		t.Errorf("Subnet() = %q", cr.Subnet())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_WithSubnet_TypedRef(t *testing.T) {
	ref := URI("/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1")
	cr := NewContainerRegistry().WithSubnet(ref)
	if cr.Subnet() != ref.URI() {
		t.Errorf("Subnet() = %q, want %q", cr.Subnet(), ref.URI())
	}
}

func TestContainerRegistry_WithSubnet_EmptyURI(t *testing.T) {
	cr := NewContainerRegistry().WithSubnet(URI(""))
	if cr.Err() == nil {
		t.Error("expected Err() != nil for empty Subnet URI")
	}
	if cr.Subnet() != "" {
		t.Errorf("Subnet() should remain empty, got %q", cr.Subnet())
	}
}

// --------------------------------------------------------------------------
// WithSecurityGroup body-ref setter
// --------------------------------------------------------------------------

func TestContainerRegistry_WithSecurityGroup_URIRef(t *testing.T) {
	uri := "/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1"
	cr := NewContainerRegistry().WithSecurityGroup(URI(uri))
	if cr.SecurityGroup() != uri {
		t.Errorf("SecurityGroup() = %q", cr.SecurityGroup())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_WithSecurityGroup_TypedRef(t *testing.T) {
	ref := URI("/projects/p-1/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1")
	cr := NewContainerRegistry().WithSecurityGroup(ref)
	if cr.SecurityGroup() != ref.URI() {
		t.Errorf("SecurityGroup() = %q, want %q", cr.SecurityGroup(), ref.URI())
	}
}

func TestContainerRegistry_WithSecurityGroup_EmptyURI(t *testing.T) {
	cr := NewContainerRegistry().WithSecurityGroup(URI(""))
	if cr.Err() == nil {
		t.Error("expected Err() != nil for empty SecurityGroup URI")
	}
	if cr.SecurityGroup() != "" {
		t.Errorf("SecurityGroup() should remain empty, got %q", cr.SecurityGroup())
	}
}

// --------------------------------------------------------------------------
// WithBlockStorage body-ref setter
// --------------------------------------------------------------------------

func TestContainerRegistry_WithBlockStorage_URIRef(t *testing.T) {
	uri := "/projects/p-1/providers/Aruba.Storage/blockstorages/bs-1"
	cr := NewContainerRegistry().WithBlockStorage(URI(uri))
	if cr.BlockStorage() != uri {
		t.Errorf("BlockStorage() = %q", cr.BlockStorage())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_WithBlockStorage_TypedRef(t *testing.T) {
	ref := URI("/projects/p-1/providers/Aruba.Storage/blockstorages/bs-1")
	cr := NewContainerRegistry().WithBlockStorage(ref)
	if cr.BlockStorage() != ref.URI() {
		t.Errorf("BlockStorage() = %q, want %q", cr.BlockStorage(), ref.URI())
	}
}

func TestContainerRegistry_WithBlockStorage_EmptyURI(t *testing.T) {
	cr := NewContainerRegistry().WithBlockStorage(URI(""))
	if cr.Err() == nil {
		t.Error("expected Err() != nil for empty BlockStorage URI")
	}
	if cr.BlockStorage() != "" {
		t.Errorf("BlockStorage() should remain empty, got %q", cr.BlockStorage())
	}
}

// --------------------------------------------------------------------------
// Registry scalars
// --------------------------------------------------------------------------

func TestContainerRegistry_WithAdminUsername(t *testing.T) {
	cr := NewContainerRegistry().WithAdminUsername("myuser")
	if cr.AdminUsername() != "myuser" {
		t.Errorf("AdminUsername() = %q", cr.AdminUsername())
	}
	if cr.Err() != nil {
		t.Errorf("Err() = %v", cr.Err())
	}
}

func TestContainerRegistry_WithSize(t *testing.T) {
	cr := NewContainerRegistry().WithSize(42)
	if cr.Size() != 42 {
		t.Errorf("Size() = %d", cr.Size())
	}
	// Verify the wire representation is the string "42"
	req := cr.RawRequest()
	if req.Properties.ConcurrentUsers == nil || *req.Properties.ConcurrentUsers != "42" {
		t.Errorf("wire ConcurrentUsers = %v", req.Properties.ConcurrentUsers)
	}
}

func TestContainerRegistry_WithBillingPeriod(t *testing.T) {
	cr := NewContainerRegistry().WithBillingPeriod("Monthly")
	if cr.BillingPeriod() != "Monthly" {
		t.Errorf("BillingPeriod() = %q", cr.BillingPeriod())
	}
}

// --------------------------------------------------------------------------
// toRequest round-trip
// --------------------------------------------------------------------------

func TestContainerRegistry_ToRequest(t *testing.T) {
	vpcURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1"
	subnetURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1"
	sgURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1"
	eipURI := "/projects/p/providers/Aruba.Network/elasticips/eip-1"
	bsURI := "/projects/p/providers/Aruba.Storage/blockstorages/bs-1"

	cr := NewContainerRegistry().
		WithName("reg-rt").
		AddTag("t1").AddTag("t2").
		WithLocation("ITBG-1").
		WithVPC(URI(vpcURI)).
		WithSubnet(URI(subnetURI)).
		WithSecurityGroup(URI(sgURI)).
		WithPublicIP(URI(eipURI)).
		WithBlockStorage(URI(bsURI)).
		WithAdminUsername("admin").
		WithSize(100).
		WithBillingPeriod("Hour")

	req := cr.RawRequest()

	if req.Metadata.Name != "reg-rt" {
		t.Errorf("Metadata.Name = %q", req.Metadata.Name)
	}
	if len(req.Metadata.Tags) != 2 {
		t.Errorf("Metadata.Tags = %v", req.Metadata.Tags)
	}
	if req.Metadata.Location.Value != "ITBG-1" {
		t.Errorf("Location.Value = %q", req.Metadata.Location.Value)
	}
	if req.Properties.VPC.URI != vpcURI {
		t.Errorf("Properties.VPC.URI = %q", req.Properties.VPC.URI)
	}
	if req.Properties.Subnet.URI != subnetURI {
		t.Errorf("Properties.Subnet.URI = %q", req.Properties.Subnet.URI)
	}
	if req.Properties.SecurityGroup.URI != sgURI {
		t.Errorf("Properties.SecurityGroup.URI = %q", req.Properties.SecurityGroup.URI)
	}
	if req.Properties.PublicIp.URI != eipURI {
		t.Errorf("Properties.PublicIp.URI = %q", req.Properties.PublicIp.URI)
	}
	if req.Properties.BlockStorage.URI != bsURI {
		t.Errorf("Properties.BlockStorage.URI = %q", req.Properties.BlockStorage.URI)
	}
	if req.Properties.AdminUser == nil || req.Properties.AdminUser.Username != "admin" {
		t.Errorf("Properties.AdminUser = %v", req.Properties.AdminUser)
	}
	if req.Properties.ConcurrentUsers == nil || *req.Properties.ConcurrentUsers != "100" {
		t.Errorf("Properties.ConcurrentUsers = %v", req.Properties.ConcurrentUsers)
	}
	if req.Properties.BillingPlan == nil || req.Properties.BillingPlan.BillingPeriod != "Hour" {
		t.Errorf("Properties.BillingPlan = %v", req.Properties.BillingPlan)
	}
}

func TestContainerRegistry_ToRequest_Empty(t *testing.T) {
	cr := NewContainerRegistry()
	req := cr.RawRequest() // must not panic

	// All optional pointer fields should be nil when not set.
	if req.Properties.AdminUser != nil {
		t.Errorf("AdminUser should be nil, got %v", req.Properties.AdminUser)
	}
	if req.Properties.ConcurrentUsers != nil {
		t.Errorf("ConcurrentUsers should be nil, got %v", req.Properties.ConcurrentUsers)
	}
	if req.Properties.BillingPlan != nil {
		t.Errorf("BillingPlan should be nil, got %v", req.Properties.BillingPlan)
	}
	// Body-ref ReferenceResource fields should have empty URIs.
	if req.Properties.VPC.URI != "" {
		t.Errorf("VPC.URI should be empty, got %q", req.Properties.VPC.URI)
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration helpers
// --------------------------------------------------------------------------

func containerRegistryTestResponse(name string) *types.ContainerRegistryResponse {
	id := "cr-1"
	uri := "/projects/p/providers/Aruba.Container/registries/cr-1"
	state := "Active"
	size := "50"
	vpcURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1"
	subnetURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1"
	sgURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1"
	eipURI := "/projects/p/providers/Aruba.Network/elasticips/eip-1"
	bsURI := "/projects/p/providers/Aruba.Storage/blockstorages/bs-1"
	return &types.ContainerRegistryResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:               &id,
			URI:              &uri,
			Name:             func() *string { s := name; return &s }(),
			Tags:             []string{"tag1"},
			LocationResponse: &types.LocationResponse{Value: "ITBG-1"},
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: "p",
			},
		},
		Properties: types.ContainerRegistryPropertiesResult{
			VPC:             types.ReferenceResource{URI: vpcURI},
			Subnet:          types.ReferenceResource{URI: subnetURI},
			SecurityGroup:   types.ReferenceResource{URI: sgURI},
			PublicIp:        types.ReferenceResource{URI: eipURI},
			BlockStorage:    types.ReferenceResource{URI: bsURI},
			AdminUser:       &types.UserCredential{Username: "admin"},
			ConcurrentUsers: &size,
			BillingPlan:     &types.BillingPeriodResource{BillingPeriod: "Hour"},
		},
		Status: types.ResourceStatus{
			State: &state,
		},
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration tests
// --------------------------------------------------------------------------

func TestContainerRegistry_FromResponseHydration(t *testing.T) {
	cr := &ContainerRegistry{}
	resp := containerRegistryTestResponse("my-registry")
	cr.fromResponse(resp)

	if cr.ID() != "cr-1" {
		t.Errorf("ID() = %q", cr.ID())
	}
	if cr.ContainerRegistryID() != "cr-1" {
		t.Errorf("ContainerRegistryID() = %q", cr.ContainerRegistryID())
	}
	if cr.URI() != "/projects/p/providers/Aruba.Container/registries/cr-1" {
		t.Errorf("URI() = %q", cr.URI())
	}
	if cr.Name() != "my-registry" {
		t.Errorf("Name() = %q", cr.Name())
	}
	if tags := cr.Tags(); len(tags) != 1 || tags[0] != "tag1" {
		t.Errorf("Tags() = %v", tags)
	}
	if cr.Region() != "ITBG-1" {
		t.Errorf("Region() = %q", cr.Region())
	}
	if cr.State() != "Active" {
		t.Errorf("State() = %q", cr.State())
	}
	if cr.VPC() != "/projects/p/providers/Aruba.Network/vpcs/vpc-1" {
		t.Errorf("VPC() = %q", cr.VPC())
	}
	if cr.Subnet() != "/projects/p/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1" {
		t.Errorf("Subnet() = %q", cr.Subnet())
	}
	if cr.SecurityGroup() != "/projects/p/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1" {
		t.Errorf("SecurityGroup() = %q", cr.SecurityGroup())
	}
	if cr.PublicIP() != "/projects/p/providers/Aruba.Network/elasticips/eip-1" {
		t.Errorf("PublicIP() = %q", cr.PublicIP())
	}
	if cr.BlockStorage() != "/projects/p/providers/Aruba.Storage/blockstorages/bs-1" {
		t.Errorf("BlockStorage() = %q", cr.BlockStorage())
	}
	if cr.AdminUsername() != "admin" {
		t.Errorf("AdminUsername() = %q", cr.AdminUsername())
	}
	if cr.Size() != 50 {
		t.Errorf("Size() = %d", cr.Size())
	}
	if cr.BillingPeriod() != "Hour" {
		t.Errorf("BillingPeriod() = %q", cr.BillingPeriod())
	}
	if cr.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", cr.ProjectID())
	}
	if cr.Raw() != resp {
		t.Error("Raw() should return the hydrated response pointer")
	}
}

func TestContainerRegistry_FromResponse_NilSafe(t *testing.T) {
	cr := &ContainerRegistry{}
	cr.fromResponse(nil) // must not panic
	if cr.ID() != "" || cr.URI() != "" || cr.Name() != "" {
		t.Error("fromResponse(nil) should be a no-op")
	}
}

func TestContainerRegistry_FromResponse_BackfillsProjectID_FromMetadata(t *testing.T) {
	resp := containerRegistryTestResponse("n")
	cr := &ContainerRegistry{}
	cr.fromResponse(resp)
	if cr.ProjectID() != "p" {
		t.Errorf("ProjectID() from metadata = %q", cr.ProjectID())
	}
}

func TestContainerRegistry_FromResponse_BackfillsProjectID_FromURI(t *testing.T) {
	id := "cr-99"
	uri := "/projects/p-uri/providers/Aruba.Container/registries/cr-99"
	resp := &types.ContainerRegistryResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:  &id,
			URI: &uri,
			// No ProjectResponseMetadata — should backfill from URI.
		},
	}
	cr := &ContainerRegistry{}
	cr.fromResponse(resp)
	if cr.ProjectID() != "p-uri" {
		t.Errorf("ProjectID() via URI backfill = %q", cr.ProjectID())
	}
}

// --------------------------------------------------------------------------
// containerRegistryIDsFromRef helper
// --------------------------------------------------------------------------

func TestContainerRegistryIDsFromRef_URIRef(t *testing.T) {
	ref := URI("/projects/p/providers/Aruba.Container/registries/cr-1")
	pid, rid, err := containerRegistryIDsFromRef(ref)
	if err != nil || pid != "p" || rid != "cr-1" {
		t.Errorf("containerRegistryIDsFromRef = (%q, %q, %v)", pid, rid, err)
	}
}

func TestContainerRegistryIDsFromRef_BadURI_NoRegistries(t *testing.T) {
	_, _, err := containerRegistryIDsFromRef(URI("/projects/p/providers/Aruba.Container/something/else"))
	if err == nil {
		t.Error("expected error for URI without /registries/<id>")
	}
}

func TestContainerRegistryIDsFromRef_BadURI_NoProject(t *testing.T) {
	_, _, err := containerRegistryIDsFromRef(URI("/providers/Aruba.Container/registries/cr-1"))
	if err == nil {
		t.Error("expected error for URI without /projects/<id>")
	}
}

// --------------------------------------------------------------------------
// fakeContainerRegistryLowLevel — body-capture tests
// --------------------------------------------------------------------------

type fakeContainerRegistryLowLevel struct {
	createFunc func(ctx context.Context, projectID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	updateFunc func(ctx context.Context, projectID, registryID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	getFunc    func(ctx context.Context, projectID, registryID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error)
	deleteFunc func(ctx context.Context, projectID, registryID string, params *types.RequestParameters) (*types.Response[any], error)
	listFunc   func(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryList], error)
}

func (f *fakeContainerRegistryLowLevel) Create(ctx context.Context, projectID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error) {
	return f.createFunc(ctx, projectID, body, params)
}
func (f *fakeContainerRegistryLowLevel) Update(ctx context.Context, projectID, registryID string, body types.ContainerRegistryRequest, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error) {
	return f.updateFunc(ctx, projectID, registryID, body, params)
}
func (f *fakeContainerRegistryLowLevel) Get(ctx context.Context, projectID, registryID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error) {
	return f.getFunc(ctx, projectID, registryID, params)
}
func (f *fakeContainerRegistryLowLevel) Delete(ctx context.Context, projectID, registryID string, params *types.RequestParameters) (*types.Response[any], error) {
	return f.deleteFunc(ctx, projectID, registryID, params)
}
func (f *fakeContainerRegistryLowLevel) List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.ContainerRegistryList], error) {
	return f.listFunc(ctx, projectID, params)
}

// --------------------------------------------------------------------------
// HTTP-mock adapter helper
// --------------------------------------------------------------------------

func buildContainerRegistryTestAdapter(t *testing.T, handler http.HandlerFunc) *containerRegistriesClientAdapter {
	t.Helper()
	server := testutil.NewMockServer(t, handler)
	return newContainerRegistriesClientAdapter(testutil.NewClient(t, server.URL))
}

const containerRegistrySuccessBody = `{` +
	`"metadata":{"id":"cr-1","name":"my-registry","uri":"/projects/p/providers/Aruba.Container/registries/cr-1","project":{"id":"p"}},` +
	`"properties":{` +
	`"vpc":{"uri":"/projects/p/providers/Aruba.Network/vpcs/vpc-1"},` +
	`"subnet":{"uri":"/projects/p/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1"},` +
	`"securityGroup":{"uri":"/projects/p/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1"},` +
	`"publicIp":{"uri":"/projects/p/providers/Aruba.Network/elasticips/eip-1"},` +
	`"blockStorage":{"uri":"/projects/p/providers/Aruba.Storage/blockstorages/bs-1"},` +
	`"adminUser":{"username":"admin"},"size":"50","billingPlan":{"billingPeriod":"Hour"}` +
	`},` +
	`"status":{"state":"Active"}}`

// --------------------------------------------------------------------------
// Create adapter tests
// --------------------------------------------------------------------------

func TestContainerRegistriesClientAdapter_Create_Success(t *testing.T) {
	var gotBody types.ContainerRegistryRequest
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		if !containsSubstring(r.URL.Path, "registries") {
			t.Errorf("path %q should contain 'registries'", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, containerRegistrySuccessBody)
	})

	cr := NewContainerRegistry().
		IntoProject(URI("/projects/p")).
		WithName("my-registry").
		WithLocation("ITBG-1").
		WithVPC(URI("/projects/p/providers/Aruba.Network/vpcs/vpc-1")).
		WithSubnet(URI("/projects/p/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1")).
		WithSecurityGroup(URI("/projects/p/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1")).
		WithPublicIP(URI("/projects/p/providers/Aruba.Network/elasticips/eip-1")).
		WithBlockStorage(URI("/projects/p/providers/Aruba.Storage/blockstorages/bs-1")).
		WithAdminUsername("admin").
		WithSize(50).
		WithBillingPeriod("Hour")

	result, err := adapter.Create(context.Background(), cr)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.ID() != "cr-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.Name() != "my-registry" {
		t.Errorf("Name() = %q", result.Name())
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	// Wire body assertions
	if gotBody.Metadata.Name != "my-registry" {
		t.Errorf("request Metadata.Name = %q", gotBody.Metadata.Name)
	}
	if gotBody.Metadata.Location.Value != "ITBG-1" {
		t.Errorf("request Metadata.Location.Value = %q", gotBody.Metadata.Location.Value)
	}
	if gotBody.Properties.VPC.URI != "/projects/p/providers/Aruba.Network/vpcs/vpc-1" {
		t.Errorf("request Properties.VPC.URI = %q", gotBody.Properties.VPC.URI)
	}
}

func TestContainerRegistriesClientAdapter_Create_NoProject(t *testing.T) {
	callCount := 0
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})

	_, err := adapter.Create(context.Background(), NewContainerRegistry().WithName("x"))
	if err == nil {
		t.Fatal("expected error when ContainerRegistry has no project")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without project")
	}
}

func TestContainerRegistriesClientAdapter_Create_MetadataValidationError(t *testing.T) {
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Missing "id" field — triggers MetadataValidationError from low-level Validate()
		fmt.Fprint(w, `{"metadata":{"name":"reg","uri":"/projects/p/providers/Aruba.Container/registries/x"},"properties":{},"status":{}}`)
	})

	cr := NewContainerRegistry().IntoProject(URI("/projects/p")).WithName("reg")
	result, err := adapter.Create(context.Background(), cr)
	if err == nil {
		t.Fatal("expected MetadataValidationError, got nil")
	}
	var mvErr *types.MetadataValidationError
	if !errors.As(err, &mvErr) {
		t.Fatalf("expected *types.MetadataValidationError, got %T: %v", err, err)
	}
	if result == nil {
		t.Fatal("result must be non-nil alongside MetadataValidationError")
	}
}

func TestContainerRegistriesClientAdapter_Create_NonTwoXX(t *testing.T) {
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Validation Failed", "name is required", 422))
	})

	cr := NewContainerRegistry().IntoProject(URI("/projects/p"))
	result, err := adapter.Create(context.Background(), cr)
	if err == nil {
		t.Fatal("expected error on 422")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("HTTPError.StatusCode = %d", httpErr.StatusCode)
	}
	if result == nil {
		t.Fatal("result must be non-nil on non-2xx")
	}
}

func TestContainerRegistriesClientAdapter_Create_WithBodyRefs_ViaFake(t *testing.T) {
	vpcURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1"
	subnetURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1/subnets/sn-1"
	sgURI := "/projects/p/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-1"
	eipURI := "/projects/p/providers/Aruba.Network/elasticips/eip-1"
	bsURI := "/projects/p/providers/Aruba.Storage/blockstorages/bs-1"

	var captured types.ContainerRegistryRequest
	resp := &types.Response[types.ContainerRegistryResponse]{
		StatusCode: http.StatusCreated,
		Data:       containerRegistryTestResponse("reg"),
	}
	fake := &fakeContainerRegistryLowLevel{
		createFunc: func(_ context.Context, _ string, body types.ContainerRegistryRequest, _ *types.RequestParameters) (*types.Response[types.ContainerRegistryResponse], error) {
			captured = body
			return resp, nil
		},
	}
	adapter := &containerRegistriesClientAdapter{low: fake}

	cr := NewContainerRegistry().
		IntoProject(URI("/projects/p")).
		WithLocation("ITBG-1").
		WithVPC(URI(vpcURI)).
		WithSubnet(URI(subnetURI)).
		WithSecurityGroup(URI(sgURI)).
		WithPublicIP(URI(eipURI)).
		WithBlockStorage(URI(bsURI)).
		WithAdminUsername("admin").
		WithSize(100)

	_, err := adapter.Create(context.Background(), cr)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if captured.Properties.VPC.URI != vpcURI {
		t.Errorf("captured VPC.URI = %q", captured.Properties.VPC.URI)
	}
	if captured.Properties.Subnet.URI != subnetURI {
		t.Errorf("captured Subnet.URI = %q", captured.Properties.Subnet.URI)
	}
	if captured.Properties.SecurityGroup.URI != sgURI {
		t.Errorf("captured SecurityGroup.URI = %q", captured.Properties.SecurityGroup.URI)
	}
	if captured.Properties.PublicIp.URI != eipURI {
		t.Errorf("captured PublicIp.URI = %q", captured.Properties.PublicIp.URI)
	}
	if captured.Properties.BlockStorage.URI != bsURI {
		t.Errorf("captured BlockStorage.URI = %q", captured.Properties.BlockStorage.URI)
	}
	if captured.Properties.AdminUser == nil || captured.Properties.AdminUser.Username != "admin" {
		t.Errorf("captured AdminUser = %v", captured.Properties.AdminUser)
	}
	if captured.Properties.ConcurrentUsers == nil || *captured.Properties.ConcurrentUsers != "100" {
		t.Errorf("captured ConcurrentUsers = %v", captured.Properties.ConcurrentUsers)
	}
}

// --------------------------------------------------------------------------
// Update adapter tests
// --------------------------------------------------------------------------

func TestContainerRegistriesClientAdapter_Update_Success(t *testing.T) {
	var capturedPath string
	var gotBody types.ContainerRegistryRequest
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, containerRegistrySuccessBody)
	})

	cr := NewContainerRegistry().
		IntoProject(URI("/projects/p")).
		WithName("my-registry").
		WithLocation("ITBG-1").
		WithVPC(URI("/projects/p/providers/Aruba.Network/vpcs/vpc-1"))

	// Hydrate to get an ID before Update.
	cr.fromResponse(containerRegistryTestResponse("my-registry"))

	result, err := adapter.Update(context.Background(), cr)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if result.ID() != "cr-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if !containsSubstring(capturedPath, "cr-1") {
		t.Errorf("path %q should contain registry ID 'cr-1'", capturedPath)
	}
}

func TestContainerRegistriesClientAdapter_Update_NoID(t *testing.T) {
	callCount := 0
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	cr := NewContainerRegistry().IntoProject(URI("/projects/p")).WithName("x")
	_, err := adapter.Update(context.Background(), cr)
	if err == nil {
		t.Fatal("expected error when ContainerRegistry has no ID")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without ID")
	}
}

func TestContainerRegistriesClientAdapter_Update_NoProject(t *testing.T) {
	callCount := 0
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	// Hydrate to get an ID but without a project.
	cr := &ContainerRegistry{}
	cr.fromResponse(containerRegistryTestResponse("n"))
	cr.projectID = "" // strip project

	_, err := adapter.Update(context.Background(), cr)
	if err == nil {
		t.Fatal("expected error when ContainerRegistry has no project")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without project")
	}
}

func TestContainerRegistriesClientAdapter_Update_NonTwoXX(t *testing.T) {
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Conflict", "registry already exists", 409))
	})

	cr := &ContainerRegistry{}
	cr.fromResponse(containerRegistryTestResponse("n"))

	result, err := adapter.Update(context.Background(), cr)
	if err == nil {
		t.Fatal("expected error on 409")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusConflict {
		t.Errorf("HTTPError.StatusCode = %d", httpErr.StatusCode)
	}
	if result == nil {
		t.Fatal("result must be non-nil on non-2xx")
	}
}

// --------------------------------------------------------------------------
// Get adapter tests
// --------------------------------------------------------------------------

func TestContainerRegistriesClientAdapter_Get_URIRef(t *testing.T) {
	var capturedPath string
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, containerRegistrySuccessBody)
	})

	ref := URI("/projects/p/providers/Aruba.Container/registries/cr-1")
	result, err := adapter.Get(context.Background(), ref)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "cr-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", result.ProjectID())
	}
	if !containsSubstring(capturedPath, "registries") {
		t.Errorf("path %q should contain 'registries'", capturedPath)
	}
}

func TestContainerRegistriesClientAdapter_Get_TypedRef(t *testing.T) {
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, containerRegistrySuccessBody)
	})

	existing := &ContainerRegistry{}
	existing.fromResponse(containerRegistryTestResponse("my-registry"))

	result, err := adapter.Get(context.Background(), existing)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "cr-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

// --------------------------------------------------------------------------
// Delete adapter tests
// --------------------------------------------------------------------------

func TestContainerRegistriesClientAdapter_Delete_Success(t *testing.T) {
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/providers/Aruba.Container/registries/cr-1"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestContainerRegistriesClientAdapter_Delete_NonTwoXX(t *testing.T) {
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "registry not found", 404))
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/providers/Aruba.Container/registries/missing"))
	if err == nil {
		t.Fatal("expected error on 404")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T", err)
	}
	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d", httpErr.StatusCode)
	}
}

// --------------------------------------------------------------------------
// List adapter tests
// --------------------------------------------------------------------------

func TestContainerRegistriesClientAdapter_List_TwoItems(t *testing.T) {
	adapter := buildContainerRegistryTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total":2,"self":"","prev":"","next":"","first":"","last":"","values":[`+
			`{"metadata":{"id":"cr-1","name":"n1","uri":"/projects/p/providers/Aruba.Container/registries/cr-1","project":{"id":"p"}},"properties":{"vpc":{"uri":"/projects/p/providers/Aruba.Network/vpcs/vpc-1"},"size":"10","billingPlan":{"billingPeriod":"Hour"}},"status":{}},`+
			`{"metadata":{"id":"cr-2","name":"n2","uri":"/projects/p/providers/Aruba.Container/registries/cr-2","project":{"id":"p"}},"properties":{"vpc":{"uri":"/projects/p/providers/Aruba.Network/vpcs/vpc-1"},"size":"20","billingPlan":{"billingPeriod":"Monthly"}},"status":{}}`+
			`]}`)
	})

	list, err := adapter.List(context.Background(), URI("/projects/p"))
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if list.Total() != 2 {
		t.Errorf("Total() = %d", list.Total())
	}
	items := list.Items()
	if len(items) != 2 {
		t.Fatalf("Items() len = %d", len(items))
	}
	if items[0].ID() != "cr-1" || items[0].Name() != "n1" {
		t.Errorf("items[0] = {%q, %q}", items[0].ID(), items[0].Name())
	}
	if items[0].Size() != 10 {
		t.Errorf("items[0].Size() = %d", items[0].Size())
	}
	if items[1].ID() != "cr-2" || items[1].BillingPeriod() != "Monthly" {
		t.Errorf("items[1] ID=%q BillingPeriod=%q", items[1].ID(), items[1].BillingPeriod())
	}
	if items[0].ProjectID() != "p" {
		t.Errorf("items[0].ProjectID() = %q", items[0].ProjectID())
	}
}

// --------------------------------------------------------------------------
// Reflective check: ContainerRegistryClient has Update method
// --------------------------------------------------------------------------

func TestContainerRegistryClient_HasUpdateMethod(t *testing.T) {
	iface := reflect.TypeOf((*ContainerRegistryClient)(nil)).Elem()
	for i := range iface.NumMethod() {
		if iface.Method(i).Name == "Update" {
			return // found — test passes
		}
	}
	t.Fatal("ContainerRegistryClient must have an Update method")
}
