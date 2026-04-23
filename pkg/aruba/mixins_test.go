package aruba

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// --------------------------------------------------------------------------
// errMixin
// --------------------------------------------------------------------------

func TestErrMixin_Empty(t *testing.T) {
	var m errMixin
	if m.Err() != nil {
		t.Errorf("expected nil error, got %v", m.Err())
	}
}

func TestErrMixin_SingleError(t *testing.T) {
	var m errMixin
	sentinel := errors.New("oops")
	m.addErr(sentinel)
	if !errors.Is(m.Err(), sentinel) {
		t.Errorf("Err() does not wrap sentinel: %v", m.Err())
	}
}

func TestErrMixin_MultipleErrors(t *testing.T) {
	var m errMixin
	e1 := errors.New("first")
	e2 := errors.New("second")
	m.addErr(e1)
	m.addErr(e2)
	err := m.Err()
	if !errors.Is(err, e1) {
		t.Errorf("Err() does not wrap e1: %v", err)
	}
	if !errors.Is(err, e2) {
		t.Errorf("Err() does not wrap e2: %v", err)
	}
}

func TestErrMixin_AddNilIsNoop(t *testing.T) {
	var m errMixin
	m.addErr(nil)
	if m.Err() != nil {
		t.Errorf("expected nil after adding nil error, got %v", m.Err())
	}
}

// --------------------------------------------------------------------------
// metadataMixin
// --------------------------------------------------------------------------

func TestMetadataMixin(t *testing.T) {
	var m metadataMixin
	m.withName("hello")
	if m.Name() != "hello" {
		t.Errorf("Name() = %q", m.Name())
	}

	m.addTag("a")
	m.addTag("b")
	m.addTag("a") // duplicate — should not appear twice
	if got := m.Tags(); len(got) != 2 {
		t.Errorf("Tags() length = %d, want 2; tags=%v", len(got), got)
	}

	m.removeTag("a")
	tags := m.Tags()
	if len(tags) != 1 || tags[0] != "b" {
		t.Errorf("after RemoveTag(a): %v", tags)
	}

	m.removeTag("nonexistent") // no-op
	if len(m.Tags()) != 1 {
		t.Errorf("RemoveTag of missing tag changed slice")
	}

	m.replaceTags("x", "y", "z")
	if got := m.Tags(); len(got) != 3 {
		t.Errorf("after ReplaceTags: %v", got)
	}

	req := m.toMetadata()
	if req.Name != "hello" {
		t.Errorf("toMetadata().Name = %q", req.Name)
	}
	if len(req.Tags) != 3 {
		t.Errorf("toMetadata().Tags = %v", req.Tags)
	}
}

// --------------------------------------------------------------------------
// regionalMixin
// --------------------------------------------------------------------------

func TestRegionalMixin(t *testing.T) {
	var m regionalMixin
	m.withLocation("eu-west")
	if m.Region() != "eu-west" {
		t.Errorf("Region() = %q", m.Region())
	}
	m.inRegion("us-east")
	if m.Region() != "us-east" {
		t.Errorf("InRegion alias: Region() = %q", m.Region())
	}
	loc := m.toLocation()
	if loc.Value != "us-east" {
		t.Errorf("toLocation().Value = %q", loc.Value)
	}
}

// --------------------------------------------------------------------------
// projectScopedMixin
// --------------------------------------------------------------------------

// typedProjectParent implements withProjectID + Ref.
type typedProjectParent struct{ id string }

func (p typedProjectParent) ProjectID() string { return p.id }
func (p typedProjectParent) URI() string       { return "/projects/" + p.id }
func (p typedProjectParent) ID() string        { return p.id }

func TestProjectScopedMixin_TypedRef(t *testing.T) {
	errs := &errMixin{}
	m := bindProjectScoped(errs)
	m.intoProject(typedProjectParent{id: "proj-1"})

	if m.ProjectID() != "proj-1" {
		t.Errorf("ProjectID() = %q", m.ProjectID())
	}
	if errs.Err() != nil {
		t.Errorf("unexpected error: %v", errs.Err())
	}
}

func TestProjectScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindProjectScoped(errs)
	m.intoProject(URI("/projects/proj-2"))

	if m.ProjectID() != "proj-2" {
		t.Errorf("ProjectID() via URI = %q", m.ProjectID())
	}
	if errs.Err() != nil {
		t.Errorf("unexpected error: %v", errs.Err())
	}
}

func TestProjectScopedMixin_MissingSegment(t *testing.T) {
	errs := &errMixin{}
	m := bindProjectScoped(errs)
	m.intoProject(URI("/network/vpcs/v")) // no "projects" segment

	if errs.Err() == nil {
		t.Error("expected error for missing project segment, got nil")
	}
}

// --------------------------------------------------------------------------
// vpcScopedMixin
// --------------------------------------------------------------------------

type typedVPCParent struct{ vpcID, projectID string }

func (p typedVPCParent) VPCID() string     { return p.vpcID }
func (p typedVPCParent) ProjectID() string { return p.projectID }
func (p typedVPCParent) URI() string {
	return "/projects/" + p.projectID + "/network/vpcs/" + p.vpcID
}
func (p typedVPCParent) ID() string { return p.vpcID }

func TestVPCScopedMixin_TypedRef(t *testing.T) {
	errs := &errMixin{}
	m := bindVPCScoped(errs)
	m.intoVPC(typedVPCParent{vpcID: "vpc-1", projectID: "proj-1"})

	if m.VPCID() != "vpc-1" {
		t.Errorf("VPCID() = %q", m.VPCID())
	}
	if m.ProjectID() != "proj-1" {
		t.Errorf("ProjectID() = %q", m.ProjectID())
	}
	if errs.Err() != nil {
		t.Errorf("unexpected error: %v", errs.Err())
	}
}

func TestVPCScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindVPCScoped(errs)
	m.intoVPC(URI("/projects/proj-2/network/vpcs/vpc-2"))

	if m.VPCID() != "vpc-2" {
		t.Errorf("VPCID() via URI = %q", m.VPCID())
	}
	if m.ProjectID() != "proj-2" {
		t.Errorf("ProjectID() via URI = %q", m.ProjectID())
	}
}

func TestVPCScopedMixin_MissingVPCSegment(t *testing.T) {
	errs := &errMixin{}
	m := bindVPCScoped(errs)
	m.intoVPC(URI("/projects/proj-1/network")) // no vpcs segment

	if errs.Err() == nil {
		t.Error("expected error for missing VPC segment")
	}
}

// --------------------------------------------------------------------------
// securityGroupScopedMixin
// --------------------------------------------------------------------------

type typedSGParent struct{ sgID, vpcID, projectID string }

func (p typedSGParent) SecurityGroupID() string { return p.sgID }
func (p typedSGParent) VPCID() string           { return p.vpcID }
func (p typedSGParent) ProjectID() string       { return p.projectID }
func (p typedSGParent) URI() string {
	return "/projects/" + p.projectID + "/network/vpcs/" + p.vpcID + "/security-groups/" + p.sgID
}
func (p typedSGParent) ID() string { return p.sgID }

func TestSecurityGroupScopedMixin_TypedRef(t *testing.T) {
	errs := &errMixin{}
	m := bindSecurityGroupScoped(errs)
	m.intoSecurityGroup(typedSGParent{sgID: "sg-1", vpcID: "vpc-1", projectID: "proj-1"})

	if m.SecurityGroupID() != "sg-1" {
		t.Errorf("SecurityGroupID() = %q", m.SecurityGroupID())
	}
	if m.VPCID() != "vpc-1" {
		t.Errorf("VPCID() = %q", m.VPCID())
	}
	if m.ProjectID() != "proj-1" {
		t.Errorf("ProjectID() = %q", m.ProjectID())
	}
	if errs.Err() != nil {
		t.Errorf("unexpected error: %v", errs.Err())
	}
}

func TestSecurityGroupScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindSecurityGroupScoped(errs)
	m.intoSecurityGroup(URI("/projects/proj-2/network/vpcs/vpc-2/security-groups/sg-2"))

	if m.SecurityGroupID() != "sg-2" || m.VPCID() != "vpc-2" || m.ProjectID() != "proj-2" {
		t.Errorf("got sg=%q vpc=%q proj=%q", m.SecurityGroupID(), m.VPCID(), m.ProjectID())
	}
}

// --------------------------------------------------------------------------
// dbaasScopedMixin
// --------------------------------------------------------------------------

func TestDBaaSScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindDBaaSScoped(errs)
	m.intoDBaaS(URI("/projects/proj-1/database/dbaas/db-1"))

	if m.DBaaSID() != "db-1" || m.ProjectID() != "proj-1" {
		t.Errorf("got dbaas=%q proj=%q", m.DBaaSID(), m.ProjectID())
	}
	if errs.Err() != nil {
		t.Errorf("unexpected error: %v", errs.Err())
	}
}

// --------------------------------------------------------------------------
// databaseScopedMixin
// --------------------------------------------------------------------------

func TestDatabaseScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindDatabaseScoped(errs)
	m.intoDatabase(URI("/projects/proj-1/database/dbaas/db-1/databases/mydb"))

	if m.DatabaseID() != "mydb" || m.DBaaSID() != "db-1" || m.ProjectID() != "proj-1" {
		t.Errorf("got db=%q dbaas=%q proj=%q", m.DatabaseID(), m.DBaaSID(), m.ProjectID())
	}
}

// --------------------------------------------------------------------------
// backupScopedMixin
// --------------------------------------------------------------------------

func TestBackupScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindBackupScoped(errs)
	m.intoBackup(URI("/projects/proj-1/storage/backups/bk-1"))

	if m.BackupID() != "bk-1" || m.ProjectID() != "proj-1" {
		t.Errorf("got backup=%q proj=%q", m.BackupID(), m.ProjectID())
	}
}

// --------------------------------------------------------------------------
// kmsScopedMixin
// --------------------------------------------------------------------------

func TestKMSScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindKMSScoped(errs)
	m.intoKMS(URI("/projects/proj-1/security/kms/kms-1"))

	if m.KMSID() != "kms-1" || m.ProjectID() != "proj-1" {
		t.Errorf("got kms=%q proj=%q", m.KMSID(), m.ProjectID())
	}
}

// --------------------------------------------------------------------------
// vpnTunnelScopedMixin
// --------------------------------------------------------------------------

func TestVPNTunnelScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindVPNTunnelScoped(errs)
	m.intoVPNTunnel(URI("/projects/proj-1/network/vpn-tunnels/t-1"))

	if m.VPNTunnelID() != "t-1" || m.ProjectID() != "proj-1" {
		t.Errorf("got tunnel=%q proj=%q", m.VPNTunnelID(), m.ProjectID())
	}
}

// --------------------------------------------------------------------------
// vpcPeeringScopedMixin
// --------------------------------------------------------------------------

func TestVPCPeeringScopedMixin_URIFallback(t *testing.T) {
	errs := &errMixin{}
	m := bindVPCPeeringScoped(errs)
	m.intoVPCPeering(URI("/projects/proj-1/network/vpcs/vpc-1/peerings/peer-1"))

	if m.VPCPeeringID() != "peer-1" || m.VPCID() != "vpc-1" || m.ProjectID() != "proj-1" {
		t.Errorf("got peering=%q vpc=%q proj=%q", m.VPCPeeringID(), m.VPCID(), m.ProjectID())
	}
}

// --------------------------------------------------------------------------
// responseMetadataMixin
// --------------------------------------------------------------------------

func TestResponseMetadataMixin_Nil(t *testing.T) {
	var m responseMetadataMixin
	if m.ID() != "" {
		t.Errorf("ID() on nil meta = %q", m.ID())
	}
	if m.RespURI() != "" {
		t.Errorf("RespURI() on nil meta = %q", m.RespURI())
	}
	if m.Project() != "" {
		t.Errorf("Project() on nil meta = %q", m.Project())
	}
	if !m.CreatedAt().IsZero() {
		t.Errorf("CreatedAt() should be zero, got %v", m.CreatedAt())
	}
	if !m.UpdatedAt().IsZero() {
		t.Errorf("UpdatedAt() should be zero, got %v", m.UpdatedAt())
	}
	if m.Version() != "" {
		t.Errorf("Version() on nil meta = %q", m.Version())
	}
}

func TestResponseMetadataMixin_Populated(t *testing.T) {
	id := "res-123"
	uri := "/projects/p/vpcs/v"
	ver := "1"
	proj := "proj-1"
	now := time.Now().UTC().Truncate(time.Second)
	later := now.Add(time.Hour)

	m := responseMetadataMixin{
		meta: &types.ResourceMetadataResponse{
			ID:      &id,
			URI:     &uri,
			Version: &ver,
			ProjectResponseMetadata: &types.ProjectResponseMetadata{
				ID: proj,
			},
			CreationDate: &now,
			UpdateDate:   &later,
		},
	}

	if m.ID() != id {
		t.Errorf("ID() = %q, want %q", m.ID(), id)
	}
	if m.RespURI() != uri {
		t.Errorf("RespURI() = %q, want %q", m.RespURI(), uri)
	}
	if m.Project() != proj {
		t.Errorf("Project() = %q, want %q", m.Project(), proj)
	}
	if m.Version() != ver {
		t.Errorf("Version() = %q, want %q", m.Version(), ver)
	}
	if !m.CreatedAt().Equal(now) {
		t.Errorf("CreatedAt() = %v, want %v", m.CreatedAt(), now)
	}
	if !m.UpdatedAt().Equal(later) {
		t.Errorf("UpdatedAt() = %v, want %v", m.UpdatedAt(), later)
	}
}

// --------------------------------------------------------------------------
// statusMixin
// --------------------------------------------------------------------------

func TestStatusMixin_Nil(t *testing.T) {
	var m statusMixin
	if m.State() != "" {
		t.Errorf("State() on nil = %q", m.State())
	}
	if m.IsDisabled() {
		t.Error("IsDisabled() on nil should be false")
	}
	if m.DisableReasons() != nil {
		t.Errorf("DisableReasons() on nil = %v", m.DisableReasons())
	}
	if m.FailureReason() != "" {
		t.Errorf("FailureReason() on nil = %q", m.FailureReason())
	}
	if m.PreviousState() != "" {
		t.Errorf("PreviousState() on nil = %q", m.PreviousState())
	}
}

func TestStatusMixin_Populated(t *testing.T) {
	state := "Active"
	prev := "Pending"
	reason := "disk full"
	m := statusMixin{
		status: &types.ResourceStatus{
			State:         &state,
			FailureReason: &reason,
			DisableStatusInfo: &types.DisableStatusInfo{
				IsDisabled: true,
				Reasons:    []string{"maintenance"},
			},
			PreviousStatus: &types.PreviousStatus{
				State: &prev,
			},
		},
	}

	if m.State() != "Active" {
		t.Errorf("State() = %q", m.State())
	}
	if !m.IsDisabled() {
		t.Error("IsDisabled() should be true")
	}
	if len(m.DisableReasons()) != 1 || m.DisableReasons()[0] != "maintenance" {
		t.Errorf("DisableReasons() = %v", m.DisableReasons())
	}
	if m.FailureReason() != "disk full" {
		t.Errorf("FailureReason() = %q", m.FailureReason())
	}
	if m.PreviousState() != "Pending" {
		t.Errorf("PreviousState() = %q", m.PreviousState())
	}
}

func TestStatusMixin_WaitStubs(t *testing.T) {
	var m statusMixin
	if err := m.WaitUntilActive(context.Background()); err == nil {
		t.Error("WaitUntilActive should return not-implemented error")
	}
	if err := m.WaitUntilState(context.Background(), "Active"); err == nil {
		t.Error("WaitUntilState should return not-implemented error")
	}
}

// --------------------------------------------------------------------------
// httpEnvelopeMixin
// --------------------------------------------------------------------------

func TestHTTPEnvelopeMixin(t *testing.T) {
	var m httpEnvelopeMixin

	title := "Not Found"
	status := int32(404)
	resp := &types.Response[struct{}]{
		StatusCode:   404,
		Headers:      http.Header{"X-Trace": []string{"abc"}},
		RawBody:      []byte(`{"title":"Not Found"}`),
		HTTPResponse: &http.Response{StatusCode: 404},
		Error: &types.ErrorResponse{
			Title:  &title,
			Status: &status,
		},
	}
	populateHTTPEnvelope(&m, resp)

	if m.StatusCode() != 404 {
		t.Errorf("StatusCode() = %d", m.StatusCode())
	}
	if m.Headers().Get("X-Trace") != "abc" {
		t.Errorf("Headers() = %v", m.Headers())
	}
	httpResp, raw := m.RawHTTP()
	if httpResp == nil || httpResp.StatusCode != 404 {
		t.Errorf("RawHTTP() http.Response = %v", httpResp)
	}
	if string(raw) != `{"title":"Not Found"}` {
		t.Errorf("RawHTTP() raw = %q", raw)
	}
	if m.RawError() == nil || *m.RawError().Title != "Not Found" {
		t.Errorf("RawError() = %v", m.RawError())
	}
}

func TestHTTPEnvelopeMixin_NilResponse(t *testing.T) {
	var m httpEnvelopeMixin
	populateHTTPEnvelope[struct{}](&m, nil)
	if m.StatusCode() != 0 {
		t.Errorf("StatusCode() after nil response = %d", m.StatusCode())
	}
}
