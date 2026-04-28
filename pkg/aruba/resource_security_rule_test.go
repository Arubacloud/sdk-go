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

var _ Ref = (*SecurityRule)(nil)

// --------------------------------------------------------------------------
// Fluent setters
// --------------------------------------------------------------------------

func TestSecurityRule_FluentSetters(t *testing.T) {
	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg-1", "my-sg", "/projects/p1/network/vpcs/v1/security-groups/sg-1", "p1"))

	rule := NewSecurityRule().
		IntoSecurityGroup(sg).
		WithName("allow-ssh").
		AddTag("t1").
		AddTag("t2").
		AddTag("t1"). // dedupe
		InRegion("ITBG-Bergamo").
		WithDirection(string(types.RuleDirectionIngress)).
		WithProtocol("TCP").
		WithPort("22").
		WithTargetCIDR("0.0.0.0/0")

	if rule.Name() != "allow-ssh" {
		t.Errorf("Name() = %q", rule.Name())
	}
	if tags := rule.Tags(); len(tags) != 2 || tags[0] != "t1" || tags[1] != "t2" {
		t.Errorf("Tags() = %v", tags)
	}
	if rule.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", rule.Region())
	}
	if rule.Direction() != types.RuleDirectionIngress {
		t.Errorf("Direction() = %q", rule.Direction())
	}
	if rule.Protocol() != "TCP" {
		t.Errorf("Protocol() = %q", rule.Protocol())
	}
	if rule.Port() != "22" {
		t.Errorf("Port() = %q", rule.Port())
	}
	if rule.TargetKind() != types.EndpointTypeIP {
		t.Errorf("TargetKind() = %q", rule.TargetKind())
	}
	if rule.TargetValue() != "0.0.0.0/0" {
		t.Errorf("TargetValue() = %q", rule.TargetValue())
	}
	if rule.SecurityGroupID() != "sg-1" {
		t.Errorf("SecurityGroupID() = %q", rule.SecurityGroupID())
	}
	if rule.VPCID() != "v1" {
		t.Errorf("VPCID() = %q", rule.VPCID())
	}
	if rule.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", rule.ProjectID())
	}
	if rule.Err() != nil {
		t.Errorf("Err() = %v", rule.Err())
	}

	rule.RemoveTag("t1")
	if tags := rule.Tags(); len(tags) != 1 || tags[0] != "t2" {
		t.Errorf("after RemoveTag Tags() = %v", tags)
	}

	rule.ReplaceTags("x", "y")
	if tags := rule.Tags(); len(tags) != 2 || tags[0] != "x" || tags[1] != "y" {
		t.Errorf("after ReplaceTags Tags() = %v", tags)
	}
}

// --------------------------------------------------------------------------
// IntoSecurityGroup — typed Ref
// --------------------------------------------------------------------------

func TestSecurityRule_IntoSecurityGroup_TypedRef(t *testing.T) {
	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg-1", "my-sg", "/projects/p1/network/vpcs/v1/security-groups/sg-1", "p1"))

	rule := NewSecurityRule().IntoSecurityGroup(sg)

	if rule.SecurityGroupID() != "sg-1" {
		t.Errorf("SecurityGroupID() = %q", rule.SecurityGroupID())
	}
	if rule.VPCID() != "v1" {
		t.Errorf("VPCID() = %q", rule.VPCID())
	}
	if rule.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", rule.ProjectID())
	}
	if rule.Err() != nil {
		t.Errorf("Err() = %v", rule.Err())
	}
}

// --------------------------------------------------------------------------
// IntoSecurityGroup — URI Ref
// --------------------------------------------------------------------------

func TestSecurityRule_IntoSecurityGroup_URIRef(t *testing.T) {
	rule := NewSecurityRule().IntoSecurityGroup(URI("/projects/p/network/vpcs/v/security-groups/sg"))

	if rule.SecurityGroupID() != "sg" {
		t.Errorf("SecurityGroupID() = %q", rule.SecurityGroupID())
	}
	if rule.VPCID() != "v" {
		t.Errorf("VPCID() = %q", rule.VPCID())
	}
	if rule.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", rule.ProjectID())
	}
	if rule.Err() != nil {
		t.Errorf("Err() = %v", rule.Err())
	}
}

// --------------------------------------------------------------------------
// IntoSecurityGroup — bad Ref
// --------------------------------------------------------------------------

func TestSecurityRule_IntoSecurityGroup_BadRef(t *testing.T) {
	rule := NewSecurityRule().IntoSecurityGroup(URI("/garbage"))
	if rule.Err() == nil {
		t.Error("expected Err() != nil for unresolvable Ref, got nil")
	}
}

// --------------------------------------------------------------------------
// Mutually exclusive target setters
// --------------------------------------------------------------------------

func TestSecurityRule_TargetMutuallyExclusive_CIDRThenSG(t *testing.T) {
	otherSG := URI("/projects/p/network/vpcs/v/security-groups/sg-other")
	rule := NewSecurityRule().
		WithTargetCIDR("10.0.0.0/8").
		WithTargetSecurityGroup(otherSG)

	if rule.Err() == nil {
		t.Fatal("expected error after setting SecurityGroup target over CIDR target")
	}
	if !strings.Contains(rule.Err().Error(), "pick one") {
		t.Errorf("error message = %q, expected 'pick one'", rule.Err().Error())
	}
	// Target must remain the first (CIDR).
	if rule.TargetKind() != types.EndpointTypeIP {
		t.Errorf("TargetKind() = %q, expected IP (first setter wins)", rule.TargetKind())
	}
}

func TestSecurityRule_TargetMutuallyExclusive_SGThenCIDR(t *testing.T) {
	otherSG := URI("/projects/p/network/vpcs/v/security-groups/sg-other")
	rule := NewSecurityRule().
		WithTargetSecurityGroup(otherSG).
		WithTargetCIDR("10.0.0.0/8")

	if rule.Err() == nil {
		t.Fatal("expected error after setting CIDR target over SecurityGroup target")
	}
	if !strings.Contains(rule.Err().Error(), "pick one") {
		t.Errorf("error message = %q, expected 'pick one'", rule.Err().Error())
	}
	// Target must remain the first (SG).
	if rule.TargetKind() != types.EndpointTypeSecurityGroup {
		t.Errorf("TargetKind() = %q, expected SecurityGroup (first setter wins)", rule.TargetKind())
	}
}

func TestSecurityRule_TargetSecurityGroup_EmptyURI(t *testing.T) {
	rule := NewSecurityRule().WithTargetSecurityGroup(URI(""))
	if rule.Err() == nil {
		t.Error("expected error when target SecurityGroup Ref has empty URI")
	}
	if rule.target != nil {
		t.Error("target should remain nil after empty-URI error")
	}
}

// --------------------------------------------------------------------------
// toRequest round-trip
// --------------------------------------------------------------------------

func TestSecurityRule_ToRequestRoundTrip(t *testing.T) {
	rule := NewSecurityRule().
		WithName("allow-ssh").
		AddTag("t1").
		AddTag("t2").
		InRegion("ITBG-Bergamo").
		WithDirection(string(types.RuleDirectionIngress)).
		WithProtocol("TCP").
		WithPort("22").
		WithTargetCIDR("0.0.0.0/0")

	req := rule.RawRequest()

	if req.Metadata.Name != "allow-ssh" {
		t.Errorf("Metadata.Name = %q", req.Metadata.Name)
	}
	if len(req.Metadata.Tags) != 2 {
		t.Errorf("Metadata.Tags = %v", req.Metadata.Tags)
	}
	if req.Metadata.Location.Value != "ITBG-Bergamo" {
		t.Errorf("Metadata.Location.Value = %q", req.Metadata.Location.Value)
	}
	if req.Properties.Direction != types.RuleDirectionIngress {
		t.Errorf("Properties.Direction = %q", req.Properties.Direction)
	}
	if req.Properties.Protocol != "TCP" {
		t.Errorf("Properties.Protocol = %q", req.Properties.Protocol)
	}
	if req.Properties.Port != "22" {
		t.Errorf("Properties.Port = %q", req.Properties.Port)
	}
	if req.Properties.Target == nil || req.Properties.Target.Kind != types.EndpointTypeIP || req.Properties.Target.Value != "0.0.0.0/0" {
		t.Errorf("Properties.Target = %v", req.Properties.Target)
	}

	// Unset target → Properties.Target must be nil.
	rule2 := NewSecurityRule().WithName("no-target").WithDirection("Ingress")
	req2 := rule2.RawRequest()
	if req2.Properties.Target != nil {
		t.Errorf("Properties.Target should be nil when not set, got %v", req2.Properties.Target)
	}
}

// --------------------------------------------------------------------------
// fromResponse hydration
// --------------------------------------------------------------------------

func securityRuleTestResponse(id, name, uri, projectID string) *types.SecurityRuleResponse {
	state := "Active"
	dir := types.RuleDirectionEgress
	proto := "UDP"
	port := "53"
	return &types.SecurityRuleResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:   &id,
			URI:  &uri,
			Name: &name,
			Tags: []string{"rule-tag"},
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: projectID,
			},
			LocationResponse: &types.LocationResponse{
				Value: "ITBG-Bergamo",
			},
		},
		Properties: types.SecurityRulePropertiesResponse{
			Direction: dir,
			Protocol:  proto,
			Port:      port,
			Target:    &types.RuleTarget{Kind: types.EndpointTypeIP, Value: "1.2.3.4/32"},
		},
		Status: types.ResourceStatus{
			State: &state,
		},
	}
}

func TestSecurityRule_FromResponseHydration(t *testing.T) {
	resp := securityRuleTestResponse(
		"r-1",
		"allow-ssh",
		"/projects/p1/providers/Aruba.Network/vpcs/v1/securitygroups/sg-1/securityrules/r-1",
		"p1",
	)
	rule := &SecurityRule{}
	rule.fromResponse(resp)

	if rule.ID() != "r-1" {
		t.Errorf("ID() = %q", rule.ID())
	}
	if rule.Name() != "allow-ssh" {
		t.Errorf("Name() = %q", rule.Name())
	}
	if tags := rule.Tags(); len(tags) != 1 || tags[0] != "rule-tag" {
		t.Errorf("Tags() = %v", tags)
	}
	if rule.Region() != "ITBG-Bergamo" {
		t.Errorf("Region() = %q", rule.Region())
	}
	if rule.State() != "Active" {
		t.Errorf("State() = %q", rule.State())
	}
	if rule.Direction() != types.RuleDirectionEgress {
		t.Errorf("Direction() = %q", rule.Direction())
	}
	if rule.Protocol() != "UDP" {
		t.Errorf("Protocol() = %q", rule.Protocol())
	}
	if rule.Port() != "53" {
		t.Errorf("Port() = %q", rule.Port())
	}
	if rule.TargetKind() != types.EndpointTypeIP {
		t.Errorf("TargetKind() = %q", rule.TargetKind())
	}
	if rule.TargetValue() != "1.2.3.4/32" {
		t.Errorf("TargetValue() = %q", rule.TargetValue())
	}
	if rule.ProjectID() != "p1" {
		t.Errorf("ProjectID() = %q", rule.ProjectID())
	}
	if rule.securityGroupID != "sg-1" {
		t.Errorf("securityGroupID (from URI) = %q", rule.securityGroupID)
	}
	if rule.vpcID != "v1" {
		t.Errorf("vpcID (from URI) = %q", rule.vpcID)
	}
	if rule.Raw() != resp {
		t.Error("Raw() should return the hydrated response pointer")
	}
}

func TestSecurityRule_FromResponsePartial(t *testing.T) {
	rule := &SecurityRule{}
	rule.fromResponse(nil)
	if rule.ID() != "" || rule.URI() != "" || rule.Name() != "" {
		t.Error("fromResponse(nil) should be a no-op")
	}
	if rule.Raw() != nil {
		t.Error("Raw() should be nil before hydration")
	}

	rule2 := &SecurityRule{}
	rule2.fromResponse(&types.SecurityRuleResponse{})
	if rule2.ID() != "" || rule2.State() != "" {
		t.Error("empty response should yield zero accessor values")
	}
	if rule2.Direction() != "" {
		t.Errorf("Direction() from zero response = %q", rule2.Direction())
	}
	if rule2.TargetKind() != "" {
		t.Errorf("TargetKind() from zero response = %q", rule2.TargetKind())
	}
}

func TestSecurityRule_FromResponseURIBackfill_HyphenForm(t *testing.T) {
	uri := "/projects/p2/network/vpcs/v2/security-groups/sg-2/security-rules/r-2"
	id := "r-2"
	name := "rule-uri"
	resp := &types.SecurityRuleResponse{
		Metadata: types.ResourceMetadataResponse{
			ID:   &id,
			URI:  &uri,
			Name: &name,
		},
	}
	rule := &SecurityRule{}
	rule.fromResponse(resp)

	if rule.ProjectID() != "p2" {
		t.Errorf("ProjectID() via URI fallback = %q", rule.ProjectID())
	}
	if rule.vpcID != "v2" {
		t.Errorf("vpcID via URI fallback = %q", rule.vpcID)
	}
	if rule.securityGroupID != "sg-2" {
		t.Errorf("securityGroupID via URI fallback = %q", rule.securityGroupID)
	}
}

// --------------------------------------------------------------------------
// Ref + ancestor ID satisfaction (runtime)
// --------------------------------------------------------------------------

func TestSecurityRule_RefSatisfaction(t *testing.T) {
	rule := &SecurityRule{}
	rule.fromResponse(securityRuleTestResponse(
		"r-99",
		"n",
		"/projects/p99/network/vpcs/v99/security-groups/sg-99/security-rules/r-99",
		"p99",
	))

	sid, ok := extractID(rule, func(r Ref) (string, bool) {
		if w, ok := r.(withSecurityRuleID); ok {
			return w.SecurityRuleID(), true
		}
		return "", false
	}, "security-rules")
	if !ok || sid != "r-99" {
		t.Errorf("extractID via withSecurityRuleID = (%q, %v)", sid, ok)
	}

	sgid, ok := extractID(rule, func(r Ref) (string, bool) {
		if w, ok := r.(withSecurityGroupID); ok {
			return w.SecurityGroupID(), true
		}
		return "", false
	}, "security-groups")
	if !ok || sgid != "sg-99" {
		t.Errorf("extractID via withSecurityGroupID = (%q, %v)", sgid, ok)
	}

	vid, ok := extractID(rule, func(r Ref) (string, bool) {
		if w, ok := r.(withVPCID); ok {
			return w.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok || vid != "v99" {
		t.Errorf("extractID via withVPCID = (%q, %v)", vid, ok)
	}

	pid, ok := extractID(rule, func(r Ref) (string, bool) {
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
// securityRuleIDsFromRef helper
// --------------------------------------------------------------------------

func TestSecurityRuleIDsFromRef_TypedRef(t *testing.T) {
	rule := &SecurityRule{}
	rule.fromResponse(securityRuleTestResponse(
		"r-1",
		"n",
		"/projects/p/network/vpcs/v/security-groups/sg/security-rules/r-1",
		"p",
	))
	pid, vid, sgid, rid, err := securityRuleIDsFromRef(rule)
	if err != nil || pid != "p" || vid != "v" || sgid != "sg" || rid != "r-1" {
		t.Errorf("securityRuleIDsFromRef typed = (%q, %q, %q, %q, %v)", pid, vid, sgid, rid, err)
	}
}

func TestSecurityRuleIDsFromRef_URIRef_APIForm(t *testing.T) {
	ref := URI("/projects/p/providers/Aruba.Network/vpcs/v/securitygroups/sg/securityrules/r")
	pid, vid, sgid, rid, err := securityRuleIDsFromRef(ref)
	if err != nil || pid != "p" || vid != "v" || sgid != "sg" || rid != "r" {
		t.Errorf("securityRuleIDsFromRef API form = (%q, %q, %q, %q, %v)", pid, vid, sgid, rid, err)
	}
}

func TestSecurityRuleIDsFromRef_URIRef_HyphenForm(t *testing.T) {
	ref := URI("/projects/p/network/vpcs/v/security-groups/sg/security-rules/r")
	pid, vid, sgid, rid, err := securityRuleIDsFromRef(ref)
	if err != nil || pid != "p" || vid != "v" || sgid != "sg" || rid != "r" {
		t.Errorf("securityRuleIDsFromRef hyphen form = (%q, %q, %q, %q, %v)", pid, vid, sgid, rid, err)
	}
}

func TestSecurityRuleIDsFromRef_BadURI_MissingRule(t *testing.T) {
	_, _, _, _, err := securityRuleIDsFromRef(URI("/projects/p/network/vpcs/v/security-groups/sg"))
	if err == nil {
		t.Error("expected error for URI without rule segment")
	}
}

func TestSecurityRuleIDsFromRef_BadURI_MissingAll(t *testing.T) {
	_, _, _, _, err := securityRuleIDsFromRef(URI("/something/else"))
	if err == nil {
		t.Error("expected error for totally invalid URI")
	}
}

// --------------------------------------------------------------------------
// securityGroupRulesClientAdapter — CRUD integration tests
// --------------------------------------------------------------------------

func buildSecurityRuleTestAdapter(t *testing.T, handler http.HandlerFunc) *securityGroupRulesClientAdapter {
	t.Helper()
	server := testutil.NewMockServer(t, handler)
	return newSecurityGroupRulesClientAdapter(testutil.NewClient(t, server.URL))
}

// activeSecurityGroupBody is returned by the SG GET poll inside waitForSecurityGroupActive.
const activeSecurityGroupBody = `{"metadata":{"id":"sg","uri":"/projects/p/providers/Aruba.Network/vpcs/v/securitygroups/sg","project":{"id":"p"}},"properties":{"default":false},"status":{"state":"Active"}}`

// withSecurityGroupActiveRoute wraps a handler so SG GET requests (for the waiter poll) are
// answered with an "Active" state. Only rule-path requests are forwarded to ruleHandler.
func withSecurityGroupActiveRoute(ruleHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/securityrules") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, activeSecurityGroupBody)
			return
		}
		ruleHandler(w, r)
	}
}

const securityRuleSuccessBody = `{` +
	`"metadata":{"id":"r-1","name":"allow-ssh","uri":"/projects/p/network/vpcs/v/securitygroups/sg/securityrules/r-1","project":{"id":"p"}},` +
	`"properties":{"direction":"Ingress","protocol":"TCP","port":"22","target":{"kind":"Ip","value":"0.0.0.0/0"}},` +
	`"status":{"state":"Active"}}`

func TestSecurityGroupRulesClientAdapter_Create_Success(t *testing.T) {
	var gotBody types.SecurityRuleRequest
	adapter := buildSecurityRuleTestAdapter(t, withSecurityGroupActiveRoute(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, securityRuleSuccessBody)
	}))

	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg", "my-sg", "/projects/p/providers/Aruba.Network/vpcs/v/securitygroups/sg", "p"))

	rule := NewSecurityRule().
		IntoSecurityGroup(sg).
		WithName("allow-ssh").
		InRegion("ITBG-Bergamo").
		WithDirection(string(types.RuleDirectionIngress)).
		WithProtocol("TCP").
		WithPort("22").
		WithTargetCIDR("0.0.0.0/0")

	result, err := adapter.Create(context.Background(), rule)
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if result.ID() != "r-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.Name() != "allow-ssh" {
		t.Errorf("Name() = %q", result.Name())
	}
	if result.StatusCode() != http.StatusCreated {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if gotBody.Metadata.Name != "allow-ssh" {
		t.Errorf("request Name = %q", gotBody.Metadata.Name)
	}
	if gotBody.Properties.Protocol != "TCP" {
		t.Errorf("request Protocol = %q", gotBody.Properties.Protocol)
	}
}

func TestSecurityGroupRulesClientAdapter_Create_NoSG(t *testing.T) {
	callCount := 0
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
	})

	_, err := adapter.Create(context.Background(), NewSecurityRule().WithName("x"))
	if err == nil {
		t.Fatal("expected error when security rule has no SecurityGroup")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without SecurityGroup")
	}
}

func TestSecurityGroupRulesClientAdapter_Create_MetadataValidationError(t *testing.T) {
	adapter := buildSecurityRuleTestAdapter(t, withSecurityGroupActiveRoute(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		// Missing "id" field — triggers MetadataValidationError
		fmt.Fprint(w, `{"metadata":{"name":"rule","uri":"/projects/p/network/vpcs/v/securitygroups/sg/securityrules/x"},"properties":{},"status":{}}`)
	}))

	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg", "my-sg", "/projects/p/providers/Aruba.Network/vpcs/v/securitygroups/sg", "p"))

	rule := NewSecurityRule().IntoSecurityGroup(sg).WithName("rule")
	result, err := adapter.Create(context.Background(), rule)
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

func TestSecurityGroupRulesClientAdapter_Create_NonTwoXX(t *testing.T) {
	adapter := buildSecurityRuleTestAdapter(t, withSecurityGroupActiveRoute(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Validation Failed", "direction is required", 422))
	}))

	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg", "my-sg", "/projects/p/providers/Aruba.Network/vpcs/v/securitygroups/sg", "p"))

	rule := NewSecurityRule().IntoSecurityGroup(sg)
	result, err := adapter.Create(context.Background(), rule)
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

func TestSecurityGroupRulesClientAdapter_Get_URIRef(t *testing.T) {
	var capturedPath string
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, securityRuleSuccessBody)
	})

	ref := URI("/projects/p/network/vpcs/v/securitygroups/sg/security-rules/r-1")
	result, err := adapter.Get(context.Background(), ref)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "r-1" {
		t.Errorf("ID() = %q", result.ID())
	}
	if result.ProjectID() != "p" {
		t.Errorf("ProjectID() = %q", result.ProjectID())
	}
	if result.VPCID() != "v" {
		t.Errorf("VPCID() = %q", result.VPCID())
	}
	if result.SecurityGroupID() != "sg" {
		t.Errorf("SecurityGroupID() = %q", result.SecurityGroupID())
	}
	if result.StatusCode() != http.StatusOK {
		t.Errorf("StatusCode() = %d", result.StatusCode())
	}
	if !strings.Contains(capturedPath, "securityrules") {
		t.Errorf("path = %q, expected securityrules segment", capturedPath)
	}
}

func TestSecurityGroupRulesClientAdapter_Get_TypedRef(t *testing.T) {
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, securityRuleSuccessBody)
	})

	existing := &SecurityRule{}
	existing.fromResponse(securityRuleTestResponse(
		"r-1",
		"allow-ssh",
		"/projects/p/network/vpcs/v/security-groups/sg/security-rules/r-1",
		"p",
	))

	result, err := adapter.Get(context.Background(), existing)
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if result.ID() != "r-1" {
		t.Errorf("ID() = %q", result.ID())
	}
}

func TestSecurityGroupRulesClientAdapter_Update_Success(t *testing.T) {
	var capturedBody types.SecurityRuleRequest
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewDecoder(r.Body).Decode(&capturedBody)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"metadata":{"id":"r-1","name":"allow-https","uri":"/projects/p/network/vpcs/v/securitygroups/sg/securityrules/r-1","project":{"id":"p"}},"properties":{"direction":"Ingress","protocol":"TCP","port":"443"},"status":{}}`)
	})

	rule := &SecurityRule{}
	rule.fromResponse(securityRuleTestResponse(
		"r-1",
		"allow-ssh",
		"/projects/p/network/vpcs/v/security-groups/sg/security-rules/r-1",
		"p",
	))
	rule.WithName("allow-https").WithPort("443")

	result, err := adapter.Update(context.Background(), rule)
	if err != nil {
		t.Fatalf("Update error: %v", err)
	}
	if result.Name() != "allow-https" {
		t.Errorf("Name() = %q", result.Name())
	}
	if capturedBody.Metadata.Name != "allow-https" {
		t.Errorf("request Name = %q", capturedBody.Metadata.Name)
	}
	if capturedBody.Properties.Port != "443" {
		t.Errorf("request Port = %q", capturedBody.Properties.Port)
	}
}

func TestSecurityGroupRulesClientAdapter_Update_NoID(t *testing.T) {
	callCount := 0
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	sg := &SecurityGroup{}
	sg.fromResponse(securityGroupTestResponse("sg", "my-sg", "/projects/p/network/vpcs/v/securitygroups/sg", "p"))

	rule := NewSecurityRule().IntoSecurityGroup(sg).WithName("x")
	_, err := adapter.Update(context.Background(), rule)
	if err == nil {
		t.Fatal("expected error when security rule has no ID")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made when ID is missing")
	}
}

func TestSecurityGroupRulesClientAdapter_Update_NoSG(t *testing.T) {
	callCount := 0
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.WriteHeader(http.StatusOK)
	})

	rule := &SecurityRule{}
	id := "r-1"
	rule.fromResponse(&types.SecurityRuleResponse{
		Metadata: types.ResourceMetadataResponse{
			ID: &id,
		},
	})

	_, err := adapter.Update(context.Background(), rule)
	if err == nil {
		t.Fatal("expected error when security rule has no SecurityGroup")
	}
	if callCount != 0 {
		t.Error("no HTTP call should be made without SecurityGroup")
	}
}

func TestSecurityGroupRulesClientAdapter_Delete_Success(t *testing.T) {
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/network/vpcs/v/securitygroups/sg/securityrules/r-1"))
	if err != nil {
		t.Fatalf("Delete error: %v", err)
	}
}

func TestSecurityGroupRulesClientAdapter_Delete_NonTwoXX(t *testing.T) {
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "security rule not found", 404))
	})

	err := adapter.Delete(context.Background(), URI("/projects/p/network/vpcs/v/securitygroups/sg/securityrules/missing"))
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

func TestSecurityGroupRulesClientAdapter_List_TwoItems(t *testing.T) {
	adapter := buildSecurityRuleTestAdapter(t, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"total":2,"self":"","prev":"","next":"","first":"","last":"","values":[`+
			`{"metadata":{"id":"r-1","name":"rule-a","uri":"/projects/p/network/vpcs/v/securitygroups/sg/securityrules/r-1","project":{"id":"p"}},"properties":{"direction":"Ingress","protocol":"TCP","port":"22"},"status":{"state":"Active"}},`+
			`{"metadata":{"id":"r-2","name":"rule-b","uri":"/projects/p/network/vpcs/v/securitygroups/sg/securityrules/r-2","project":{"id":"p"}},"properties":{"direction":"Egress","protocol":"UDP","port":"53"},"status":{"state":"Active"}}`+
			`]}`)
	})

	sg := URI("/projects/p/network/vpcs/v/security-groups/sg")
	list, err := adapter.List(context.Background(), sg)
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
	if items[0].ID() != "r-1" || items[0].Name() != "rule-a" {
		t.Errorf("items[0] = {%q, %q}", items[0].ID(), items[0].Name())
	}
	if items[0].Direction() != types.RuleDirectionIngress {
		t.Errorf("items[0].Direction() = %q", items[0].Direction())
	}
	if items[1].ID() != "r-2" || items[1].Direction() != types.RuleDirectionEgress {
		t.Errorf("items[1] = {%q, %q}", items[1].ID(), items[1].Direction())
	}
	if items[0].SecurityGroupID() != "sg" {
		t.Errorf("items[0].SecurityGroupID() = %q", items[0].SecurityGroupID())
	}
	if items[0].VPCID() != "v" {
		t.Errorf("items[0].VPCID() = %q", items[0].VPCID())
	}
	if items[0].ProjectID() != "p" {
		t.Errorf("items[0].ProjectID() = %q", items[0].ProjectID())
	}
}
