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

var _ Ref = (*VPCPeeringRoute)(nil)

// --------------------------------------------------------------------------
// Fluent setters
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_FluentSetters(t *testing.T) {
	parent := &VPCPeering{}
	parent.fromResponse(vpcPeeringTestResponse("peer-1", "my-peering",
		"/projects/p1/providers/Aruba.Network/vpcs/v1/vpcPeerings/peer-1", "p1"))

	r := NewVPCPeeringRoute().
		IntoVPCPeering(parent).
		WithName("my-route").
		AddTag("route").
		AddTag("billing").
		AddTag("route"). // dedupe
		InRegion("ITBG-Bergamo").
		WithLocalCIDR("10.0.0.0/24").
		WithRemoteCIDR("192.168.0.0/24").
		WithBillingPeriod("Hourly")

	if r.Name() != "my-route" {
		t.Errorf("Name() = %q", r.Name())
	}
	if tags := r.Tags(); len(tags) != 2 || tags[0] != "route" || tags[1] != "billing" {
		t.Errorf("Tags() = %v", tags)
	}
	if r.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", r.Region())
	}
	if r.LocalCIDR() != "10.0.0.0/24" {
		t.Errorf("LocalCIDR() = %q", r.LocalCIDR())
	}
	if r.RemoteCIDR() != "192.168.0.0/24" {
		t.Errorf("RemoteCIDR() = %q", r.RemoteCIDR())
	}
	if r.BillingPeriod() != "Hourly" {
		t.Errorf("BillingPeriod() = %q", r.BillingPeriod())
	}
	if r.VPCPeeringID() != "peer-1" {
		t.Errorf("VPCPeeringID() = %q", r.VPCPeeringID())
	}
	if r.VPCID() != "v1" {
		t.Errorf("VPCID() = %q", r.VPCID())
	}
	if r.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", r.ProjectID())
	}
	if r.Err() != nil {
		t.Errorf("Err() = %v", r.Err())
	}

	r.RemoveTag("route")
	if tags := r.Tags(); len(tags) != 1 || tags[0] != "billing" {
		t.Errorf("after RemoveTag Tags() = %v", tags)
	}

	r.ReplaceTags("x", "y")
	if tags := r.Tags(); len(tags) != 2 || tags[0] != "x" || tags[1] != "y" {
		t.Errorf("after ReplaceTags Tags() = %v", tags)
	}
}

// --------------------------------------------------------------------------
// IntoVPCPeering — typed Ref
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_IntoVPCPeering_TypedRef(t *testing.T) {
	parent := &VPCPeering{}
	parent.fromResponse(vpcPeeringTestResponse("peer-1", "my-peering",
		"/projects/p1/providers/Aruba.Network/vpcs/v1/vpcPeerings/peer-1", "p1"))

	r := NewVPCPeeringRoute().IntoVPCPeering(parent)

	if r.VPCPeeringID() != "peer-1" {
		t.Errorf("VPCPeeringID() = %q", r.VPCPeeringID())
	}
	if r.VPCID() != "v1" {
		t.Errorf("VPCID() = %q", r.VPCID())
	}
	if r.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", r.ProjectID())
	}
	if r.Err() != nil {
		t.Errorf("Err() = %v", r.Err())
	}
}

// --------------------------------------------------------------------------
// IntoVPCPeering — URI Ref (lowercase/mixin form)
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_IntoVPCPeering_URIRef_LowerCase(t *testing.T) {
	r := NewVPCPeeringRoute().IntoVPCPeering(URI("/projects/p/network/vpcs/v/peerings/peer"))

	if r.VPCPeeringID() != "peer" {
		t.Errorf("VPCPeeringID() = %q", r.VPCPeeringID())
	}
	if r.VPCID() != "v" {
		t.Errorf("VPCID() = %q", r.VPCID())
	}
	if r.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", r.ProjectID())
	}
	if r.Err() != nil {
		t.Errorf("Err() = %v", r.Err())
	}
}

// --------------------------------------------------------------------------
// IntoVPCPeering — URI Ref (camelCase/production form — validates mixin extension)
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_IntoVPCPeering_URIRef_CamelCase(t *testing.T) {
	r := NewVPCPeeringRoute().IntoVPCPeering(
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer"))

	if r.VPCPeeringID() != "peer" {
		t.Errorf("VPCPeeringID() = %q", r.VPCPeeringID())
	}
	if r.VPCID() != "v" {
		t.Errorf("VPCID() = %q", r.VPCID())
	}
	if r.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", r.ProjectID())
	}
	if r.Err() != nil {
		t.Errorf("Err() = %v", r.Err())
	}
}

// --------------------------------------------------------------------------
// IntoVPCPeering — bad Ref
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_IntoVPCPeering_BadRef(t *testing.T) {
	r := NewVPCPeeringRoute().IntoVPCPeering(URI("/garbage"))
	if r.Err() == nil {
		t.Error("expected Err() != nil for unresolvable Ref, got nil")
	}
}

// --------------------------------------------------------------------------
// toRequest round-trip
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_ToRequestRoundTrip(t *testing.T) {
	r := NewVPCPeeringRoute().
		WithName("my-route").
		AddTag("t1").
		AddTag("t2").
		InRegion("ITBG-Bergamo").
		WithLocalCIDR("10.0.0.0/24").
		WithRemoteCIDR("192.168.0.0/24").
		WithBillingPeriod("Hourly")

	req := r.RawRequest()

	if req.Metadata.Name != "my-route" {
		t.Errorf("Metadata.Name = %q", req.Metadata.Name)
	}
	if len(req.Metadata.Tags) != 2 {
		t.Errorf("Metadata.Tags = %v", req.Metadata.Tags)
	}
	if req.Metadata.Location.Value != "ITBG-Bergamo" {
		t.Errorf("Metadata.Location.Value = %q", req.Metadata.Location.Value)
	}
	if req.Properties.LocalNetworkAddress != "10.0.0.0/24" {
		t.Errorf("Properties.LocalNetworkAddress = %q", req.Properties.LocalNetworkAddress)
	}
	if req.Properties.RemoteNetworkAddress != "192.168.0.0/24" {
		t.Errorf("Properties.RemoteNetworkAddress = %q", req.Properties.RemoteNetworkAddress)
	}
	if req.Properties.BillingPeriod == nil || *req.Properties.BillingPeriod != "Hourly" {
		t.Errorf("Properties.BillingPeriod = %v", req.Properties.BillingPeriod)
	}

	// Unset CIDRs must produce empty strings.
	r2 := NewVPCPeeringRoute().WithName("bare")
	req2 := r2.RawRequest()
	if req2.Properties.LocalNetworkAddress != "" {
		t.Errorf("LocalNetworkAddress should be empty when not set, got %q", req2.Properties.LocalNetworkAddress)
	}
	if req2.Properties.RemoteNetworkAddress != "" {
		t.Errorf("RemoteNetworkAddress should be empty when not set, got %q", req2.Properties.RemoteNetworkAddress)
	}
}

// --------------------------------------------------------------------------
// BillingPlan always emitted (value type, no omitempty)
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_ToRequest_BillingPeriodAlwaysEmitted(t *testing.T) {
	r := NewVPCPeeringRoute().WithName("bare")
	req := r.RawRequest()
	// BillingPeriod is always emitted — defaults to Hour when not explicitly set.
	if req.Properties.BillingPeriod == nil || *req.Properties.BillingPeriod != BillingPeriodHour {
		t.Errorf("BillingPeriod should default to Hour when not set, got %v", req.Properties.BillingPeriod)
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration
// --------------------------------------------------------------------------

func vpcPeeringRouteTestResponse(id, name, uri, projectID string) *types.VPCPeeringRouteResponse {
	state := "Active"
	loc := &types.LocationResponse{Value: "ITBG-Bergamo"}
	return &types.VPCPeeringRouteResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:               &id,
			URI:              &uri,
			Name:             &name,
			Tags:             []string{"route-tag"},
			LocationResponse: loc,
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: projectID,
			},
		},
		Properties: types.VPCPeeringRoutePropertiesResponse{
			LocalNetworkAddress:  "10.0.0.0/24",
			RemoteNetworkAddress: "192.168.0.0/24",
			BillingPeriod:        func() *types.BillingPeriod { v := types.BillingPeriod("Hourly"); return &v }(),
		},
		Status: types.ResourceStatus{
			State: &state,
		},
	}
}

func TestVPCPeeringRoute_FromResponseHydration(t *testing.T) {
	r := &VPCPeeringRoute{}
	resp := vpcPeeringRouteTestResponse("route-1", "my-route",
		"/projects/p1/providers/Aruba.Network/vpcs/v1/vpcPeerings/peer-1/vpcPeeringRoutes/route-1", "p1")
	r.fromResponse(resp)

	if r.ID() != "route-1" {
		t.Errorf("ID() = %q", r.ID())
	}
	if r.URI() != "/projects/p1/providers/Aruba.Network/vpcs/v1/vpcPeerings/peer-1/vpcPeeringRoutes/route-1" {
		t.Errorf("URI() = %q", r.URI())
	}
	if r.VPCPeeringRouteID() != "route-1" {
		t.Errorf("VPCPeeringRouteID() = %q", r.VPCPeeringRouteID())
	}
	if r.Name() != "my-route" {
		t.Errorf("Name() = %q", r.Name())
	}
	if tags := r.Tags(); len(tags) != 1 || tags[0] != "route-tag" {
		t.Errorf("Tags() = %v", tags)
	}
	if r.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", r.Region())
	}
	if r.State() != "Active" {
		t.Errorf("State() = %q", r.State())
	}
	if r.LocalCIDR() != "10.0.0.0/24" {
		t.Errorf("LocalCIDR() = %q", r.LocalCIDR())
	}
	if r.RemoteCIDR() != "192.168.0.0/24" {
		t.Errorf("RemoteCIDR() = %q", r.RemoteCIDR())
	}
	if r.BillingPeriod() != "Hourly" {
		t.Errorf("BillingPeriod() = %q", r.BillingPeriod())
	}
	if r.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", r.ProjectID())
	}
	if r.VPCID() != "v1" {
		t.Errorf("VPCID() via URI fallback = %q", r.VPCID())
	}
	if r.VPCPeeringID() != "peer-1" {
		t.Errorf("VPCPeeringID() via URI fallback = %q", r.VPCPeeringID())
	}
	if r.Raw() != resp {
		t.Error("Raw() should return the hydrated response pointer")
	}
}

func TestVPCPeeringRoute_FromResponsePartial(t *testing.T) {
	// nil response is a no-op
	r := &VPCPeeringRoute{}
	r.fromResponse(nil)
	if r.ID() != "" || r.URI() != "" || r.Name() != "" {
		t.Error("fromResponse(nil) should be a no-op")
	}
	if r.Raw() != nil {
		t.Error("Raw() should be nil before hydration")
	}

	// empty response — accessors must not panic; zero values expected
	r2 := &VPCPeeringRoute{}
	r2.fromResponse(&types.VPCPeeringRouteResponse{})
	if r2.ID() != "" || r2.URI() != "" || r2.State() != "" {
		t.Error("empty response should yield zero accessor values")
	}
	if r2.LocalCIDR() != "" || r2.RemoteCIDR() != "" || r2.BillingPeriod() != "" {
		t.Error("empty response should yield zero CIDR/billing values")
	}
}

func TestVPCPeeringRoute_FromResponseURIBackfill_HyphenForm(t *testing.T) {
	uri := "/projects/p2/network/vpcs/v2/peerings/peer-2/vpc-peering-routes/route-2"
	id := "route-2"
	name := "uri-route"
	resp := &types.VPCPeeringRouteResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:   &id,
			URI:  &uri,
			Name: &name,
			// ProjectResponseMetadata intentionally nil
		},
	}
	r := &VPCPeeringRoute{}
	r.fromResponse(resp)

	if r.ProjectID() != "p2" {
		t.Errorf("ProjectID() via URI fallback = %q", r.ProjectID())
	}
	if r.VPCID() != "v2" {
		t.Errorf("VPCID() via URI fallback = %q", r.VPCID())
	}
	if r.VPCPeeringID() != "peer-2" {
		t.Errorf("VPCPeeringID() via URI fallback = %q", r.VPCPeeringID())
	}
}

// --------------------------------------------------------------------------
// Ref + ancestor ID satisfaction (runtime)
// --------------------------------------------------------------------------

func TestVPCPeeringRoute_RefSatisfaction(t *testing.T) {
	r := &VPCPeeringRoute{}
	r.fromResponse(vpcPeeringRouteTestResponse("route-99", "n",
		"/projects/p99/providers/Aruba.Network/vpcs/v99/vpcPeerings/peer-99/vpcPeeringRoutes/route-99", "p99"))

	// withVPCPeeringRouteID typed path
	rid, ok := extractID(r, func(ref Ref) (string, bool) {
		if w, ok := ref.(withVPCPeeringRouteID); ok {
			return w.VPCPeeringRouteID(), true
		}
		return "", false
	}, "vpc-peering-routes")
	if !ok || rid != "route-99" {
		t.Errorf("extractID via withVPCPeeringRouteID = (%q, %v)", rid, ok)
	}

	// withVPCPeeringID typed path
	pid, ok := extractID(r, func(ref Ref) (string, bool) {
		if w, ok := ref.(withVPCPeeringID); ok {
			return w.VPCPeeringID(), true
		}
		return "", false
	}, "vpc-peerings")
	if !ok || pid != "peer-99" {
		t.Errorf("extractID via withVPCPeeringID = (%q, %v)", pid, ok)
	}

	// withVPCID typed path
	vid, ok := extractID(r, func(ref Ref) (string, bool) {
		if w, ok := ref.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid != "v99" {
		t.Errorf("extractID via withVPCID = (%q, %v)", vid, ok)
	}

	// withProjectID typed path
	projID, ok := extractID(r, func(ref Ref) (string, bool) {
		if w, ok := ref.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || projID != "p99" {
		t.Errorf("extractID via withProjectID = (%q, %v)", projID, ok)
	}
}

// --------------------------------------------------------------------------
// vpcPeeringRouteIDsFromRef helper
// --------------------------------------------------------------------------

func TestVPCPeeringRouteIDsFromRef_TypedRef(t *testing.T) {
	r := &VPCPeeringRoute{}
	r.fromResponse(vpcPeeringRouteTestResponse("route-1", "n",
		"/projects/p/network/vpcs/v/peerings/peer-1/vpc-peering-routes/route-1", "p"))
	projID, vid, peerid, rid, err := vpcPeeringRouteIDsFromRef(r)
	if err != nil || projID != "p" || vid != "v" || peerid != "peer-1" || rid != "route-1" {
		t.Errorf("vpcPeeringRouteIDsFromRef typed = (%q, %q, %q, %q, %v)", projID, vid, peerid, rid, err)
	}
}

func TestVPCPeeringRouteIDsFromRef_URIRef_CamelCase(t *testing.T) {
	ref := URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1")
	projID, vid, peerid, rid, err := vpcPeeringRouteIDsFromRef(ref)
	if err != nil || projID != "p" || vid != "v" || peerid != "peer-1" || rid != "route-1" {
		t.Errorf("vpcPeeringRouteIDsFromRef camelCase = (%q, %q, %q, %q, %v)", projID, vid, peerid, rid, err)
	}
}

func TestVPCPeeringRouteIDsFromRef_URIRef_LowerCase(t *testing.T) {
	ref := URI("/projects/p/network/vpcs/v/peerings/peer-1/vpc-peering-routes/route-1")
	projID, vid, peerid, rid, err := vpcPeeringRouteIDsFromRef(ref)
	if err != nil || projID != "p" || vid != "v" || peerid != "peer-1" || rid != "route-1" {
		t.Errorf("vpcPeeringRouteIDsFromRef lowercase = (%q, %q, %q, %q, %v)", projID, vid, peerid, rid, err)
	}
}

func TestVPCPeeringRouteIDsFromRef_BadURI_MissingRoute(t *testing.T) {
	_, _, _, _, err := vpcPeeringRouteIDsFromRef(
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1"))
	if err == nil {
		t.Error("expected error for URI without route segment")
	}
}

func TestVPCPeeringRouteIDsFromRef_BadURI_MissingPeering(t *testing.T) {
	_, _, _, _, err := vpcPeeringRouteIDsFromRef(
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeeringRoutes/route-1"))
	if err == nil {
		t.Error("expected error for URI with route but no peering segment")
	}
}

func TestVPCPeeringRouteIDsFromRef_BadURI_MissingAll(t *testing.T) {
	_, _, _, _, err := vpcPeeringRouteIDsFromRef(URI("/something/else"))
	if err == nil {
		t.Error("expected error for totally invalid URI")
	}
}

// --------------------------------------------------------------------------
// vpcPeeringRoutesClientAdapter — CRUD integration tests
// --------------------------------------------------------------------------

func buildVPCPeeringRouteTestAdapter(t *testing.T, handler http.HandlerFunc) *vpcPeeringRoutesClientAdapter {
	t.Helper()
	server := testutil.NewMockServer(t, handler)
	return newVPCPeeringRoutesClientAdapter(testutil.NewClient(t, server.URL))
}

const vpcPeeringRouteSuccessBody = `{` +
	`"metadata":{"id":"route-1","name":"my-route","uri":"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1","project":{"id":"p"}},` +
	`"properties":{"localNetworkAddress":"10.0.0.0/24","remoteNetworkAddress":"192.168.0.0/24","billingPeriod":"Hourly"},` +
	`"status":{"state":"Active"}}`

func TestVPCPeeringRoutesClientAdapter_Create_Success(t *testing.T) {
	var gotBody types.VPCPeeringRouteRequest
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, vpcPeeringRouteSuccessBody)
	})

	peering := &VPCPeering{}
	peering.fromResponse(vpcPeeringTestResponse("peer-1", "my-peering",
		"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1", "p"))

	route := NewVPCPeeringRoute().
		IntoVPCPeering(peering).
		WithName("my-route").
		InRegion("ITBG-Bergamo").
		WithLocalCIDR("10.0.0.0/24").
		WithRemoteCIDR("192.168.0.0/24").
		WithBillingPeriod("Hourly")

	result, err := adapter.Create(context.Background(), route)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.ID() != "route-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.Name() != "my-route" {
		t.Errorf("Name() = %q", result.Name())
	}
	if result.LocalCIDR() != "10.0.0.0/24" {
		t.Errorf("LocalCIDR() = %q", result.LocalCIDR())
	}
	if result.BillingPeriod() != "Hourly" {
		t.Errorf("BillingPeriod() = %q", result.BillingPeriod())
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if gotBody.Metadata.Name != "my-route" {
		t.Errorf("request Name = %q", gotBody.Metadata.Name)
	}
	if gotBody.Metadata.Location.Value != "ITBG-Bergamo" {
		t.Errorf("request Location = %q", gotBody.Metadata.Location.Value)
	}
	if gotBody.Properties.LocalNetworkAddress != "10.0.0.0/24" {
		t.Errorf("request LocalNetworkAddress = %q", gotBody.Properties.LocalNetworkAddress)
	}
	if gotBody.Properties.BillingPeriod == nil || *gotBody.Properties.BillingPeriod != "Hourly" {
		t.Errorf("request BillingPeriod = %v", gotBody.Properties.BillingPeriod)
	}
}

func TestVPCPeeringRoutesClientAdapter_Create_NoPeering(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})

	_, err := adapter.Create(context.Background(), NewVPCPeeringRoute().WithName("x"))
	if err == nil {
		t.Fatal("expected error when route has no parent peering")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without parent peering")
	}
}

func TestVPCPeeringRoutesClientAdapter_Create_MetadataValidationError(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Missing "id" field — triggers MetadataValidationError
		fmt.Fprint(w, `{"metadata":{"name":"route","uri":"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/x"},"properties":{},"status":{}}`)
	})

	peering := &VPCPeering{}
	peering.fromResponse(vpcPeeringTestResponse("peer-1", "my-peering",
		"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1", "p"))

	route := NewVPCPeeringRoute().IntoVPCPeering(peering).WithName("route")
	result, err := adapter.Create(context.Background(), route)
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

func TestVPCPeeringRoutesClientAdapter_Create_NonTwoXX(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Validation Failed", "name is required", 422))
	})

	peering := &VPCPeering{}
	peering.fromResponse(vpcPeeringTestResponse("peer-1", "my-peering",
		"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1", "p"))

	route := NewVPCPeeringRoute().IntoVPCPeering(peering)
	result, err := adapter.Create(context.Background(), route)
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

func TestVPCPeeringRoutesClientAdapter_Get_URIRef_CamelCase(t *testing.T) {
	var capturedPath string
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, vpcPeeringRouteSuccessBody)
	})

	ref := URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1")
	result, err := adapter.Get(context.Background(), ref)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "route-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", result.ProjectID())
	}
	if result.VPCID() != "v" {
		t.Errorf("VPCID() = %q", result.VPCID())
	}
	if result.VPCPeeringID() != "peer-1" {
		t.Errorf("VPCPeeringID() = %q", result.VPCPeeringID())
	}
	if result.StatusCode() != http.StatusOK {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if !strings.Contains(capturedPath, "vpcPeeringRoutes") {
		t.Errorf("path = %q, expected vpcPeeringRoutes segment", capturedPath)
	}
}

func TestVPCPeeringRoutesClientAdapter_Get_TypedRef(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, vpcPeeringRouteSuccessBody)
	})

	existing := &VPCPeeringRoute{}
	existing.fromResponse(vpcPeeringRouteTestResponse("route-1", "n",
		"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1", "p"))

	result, err := adapter.Get(context.Background(), existing)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "route-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestVPCPeeringRoutesClientAdapter_Update_Success(t *testing.T) {
	var capturedBody types.VPCPeeringRouteRequest
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"metadata":{"id":"route-1","name":"renamed","uri":"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1","project":{"id":"p"}},"properties":{},"status":{}}`)
	})

	r := &VPCPeeringRoute{}
	r.fromResponse(vpcPeeringRouteTestResponse("route-1", "orig",
		"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1", "p"))
	r.WithName("renamed")

	result, err := adapter.Update(context.Background(), r)
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

func TestVPCPeeringRoutesClientAdapter_Update_NoID(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	r := NewVPCPeeringRoute().
		IntoVPCPeering(URI("/projects/p/network/vpcs/v/peerings/peer-1")).
		WithName("x")
	_, err := adapter.Update(context.Background(), r)
	if err == nil {
		t.Fatal("expected error when route has no ID")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made when ID is missing")
	}
}

func TestVPCPeeringRoutesClientAdapter_Update_NoPeering(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	r := &VPCPeeringRoute{}
	id := "route-1"
	r.fromResponse(&types.VPCPeeringRouteResponse{
		Metadata: types.ResourceMetadataResponse{
			ID: &id,
		},
	})

	_, err := adapter.Update(context.Background(), r)
	if err == nil {
		t.Fatal("expected error when route has no parent peering")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without parent peering")
	}
}

func TestVPCPeeringRoutesClientAdapter_Delete_BadRef(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})
	err := adapter.Delete(context.Background(), URI("/something/else"))
	if err == nil {
		t.Fatal("expected error for bad Ref")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made for bad Ref")
	}
}

func TestVPCPeeringRoutesClientAdapter_Delete_Success(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := adapter.Delete(context.Background(),
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestVPCPeeringRoutesClientAdapter_Delete_NonTwoXX(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "vpc peering route not found", 404))
	})

	err := adapter.Delete(context.Background(),
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/missing"))
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

// InRegion exercises the 0% branch.
func TestVPCPeeringRoute_InRegion(t *testing.T) {
	r := NewVPCPeeringRoute().
		AddTag("a").
		AddTag("b").
		RemoveTag("a").
		ReplaceTags("x", "y").
		InRegion("ITMI-Milano-1")

	if r.Region() != "ITMI-Milano-1" {
		t.Errorf("Region() = %q", r.Region())
	}
	if tags := r.Tags(); len(tags) != 2 || tags[0] != "x" || tags[1] != "y" {
		t.Errorf("Tags() = %v", tags)
	}
}

func TestVPCPeeringRoutesClientAdapter_Get_NonTwoXX(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "vpc peering route not found", 404))
	})

	ref := URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/missing")
	result, err := adapter.Get(context.Background(), ref)
	if err == nil {
		t.Fatal("expected error on 404")
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d", httpErr.StatusCode)
	}
	if result == nil {
		t.Fatal("result must be non-nil on non-2xx")
	}
}

func TestVPCPeeringRoutesClientAdapter_Update_NonTwoXX(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "vpc peering route not found", 404))
	})

	r := &VPCPeeringRoute{}
	r.fromResponse(vpcPeeringRouteTestResponse("route-1", "my-route",
		"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1", "p"))
	_, err := adapter.Update(context.Background(), r)
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
}

func TestVPCPeeringRoutesClientAdapter_List_NonTwoXX(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Forbidden", "access denied", 403))
	})

	_, err := adapter.List(context.Background(),
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1"))
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		t.Fatalf("expected *HTTPError, got %T: %v", err, err)
	}
}

func TestVPCPeeringRouteIDsFromRef_BadURI_MissingVPC(t *testing.T) {
	// URI has vpc-peering-routes+vpc-peerings but no vpcs segment
	_, _, _, _, err := vpcPeeringRouteIDsFromRef(URI("/projects/p/vpc-peerings/peer/vpc-peering-routes/route"))
	if err == nil {
		t.Error("expected error for URI without /vpcs/<id>")
	}
}

func TestVPCPeeringRouteIDsFromRef_BadURI_MissingProject(t *testing.T) {
	// URI has vpc-peering-routes+vpc-peerings+vpcs but no projects
	_, _, _, _, err := vpcPeeringRouteIDsFromRef(URI("/providers/Aruba.Network/vpcs/v/vpc-peerings/peer/vpc-peering-routes/route"))
	if err == nil {
		t.Error("expected error for URI without /projects/<id>")
	}
}

func TestVPCPeeringRoutesClientAdapter_Create_WithBuilderError(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})
	route := NewVPCPeeringRoute().IntoVPCPeering(URI("/garbage"))
	_, err := adapter.Create(context.Background(), route)
	if err == nil {
		t.Fatal("expected error for builder error")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made when builder has errors")
	}
}

func TestVPCPeeringRoutesClientAdapter_Get_BadRef(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})
	result, err := adapter.Get(context.Background(), URI("/something/else"))
	if err == nil {
		t.Fatal("expected error for bad Ref")
	}
	if result != nil {
		t.Error("result should be nil on bad Ref")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made for bad Ref")
	}
}

func TestVPCPeeringRoutesClientAdapter_Get_TransportError(t *testing.T) {
	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Error("server doesn't support hijacking")
			return
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	})
	adapter := newVPCPeeringRoutesClientAdapter(testutil.NewClient(t, server.URL))
	result, err := adapter.Get(context.Background(),
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpc-peerings/peer/vpc-peering-routes/route"))
	if err == nil {
		t.Fatal("expected transport error")
	}
	_ = result
}

func TestVPCPeeringRoutesClientAdapter_Update_WithBuilderError(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})
	route := NewVPCPeeringRoute().IntoVPCPeering(URI("/garbage"))
	_, err := adapter.Update(context.Background(), route)
	if err == nil {
		t.Fatal("expected error for builder error")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made when builder has errors")
	}
}

func TestVPCPeeringRoutesClientAdapter_Update_TransportError(t *testing.T) {
	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Error("server doesn't support hijacking")
			return
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	})
	adapter := newVPCPeeringRoutesClientAdapter(testutil.NewClient(t, server.URL))
	route := &VPCPeeringRoute{}
	route.fromResponse(vpcPeeringRouteTestResponse("route-1", "route-a",
		"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1", "p"))
	_, err := adapter.Update(context.Background(), route)
	if err == nil {
		t.Fatal("expected transport error")
	}
}

func TestVPCPeeringRoutesClientAdapter_Delete_TransportError(t *testing.T) {
	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Error("server doesn't support hijacking")
			return
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	})
	adapter := newVPCPeeringRoutesClientAdapter(testutil.NewClient(t, server.URL))
	err := adapter.Delete(context.Background(),
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpc-peerings/peer/vpc-peering-routes/route"))
	if err == nil {
		t.Fatal("expected transport error")
	}
}

func TestVPCPeeringRoutesClientAdapter_List_BadPeeringRef(t *testing.T) {
	callCount := 0
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})
	_, err := adapter.List(context.Background(), URI("/something/else"))
	if err == nil {
		t.Fatal("expected error for bad peering Ref")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made for bad peering Ref")
	}
}

func TestVPCPeeringRoutesClientAdapter_List_TransportError(t *testing.T) {
	server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Error("server doesn't support hijacking")
			return
		}
		conn, _, _ := hj.Hijack()
		conn.Close()
	})
	adapter := newVPCPeeringRoutesClientAdapter(testutil.NewClient(t, server.URL))
	_, err := adapter.List(context.Background(),
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpc-peerings/peer"))
	if err == nil {
		t.Fatal("expected transport error")
	}
}

func TestVPCPeeringRoutesClientAdapter_List_AncestorIDBackfill(t *testing.T) {
	// Items without ancestor IDs in metadata/URI: triggers vpcPeeringID/vpcID/projectID backfill
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total":1,"self":"","prev":"","next":"","first":"","last":"","values":[`+
			`{"metadata":{"id":"route-x","name":"route-x"},"properties":{"localNetworkAddress":"10.0.0.0/24","remoteNetworkAddress":"192.168.0.0/24"},"status":{}}`+
			`]}`)
	})

	list, err := adapter.List(context.Background(),
		URI("/projects/proj-x/providers/Aruba.Network/vpcs/vpc-x/vpc-peerings/peer-x"))
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	items := list.Items()
	if len(items) != 1 {
		t.Fatalf("Items() len = %d", len(items))
	}
	if items[0].ProjectID() != "proj-x" {
		t.Errorf("ProjectID() after backfill = %q, want %q", items[0].ProjectID(), "proj-x")
	}
	if items[0].VPCID() != "vpc-x" {
		t.Errorf("VPCID() after backfill = %q, want %q", items[0].VPCID(), "vpc-x")
	}
	if items[0].VPCPeeringID() != "peer-x" {
		t.Errorf("VPCPeeringID() after backfill = %q, want %q", items[0].VPCPeeringID(), "peer-x")
	}
}

func TestVPCPeeringRoutesClientAdapter_List_TwoItems(t *testing.T) {
	adapter := buildVPCPeeringRouteTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total":2,"self":"","prev":"","next":"","first":"","last":"","values":[`+
			`{"metadata":{"id":"route-1","name":"route-a","uri":"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1","project":{"id":"p"}},"properties":{"localNetworkAddress":"10.0.0.0/24","remoteNetworkAddress":"192.168.0.0/24"},"status":{"state":"Active"}},`+
			`{"metadata":{"id":"route-2","name":"route-b","uri":"/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-2","project":{"id":"p"}},"properties":{"localNetworkAddress":"10.1.0.0/24","remoteNetworkAddress":"192.168.1.0/24"},"status":{"state":"Inactive"}}`+
			`]}`)
	})

	list, err := adapter.List(context.Background(),
		URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1"))
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
	if items[0].ID() != "route-1" || items[0].Name() != "route-a" {
		t.Errorf("items[0] = {%q, %q}", items[0].ID(), items[0].Name())
	}
	if items[1].ID() != "route-2" || items[1].State() != "Inactive" {
		t.Errorf("items[1] ID=%q State=%q", items[1].ID(), items[1].State())
	}
	if items[0].VPCPeeringID() != "peer-1" {
		t.Errorf("items[0].VPCPeeringID() = %q", items[0].VPCPeeringID())
	}
	if items[0].VPCID() != "v" {
		t.Errorf("items[0].VPCID() = %q", items[0].VPCID())
	}
	if items[0].ProjectID() != "p" {
		t.Errorf("items[0].ProjectID() = %q", items[0].ProjectID())
	}
}

func TestVPCPeeringRoute_FromResponse_SetsTerminalStates(t *testing.T) {
	r := &VPCPeeringRoute{}
	state := "Active"
	r.fromResponse(&types.VPCPeeringRouteResponse{
		Status: types.ResourceStatus{State: &state},
	})
	if len(r.terminalStates) == 0 {
		t.Error("fromResponse should set terminalStates on the wrapper")
	}
	if !r.terminalStates["Active"] {
		t.Error("terminalStates[Active] should be true for VPCPeeringRoute")
	}
	if r.terminalStates["Error"] {
		t.Error("terminalStates[Error] should be false for VPCPeeringRoute")
	}
}

func TestVPCPeeringRoutesClientAdapter_Get_InjectsRefresh(t *testing.T) {
	server := testutil.NewMockServer(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, vpcPeeringRouteSuccessBody)
	})
	adapter := newVPCPeeringRoutesClientAdapter(testutil.NewClient(t, server.URL))
	route, err := adapter.Get(context.Background(), URI("/projects/p/providers/Aruba.Network/vpcs/v/vpcPeerings/peer-1/vpcPeeringRoutes/route-1"))
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if !refreshIsSet(&route.statusMixin) {
		t.Error("Get should inject a refresh callback into the returned VPCPeeringRoute")
	}
}
