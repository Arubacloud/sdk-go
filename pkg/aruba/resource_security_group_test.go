package aruba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// --------------------------------------------------------------------------
// Compile-time Ref satisfaction
// --------------------------------------------------------------------------

var _ Ref = (*SecurityGroup)(nil)

// --------------------------------------------------------------------------
// Fluent setters
// --------------------------------------------------------------------------

func TestSecurityGroup_FluentSetters(t *testing.T) {
	parent := &VPC{}
	parent.fromResponse(vpcTestResponse("v1", "my-vpc", "/projects/p1/providers/Aruba.Network/vpcs/v1", "p1"))

	sg := NewSecurityGroup().
		IntoVPC(parent).
		WithName("my-sg").
		AddTag("security").
		AddTag("network").
		AddTag("security"). // dedupe
		WithDefault(true)

	if sg.Name() != "my-sg" {
		t.Errorf("Name() = %q", sg.Name())
	}
	if tags := sg.Tags(); len(tags) != 2 || tags[0] != "security" || tags[1] != "network" {
		t.Errorf("Tags() = %v", tags)
	}
	if sg.Default() != true {
		t.Errorf("Default() = %v", sg.Default())
	}
	if sg.VPCID() != "v1" {
		t.Errorf("VPCID() = %q", sg.VPCID())
	}
	if sg.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", sg.ProjectID())
	}
	if sg.Err() != nil {
		t.Errorf("Err() = %v", sg.Err())
	}

	sg.RemoveTag("security")
	if tags := sg.Tags(); len(tags) != 1 || tags[0] != "network" {
		t.Errorf("after RemoveTag Tags() = %v", tags)
	}

	sg.ReplaceTags("x", "y")
	if tags := sg.Tags(); len(tags) != 2 || tags[0] != "x" || tags[1] != "y" {
		t.Errorf("after ReplaceTags Tags() = %v", tags)
	}
}

// --------------------------------------------------------------------------
// IntoVPC with bad Ref
// --------------------------------------------------------------------------

func TestSecurityGroup_IntoVPC_BadRef(t *testing.T) {
	sg := NewSecurityGroup().IntoVPC(URI("/garbage"))
	if sg.Err() == nil {
		t.Error("expected Err() != nil for unresolvable Ref, got nil")
	}
}

// --------------------------------------------------------------------------
// toRequest round-trip
// --------------------------------------------------------------------------

func TestSecurityGroup_ToRequestRoundTrip(t *testing.T) {
	sg := NewSecurityGroup().
		WithName("sg-1").
		AddTag("t1").
		AddTag("t2").
		WithDefault(false)

	req := sg.RawRequest()

	if req.Metadata.Name != "sg-1" {
		t.Errorf("Metadata.Name = %q", req.Metadata.Name)
	}
	if len(req.Metadata.Tags) != 2 {
		t.Errorf("Metadata.Tags = %v", req.Metadata.Tags)
	}
	if req.Properties.Default == nil || *req.Properties.Default != false {
		t.Errorf("Properties.Default = %v", req.Properties.Default)
	}

	// When WithDefault is not called, Default pointer must be nil (omitempty).
	sg2 := NewSecurityGroup().WithName("bare")
	req2 := sg2.RawRequest()
	if req2.Properties.Default != nil {
		t.Errorf("Default should be nil when not set, got %v", *req2.Properties.Default)
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration
// --------------------------------------------------------------------------

func securityGroupTestResponse(id, name, uri, projectID string) *types.SecurityGroupResponse {
	state := "Active"
	return &types.SecurityGroupResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:   &id,
			URI:  &uri,
			Name: &name,
			Tags: []string{"sg-tag"},
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: projectID,
			},
		},
		Properties: types.SecurityGroupPropertiesResponse{
			Default: true,
			LinkedResources: []types.LinkedResource{
				{URI: "/projects/p/providers/Aruba.Compute/cloudservers/cs1", StrictCorrelation: true},
			},
		},
		Status: types.ResourceStatus{
			State: &state,
			DisableStatusInfo: &types.DisableStatusInfo{
				IsDisabled: false,
			},
		},
	}
}

func TestSecurityGroup_FromResponseHydration(t *testing.T) {
	sg := &SecurityGroup{}
	resp := securityGroupTestResponse("sg-1", "my-sg", "/projects/p1/network/vpcs/v1/security-groups/sg-1", "p1")
	sg.fromResponse(resp)

	if sg.ID() != "sg-1" {
		t.Errorf("ID() = %q", sg.ID())
	}
	if sg.URI() != "/projects/p1/network/vpcs/v1/security-groups/sg-1" {
		t.Errorf("URI() = %q", sg.URI())
	}
	if sg.SecurityGroupID() != "sg-1" {
		t.Errorf("SecurityGroupID() = %q", sg.SecurityGroupID())
	}
	if sg.Name() != "my-sg" {
		t.Errorf("Name() = %q", sg.Name())
	}
	if tags := sg.Tags(); len(tags) != 1 || tags[0] != "sg-tag" {
		t.Errorf("Tags() = %v", tags)
	}
	if sg.State() != "Active" {
		t.Errorf("State() = %q", sg.State())
	}
	if sg.IsDisabled() {
		t.Error("IsDisabled() should be false")
	}
	if linked := sg.LinkedResources(); len(linked) != 1 {
		t.Errorf("LinkedResources() len = %d", len(linked))
	}
	if sg.Default() != true {
		t.Errorf("Default() = %v", sg.Default())
	}
	if sg.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", sg.ProjectID())
	}
	if sg.VPCID() != "v1" {
		t.Errorf("VPCID() = %q", sg.VPCID())
	}
	if sg.Raw() != resp {
		t.Error("Raw() should return the hydrated response pointer")
	}
}

func TestSecurityGroup_FromResponsePartial(t *testing.T) {
	// nil response is a no-op
	sg := &SecurityGroup{}
	sg.fromResponse(nil)
	if sg.ID() != "" || sg.URI() != "" || sg.Name() != "" {
		t.Error("fromResponse(nil) should be a no-op")
	}
	if sg.Raw() != nil {
		t.Error("Raw() should be nil before hydration")
	}

	// empty response — accessors must not panic; zero values expected
	sg2 := &SecurityGroup{}
	sg2.fromResponse(&types.SecurityGroupResponse{})
	if sg2.ID() != "" || sg2.URI() != "" || sg2.State() != "" {
		t.Error("empty response should yield zero accessor values")
	}
	// Default is plain bool: response's zero value is false → our *bool points to false
	if sg2.Default() != false {
		t.Errorf("Default() from zero response = %v", sg2.Default())
	}
}

func TestSecurityGroup_FromResponseURIBackfill(t *testing.T) {
	uri := "/projects/p2/network/vpcs/v2/security-groups/sg-2"
	id := "sg-2"
	name := "uri-sg"
	resp := &types.SecurityGroupResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:   &id,
			URI:  &uri,
			Name: &name,
			// ProjectResponseMetadata intentionally nil
		},
	}
	sg := &SecurityGroup{}
	sg.fromResponse(resp)

	if sg.ProjectID() != "p2" {
		t.Errorf("ProjectID() via URI fallback = %q", sg.ProjectID())
	}
	if sg.VPCID() != "v2" {
		t.Errorf("VPCID() via URI fallback = %q", sg.VPCID())
	}
}

// --------------------------------------------------------------------------
// Ref + ancestor ID satisfaction (runtime)
// --------------------------------------------------------------------------

func TestSecurityGroup_RefSatisfaction(t *testing.T) {
	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg-99", "n", "/projects/p99/network/vpcs/v99/security-groups/sg-99", "p99"))

	// withSecurityGroupID typed path
	sid, ok := extractID(sg, func(r Ref) (string, bool) {
		if w, ok := r.(withSecurityGroupID); ok {
			return w.SecurityGroupID(), true
		}
		return "", false
	}, "security-groups")
	if !ok || sid != "sg-99" {
		t.Errorf("extractID via withSecurityGroupID = (%q, %v)", sid, ok)
	}

	// withVPCID typed path
	vid, ok := extractID(sg, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid != "v99" {
		t.Errorf("extractID via withVPCID = (%q, %v)", vid, ok)
	}

	// withProjectID typed path
	pid, ok := extractID(sg, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid != "p99" {
		t.Errorf("extractID via withProjectID = (%q, %v)", pid, ok)
	}
}

// --------------------------------------------------------------------------
// securityGroupIDsFromRef helper
// --------------------------------------------------------------------------

func TestSecurityGroupIDsFromRef_TypedRef(t *testing.T) {
	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg-1", "n", "/projects/p/network/vpcs/v/security-groups/sg-1", "p"))
	pid, vid, sgid, err := securityGroupIDsFromRef(sg)
	if err != nil || pid != "p" || vid != "v" || sgid != "sg-1" {
		t.Errorf("securityGroupIDsFromRef typed = (%q, %q, %q, %v)", pid, vid, sgid, err)
	}
}

func TestSecurityGroupIDsFromRef_URIRef_APIForm(t *testing.T) {
	ref := URI("/projects/p/providers/Aruba.Network/vpcs/v/security-groups/sg-1")
	pid, vid, sgid, err := securityGroupIDsFromRef(ref)
	if err != nil || pid != "p" || vid != "v" || sgid != "sg-1" {
		t.Errorf("securityGroupIDsFromRef API form = (%q, %q, %q, %v)", pid, vid, sgid, err)
	}
}

func TestSecurityGroupIDsFromRef_URIRef_NetworkForm(t *testing.T) {
	ref := URI("/projects/p/network/vpcs/v/security-groups/sg-1")
	pid, vid, sgid, err := securityGroupIDsFromRef(ref)
	if err != nil || pid != "p" || vid != "v" || sgid != "sg-1" {
		t.Errorf("securityGroupIDsFromRef network form = (%q, %q, %q, %v)", pid, vid, sgid, err)
	}
}

func TestSecurityGroupIDsFromRef_BadURI_MissingSG(t *testing.T) {
	_, _, _, err := securityGroupIDsFromRef(URI("/projects/p/network/vpcs/v"))
	if err == nil {
		t.Error("expected error for URI without /security-groups/<id>")
	}
}

func TestSecurityGroupIDsFromRef_BadURI_MissingAll(t *testing.T) {
	_, _, _, err := securityGroupIDsFromRef(URI("/something/else"))
	if err == nil {
		t.Error("expected error for totally invalid URI")
	}
}

// --------------------------------------------------------------------------
// securityGroupsClientAdapter — CRUD integration tests
// --------------------------------------------------------------------------

func buildSecurityGroupTestAdapter(t *testing.T, handler http.HandlerFunc) *securityGroupsClientAdapter {
	t.Helper()
	server := testutil.NewMockServer(t, handler)
	return newSecurityGroupsClientAdapter(testutil.NewClient(t, server.URL))
}

const securityGroupSuccessBody = `{` +
	`"metadata":{"id":"sg-1","name":"my-sg","uri":"/projects/p/network/vpcs/v/security-groups/sg-1","project":{"id":"p"}},` +
	`"properties":{"default":false},` +
	`"status":{"state":"Active"}}`

// withVPCActiveRouteForSG wraps a handler so VPC GET requests are answered with
// an "Active" state. The internal security-group client calls waitForVPCActive
// before Create, which polls the VPC endpoint; without this, the test would time out.
func withVPCActiveRouteForSG(sgHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/securitygroups") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, activeVPCBody)
			return
		}
		sgHandler(w, r)
	}
}

func TestSecurityGroupsClientAdapter_Create_Success(t *testing.T) {
	var gotBody types.SecurityGroupRequest
	adapter := buildSecurityGroupTestAdapter(t, withVPCActiveRouteForSG(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, securityGroupSuccessBody)
	}))

	vpc := &VPC{}
	vpc.fromResponse(vpcTestResponse("v", "my-vpc", "/projects/p/providers/Aruba.Network/vpcs/v", "p"))

	sg := NewSecurityGroup().
		IntoVPC(vpc).
		WithName("my-sg").
		WithDefault(false)

	result, err := adapter.Create(context.Background(), sg)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.ID() != "sg-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.Name() != "my-sg" {
		t.Errorf("Name() = %q", result.Name())
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if gotBody.Metadata.Name != "my-sg" {
		t.Errorf("request Name = %q", gotBody.Metadata.Name)
	}
	if gotBody.Properties.Default == nil || *gotBody.Properties.Default != false {
		t.Errorf("request Default = %v", gotBody.Properties.Default)
	}
}

func TestSecurityGroupsClientAdapter_Create_NoVPC(t *testing.T) {
	callCount := 0
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})

	_, err := adapter.Create(context.Background(), NewSecurityGroup().WithName("x"))
	if err == nil {
		t.Fatal("expected error when security group has no VPC")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without VPC")
	}
}

func TestSecurityGroupsClientAdapter_Create_MetadataValidationError(t *testing.T) {
	adapter := buildSecurityGroupTestAdapter(t, withVPCActiveRouteForSG(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Missing "id" field — triggers MetadataValidationError
		fmt.Fprint(w, `{"metadata":{"name":"sg","uri":"/projects/p/network/vpcs/v/security-groups/x"},"properties":{},"status":{}}`)
	}))

	vpc := &VPC{}
	vpc.fromResponse(vpcTestResponse("v", "my-vpc", "/projects/p/providers/Aruba.Network/vpcs/v", "p"))

	sg := NewSecurityGroup().IntoVPC(vpc).WithName("sg")
	result, err := adapter.Create(context.Background(), sg)
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

func TestSecurityGroupsClientAdapter_Create_NonTwoXX(t *testing.T) {
	adapter := buildSecurityGroupTestAdapter(t, withVPCActiveRouteForSG(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Validation Failed", "name is required", 422))
	}))

	vpc := &VPC{}
	vpc.fromResponse(vpcTestResponse("v", "my-vpc", "/projects/p/providers/Aruba.Network/vpcs/v", "p"))

	sg := NewSecurityGroup().IntoVPC(vpc)
	result, err := adapter.Create(context.Background(), sg)
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

func TestSecurityGroupsClientAdapter_Get_URIRef(t *testing.T) {
	var capturedPath string
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, securityGroupSuccessBody)
	})

	ref := URI("/projects/p/network/vpcs/v/security-groups/sg-1")
	result, err := adapter.Get(context.Background(), ref)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "sg-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", result.ProjectID())
	}
	if result.VPCID() != "v" {
		t.Errorf("VPCID() = %q", result.VPCID())
	}
	if result.StatusCode() != http.StatusOK {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if !strings.Contains(capturedPath, "securitygroups") {
		t.Errorf("path = %q, expected securitygroups segment", capturedPath)
	}
}

func TestSecurityGroupsClientAdapter_Get_TypedRef(t *testing.T) {
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, securityGroupSuccessBody)
	})

	existing := &SecurityGroup{}
	existing.fromResponse(securityGroupTestResponse("sg-1", "n", "/projects/p/network/vpcs/v/security-groups/sg-1", "p"))

	result, err := adapter.Get(context.Background(), existing)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "sg-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestSecurityGroupsClientAdapter_Update_Success(t *testing.T) {
	var capturedBody types.SecurityGroupRequest
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"metadata":{"id":"sg-1","name":"renamed","uri":"/projects/p/network/vpcs/v/security-groups/sg-1","project":{"id":"p"}},"properties":{"default":true},"status":{}}`)
	})

	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg-1", "orig", "/projects/p/network/vpcs/v/security-groups/sg-1", "p"))
	sg.WithName("renamed").WithDefault(true)

	result, err := adapter.Update(context.Background(), sg)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if result.Name() != "renamed" {
		t.Errorf("Name() = %q", result.Name())
	}
	if capturedBody.Metadata.Name != "renamed" {
		t.Errorf("request Name = %q", capturedBody.Metadata.Name)
	}
}

func TestSecurityGroupsClientAdapter_Update_NoID(t *testing.T) {
	callCount := 0
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	sg := NewSecurityGroup().IntoVPC(URI("/projects/p/network/vpcs/v")).WithName("x")
	_, err := adapter.Update(context.Background(), sg)
	if err == nil {
		t.Fatal("expected error when security group has no ID")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made when ID is missing")
	}
}

func TestSecurityGroupsClientAdapter_Update_NoVPC(t *testing.T) {
	callCount := 0
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	sg := &SecurityGroup{}
	id := "sg-1"
	sg.fromResponse(&types.SecurityGroupResponse{
		Metadata: types.ResourceMetadataResponse{
			ID: &id,
		},
	})

	_, err := adapter.Update(context.Background(), sg)
	if err == nil {
		t.Fatal("expected error when security group has no VPC")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without VPC")
	}
}

func TestSecurityGroupsClientAdapter_Delete_Success(t *testing.T) {
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/network/vpcs/v/security-groups/sg-1"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestSecurityGroupsClientAdapter_Delete_NonTwoXX(t *testing.T) {
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "security group not found", 404))
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/network/vpcs/v/security-groups/missing"))
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

func TestSecurityGroupsClientAdapter_List_TwoItems(t *testing.T) {
	adapter := buildSecurityGroupTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total":2,"self":"","prev":"","next":"","first":"","last":"","values":[`+
			`{"metadata":{"id":"sg-1","name":"sg-a","uri":"/projects/p/network/vpcs/v/security-groups/sg-1","project":{"id":"p"}},"properties":{"default":false},"status":{"state":"Active"}},`+
			`{"metadata":{"id":"sg-2","name":"sg-b","uri":"/projects/p/network/vpcs/v/security-groups/sg-2","project":{"id":"p"}},"properties":{"default":true},"status":{"state":"Inactive"}}`+
			`]}`)
	})

	list, err := adapter.List(context.Background(), URI("/projects/p/network/vpcs/v"))
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
	if items[0].ID() != "sg-1" || items[0].Name() != "sg-a" {
		t.Errorf("items[0] = {%q, %q}", items[0].ID(), items[0].Name())
	}
	if items[0].Default() != false {
		t.Errorf("items[0].Default() = %v", items[0].Default())
	}
	if items[1].ID() != "sg-2" || items[1].State() != "Inactive" {
		t.Errorf("items[1] ID=%q State=%q", items[1].ID(), items[1].State())
	}
	if items[0].VPCID() != "v" {
		t.Errorf("items[0].VPCID() = %q", items[0].VPCID())
	}
	if items[0].ProjectID() != "p" {
		t.Errorf("items[0].ProjectID() = %q", items[0].ProjectID())
	}
}
