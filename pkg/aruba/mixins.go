package aruba

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// --------------------------------------------------------------------------
// errMixin — setter-time error accumulator
// --------------------------------------------------------------------------

type errMixin struct {
	errs []error
}

func (m *errMixin) addErr(err error) {
	if err != nil {
		m.errs = append(m.errs, err)
	}
}

// Err returns the joined setter-time errors, or nil if none were recorded.
func (m *errMixin) Err() error {
	return errors.Join(m.errs...)
}

// validateAll is an alias for Err; used internally before issuing an HTTP call.
func (m *errMixin) validateAll() error {
	return m.Err()
}

// --------------------------------------------------------------------------
// metadataMixin — resource name and tags
// --------------------------------------------------------------------------

type metadataMixin struct {
	name string
	tags []string
}

func (m *metadataMixin) withName(name string) {
	m.name = name
}

func (m *metadataMixin) addTag(tag string) {
	for _, t := range m.tags {
		if t == tag {
			return
		}
	}
	m.tags = append(m.tags, tag)
}

func (m *metadataMixin) removeTag(tag string) {
	out := m.tags[:0]
	for _, t := range m.tags {
		if t != tag {
			out = append(out, t)
		}
	}
	m.tags = out
}

func (m *metadataMixin) replaceTags(tags ...string) {
	m.tags = append([]string(nil), tags...)
}

// Name returns the name set via withName.
func (m *metadataMixin) Name() string { return m.name }

// Tags returns a copy of the current tag slice.
func (m *metadataMixin) Tags() []string {
	if len(m.tags) == 0 {
		return nil
	}
	out := make([]string, len(m.tags))
	copy(out, m.tags)
	return out
}

func (m *metadataMixin) toMetadata() types.ResourceMetadataRequest {
	return types.ResourceMetadataRequest{Name: m.name, Tags: m.Tags()}
}

// --------------------------------------------------------------------------
// regionalMixin — resource location / region
// --------------------------------------------------------------------------

type regionalMixin struct {
	location string
}

func (m *regionalMixin) withLocation(loc string) { m.location = loc }
func (m *regionalMixin) inRegion(region string)   { m.location = region }

// Region returns the location value.
func (m *regionalMixin) Region() string { return m.location }

func (m *regionalMixin) toLocation() types.LocationRequest {
	return types.LocationRequest{Value: m.location}
}

// --------------------------------------------------------------------------
// Scoped mixins — parent hierarchy
// --------------------------------------------------------------------------

// projectScopedMixin — direct child of a Project.
type projectScopedMixin struct {
	projectID string
	errSink   *errMixin
}

func bindProjectScoped(errSink *errMixin) projectScopedMixin {
	return projectScopedMixin{errSink: errSink}
}

func (m *projectScopedMixin) intoProject(parent Ref) {
	id, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoProject: cannot determine project ID from Ref %q", parent.URI()))
		return
	}
	m.projectID = id
}

// ProjectID returns the parent project's ID.
func (m *projectScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// vpcScopedMixin — direct child of a VPC; inherits projectID from its VPC parent.
type vpcScopedMixin struct {
	vpcID     string
	projectID string
	errSink   *errMixin
}

func bindVPCScoped(errSink *errMixin) vpcScopedMixin {
	return vpcScopedMixin{errSink: errSink}
}

func (m *vpcScopedMixin) intoVPC(parent Ref) {
	vpcID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withVPCID); ok {
			return p.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoVPC: cannot determine VPC ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoVPC: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.vpcID = vpcID
	m.projectID = projectID
}

// VPCID returns the parent VPC's ID.
func (m *vpcScopedMixin) VPCID() string { return m.vpcID }

// ProjectID returns the inherited project ID.
func (m *vpcScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// securityGroupScopedMixin — direct child of a SecurityGroup; inherits vpcID and projectID.
type securityGroupScopedMixin struct {
	securityGroupID string
	vpcID           string
	projectID       string
	errSink         *errMixin
}

func bindSecurityGroupScoped(errSink *errMixin) securityGroupScopedMixin {
	return securityGroupScopedMixin{errSink: errSink}
}

func (m *securityGroupScopedMixin) intoSecurityGroup(parent Ref) {
	sgID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withSecurityGroupID); ok {
			return p.SecurityGroupID(), true
		}
		return "", false
	}, "security-groups")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoSecurityGroup: cannot determine security group ID from Ref %q", parent.URI()))
	}

	vpcID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withVPCID); ok {
			return p.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoSecurityGroup: cannot determine VPC ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoSecurityGroup: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.securityGroupID = sgID
	m.vpcID = vpcID
	m.projectID = projectID
}

// SecurityGroupID returns the parent security group's ID.
func (m *securityGroupScopedMixin) SecurityGroupID() string { return m.securityGroupID }

// VPCID returns the inherited VPC ID.
func (m *securityGroupScopedMixin) VPCID() string { return m.vpcID }

// ProjectID returns the inherited project ID.
func (m *securityGroupScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// dbaasScopedMixin — direct child of a DBaaS instance; inherits projectID.
type dbaasScopedMixin struct {
	dbaasID   string
	projectID string
	errSink   *errMixin
}

func bindDBaaSScoped(errSink *errMixin) dbaasScopedMixin {
	return dbaasScopedMixin{errSink: errSink}
}

func (m *dbaasScopedMixin) intoDBaaS(parent Ref) {
	dbaasID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withDBaaSID); ok {
			return p.DBaaSID(), true
		}
		return "", false
	}, "dbaas")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoDBaaS: cannot determine DBaaS ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoDBaaS: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.dbaasID = dbaasID
	m.projectID = projectID
}

// DBaaSID returns the parent DBaaS instance's ID.
func (m *dbaasScopedMixin) DBaaSID() string { return m.dbaasID }

// ProjectID returns the inherited project ID.
func (m *dbaasScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// databaseScopedMixin — direct child of a Database; inherits dbaasID and projectID.
type databaseScopedMixin struct {
	databaseID string
	dbaasID    string
	projectID  string
	errSink    *errMixin
}

func bindDatabaseScoped(errSink *errMixin) databaseScopedMixin {
	return databaseScopedMixin{errSink: errSink}
}

func (m *databaseScopedMixin) intoDatabase(parent Ref) {
	dbID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withDatabaseID); ok {
			return p.DatabaseID(), true
		}
		return "", false
	}, "databases")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoDatabase: cannot determine database ID from Ref %q", parent.URI()))
	}

	dbaasID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withDBaaSID); ok {
			return p.DBaaSID(), true
		}
		return "", false
	}, "dbaas")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoDatabase: cannot determine DBaaS ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoDatabase: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.databaseID = dbID
	m.dbaasID = dbaasID
	m.projectID = projectID
}

// DatabaseID returns the parent database's ID.
func (m *databaseScopedMixin) DatabaseID() string { return m.databaseID }

// DBaaSID returns the inherited DBaaS ID.
func (m *databaseScopedMixin) DBaaSID() string { return m.dbaasID }

// ProjectID returns the inherited project ID.
func (m *databaseScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// backupScopedMixin — direct child of a StorageBackup; inherits projectID.
type backupScopedMixin struct {
	backupID  string
	projectID string
	errSink   *errMixin
}

func bindBackupScoped(errSink *errMixin) backupScopedMixin {
	return backupScopedMixin{errSink: errSink}
}

func (m *backupScopedMixin) intoBackup(parent Ref) {
	backupID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withBackupID); ok {
			return p.BackupID(), true
		}
		return "", false
	}, "backups")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoBackup: cannot determine backup ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoBackup: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.backupID = backupID
	m.projectID = projectID
}

// BackupID returns the parent backup's ID.
func (m *backupScopedMixin) BackupID() string { return m.backupID }

// ProjectID returns the inherited project ID.
func (m *backupScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// kmsScopedMixin — direct child of a KMS instance; inherits projectID.
type kmsScopedMixin struct {
	kmsID     string
	projectID string
	errSink   *errMixin
}

func bindKMSScoped(errSink *errMixin) kmsScopedMixin {
	return kmsScopedMixin{errSink: errSink}
}

func (m *kmsScopedMixin) intoKMS(parent Ref) {
	kmsID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withKMSID); ok {
			return p.KMSID(), true
		}
		return "", false
	}, "kms")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoKMS: cannot determine KMS ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoKMS: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.kmsID = kmsID
	m.projectID = projectID
}

// KMSID returns the parent KMS instance's ID.
func (m *kmsScopedMixin) KMSID() string { return m.kmsID }

// ProjectID returns the inherited project ID.
func (m *kmsScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// vpnTunnelScopedMixin — direct child of a VPN tunnel; inherits projectID.
type vpnTunnelScopedMixin struct {
	vpnTunnelID string
	projectID   string
	errSink     *errMixin
}

func bindVPNTunnelScoped(errSink *errMixin) vpnTunnelScopedMixin {
	return vpnTunnelScopedMixin{errSink: errSink}
}

func (m *vpnTunnelScopedMixin) intoVPNTunnel(parent Ref) {
	tunnelID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withVPNTunnelID); ok {
			return p.VPNTunnelID(), true
		}
		return "", false
	}, "vpn-tunnels")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoVPNTunnel: cannot determine VPN tunnel ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoVPNTunnel: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.vpnTunnelID = tunnelID
	m.projectID = projectID
}

// VPNTunnelID returns the parent VPN tunnel's ID.
func (m *vpnTunnelScopedMixin) VPNTunnelID() string { return m.vpnTunnelID }

// ProjectID returns the inherited project ID.
func (m *vpnTunnelScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------

// vpcPeeringScopedMixin — direct child of a VPC peering; inherits vpcID and projectID.
type vpcPeeringScopedMixin struct {
	vpcPeeringID string
	vpcID        string
	projectID    string
	errSink      *errMixin
}

func bindVPCPeeringScoped(errSink *errMixin) vpcPeeringScopedMixin {
	return vpcPeeringScopedMixin{errSink: errSink}
}

func (m *vpcPeeringScopedMixin) intoVPCPeering(parent Ref) {
	peeringID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withVPCPeeringID); ok {
			return p.VPCPeeringID(), true
		}
		return "", false
	}, "peerings")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoVPCPeering: cannot determine VPC peering ID from Ref %q", parent.URI()))
	}

	vpcID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withVPCID); ok {
			return p.VPCID(), true
		}
		return "", false
	}, "vpcs")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoVPCPeering: cannot determine VPC ID from Ref %q", parent.URI()))
	}

	projectID, ok := extractID(parent, func(r Ref) (string, bool) {
		if p, ok := r.(withProjectID); ok {
			return p.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok {
		m.errSink.addErr(fmt.Errorf("IntoVPCPeering: cannot determine project ID from Ref %q", parent.URI()))
	}

	m.vpcPeeringID = peeringID
	m.vpcID = vpcID
	m.projectID = projectID
}

// VPCPeeringID returns the parent VPC peering's ID.
func (m *vpcPeeringScopedMixin) VPCPeeringID() string { return m.vpcPeeringID }

// VPCID returns the inherited VPC ID.
func (m *vpcPeeringScopedMixin) VPCID() string { return m.vpcID }

// ProjectID returns the inherited project ID.
func (m *vpcPeeringScopedMixin) ProjectID() string { return m.projectID }

// --------------------------------------------------------------------------
// responseMetadataMixin — post-server-reply metadata
// --------------------------------------------------------------------------

type responseMetadataMixin struct {
	meta *types.ResourceMetadataResponse
}

func (m *responseMetadataMixin) setMeta(meta *types.ResourceMetadataResponse) {
	m.meta = meta
}

// ID returns the resource's server-assigned ID, or "" if not yet received.
func (m *responseMetadataMixin) ID() string {
	if m.meta == nil || m.meta.ID == nil {
		return ""
	}
	return *m.meta.ID
}

// RespURI returns the resource's server-assigned URI, or "" if not yet received.
// Named RespURI to avoid collision with the Ref.URI() method on wrapper types that
// derive their URI from the response.
func (m *responseMetadataMixin) RespURI() string {
	if m.meta == nil || m.meta.URI == nil {
		return ""
	}
	return *m.meta.URI
}

// Project returns the owning project's ID from the response metadata, or "".
func (m *responseMetadataMixin) Project() string {
	if m.meta == nil || m.meta.ProjectResponseMetadata == nil {
		return ""
	}
	return m.meta.ProjectResponseMetadata.ID
}

// CreatedAt returns the resource creation time, or zero time.
func (m *responseMetadataMixin) CreatedAt() time.Time {
	if m.meta == nil || m.meta.CreationDate == nil {
		return time.Time{}
	}
	return *m.meta.CreationDate
}

// UpdatedAt returns the last update time, or zero time.
func (m *responseMetadataMixin) UpdatedAt() time.Time {
	if m.meta == nil || m.meta.UpdateDate == nil {
		return time.Time{}
	}
	return *m.meta.UpdateDate
}

// Version returns the resource version string, or "".
func (m *responseMetadataMixin) Version() string {
	if m.meta == nil || m.meta.Version == nil {
		return ""
	}
	return *m.meta.Version
}

// Raw returns the underlying *types.ResourceMetadataResponse, or nil.
func (m *responseMetadataMixin) Raw() *types.ResourceMetadataResponse {
	return m.meta
}

// --------------------------------------------------------------------------
// statusMixin — resource lifecycle state
// --------------------------------------------------------------------------

// WaitOption configures WaitUntilActive / WaitUntilState behaviour.
// Full implementation is deferred to the async-poll issue; the type exists now
// so call sites compile without changes once the real implementation lands.
type WaitOption func(*waitOptions)

type waitOptions struct{}

type statusMixin struct {
	status  *types.ResourceStatus
	refresh func(ctx context.Context) error
}

func (m *statusMixin) setStatus(s *types.ResourceStatus) { m.status = s }

// State returns the current lifecycle state string, or "".
func (m *statusMixin) State() string {
	if m.status == nil || m.status.State == nil {
		return ""
	}
	return *m.status.State
}

// IsDisabled returns true when the server has disabled this resource.
func (m *statusMixin) IsDisabled() bool {
	if m.status == nil || m.status.DisableStatusInfo == nil {
		return false
	}
	return m.status.DisableStatusInfo.IsDisabled
}

// DisableReasons returns the reasons for disabling, or nil.
func (m *statusMixin) DisableReasons() []string {
	if m.status == nil || m.status.DisableStatusInfo == nil {
		return nil
	}
	return m.status.DisableStatusInfo.Reasons
}

// FailureReason returns the failure reason string, or "".
func (m *statusMixin) FailureReason() string {
	if m.status == nil || m.status.FailureReason == nil {
		return ""
	}
	return *m.status.FailureReason
}

// PreviousState returns the previous lifecycle state string, or "".
func (m *statusMixin) PreviousState() string {
	if m.status == nil || m.status.PreviousStatus == nil || m.status.PreviousStatus.State == nil {
		return ""
	}
	return *m.status.PreviousStatus.State
}

// WaitUntilActive blocks until the resource reaches the Active state.
// Full implementation is deferred to the async-poll issue.
func (m *statusMixin) WaitUntilActive(_ context.Context, _ ...WaitOption) error {
	return errors.New("WaitUntilActive: not yet implemented (see async-poll issue)")
}

// WaitUntilState blocks until the resource reaches the given state.
// Full implementation is deferred to the async-poll issue.
func (m *statusMixin) WaitUntilState(_ context.Context, _ string, _ ...WaitOption) error {
	return errors.New("WaitUntilState: not yet implemented (see async-poll issue)")
}

// --------------------------------------------------------------------------
// linkedMixin — linked resources
// --------------------------------------------------------------------------

type linkedMixin struct {
	linked []types.LinkedResource
}

func (m *linkedMixin) setLinked(l []types.LinkedResource) { m.linked = l }

// LinkedResources returns the slice of linked resources.
func (m *linkedMixin) LinkedResources() []types.LinkedResource { return m.linked }

// --------------------------------------------------------------------------
// httpEnvelopeMixin — HTTP response metadata
// --------------------------------------------------------------------------

type httpEnvelopeMixin struct {
	statusCode int
	headers    http.Header
	rawBody    []byte
	httpResp   *http.Response
	errResp    *types.ErrorResponse
}

// populateHTTPEnvelope fills an httpEnvelopeMixin from a typed *types.Response[T].
// Defined as a package-level generic function because Go does not allow generic methods.
func populateHTTPEnvelope[T any](m *httpEnvelopeMixin, resp *types.Response[T]) {
	if resp == nil {
		return
	}
	m.statusCode = resp.StatusCode
	m.headers = resp.Headers
	m.rawBody = resp.RawBody
	m.httpResp = resp.HTTPResponse
	m.errResp = resp.Error
}

// StatusCode returns the HTTP status code, or 0 before any response.
func (m *httpEnvelopeMixin) StatusCode() int { return m.statusCode }

// Headers returns the HTTP response headers, or nil.
func (m *httpEnvelopeMixin) Headers() http.Header { return m.headers }

// RawHTTP returns the underlying *http.Response and raw body bytes.
func (m *httpEnvelopeMixin) RawHTTP() (*http.Response, []byte) {
	return m.httpResp, m.rawBody
}

// RawError returns the parsed error response body for non-2xx replies, or nil.
func (m *httpEnvelopeMixin) RawError() *types.ErrorResponse { return m.errResp }
