package aruba

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/async"
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

// --------------------------------------------------------------------------
// metadataMixin — resource name and tags
// --------------------------------------------------------------------------

type metadataMixin struct {
	name string
	tags []string
}

func (m *metadataMixin) named(name string) {
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

// Name returns the name set via named.
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
	region Region
}

func (m *regionalMixin) inRegion(region Region) { m.region = region }

// Region returns the region value.
func (m *regionalMixin) Region() Region { return m.region }

func (m *regionalMixin) toLocation() types.LocationRequest {
	return types.LocationRequest{Value: m.region}
}

// --------------------------------------------------------------------------
// zonalMixin — resource zone (extends regionalMixin)
// --------------------------------------------------------------------------

// zonalMixin extends regionalMixin with zone tracking. Zones are always within
// a region (e.g. "ITBG-1" lives in region "ITBG"), so zonalMixin embeds
// regionalMixin and inherits its setter/getter/toLocation helper.
//
// The zone wire field is NOT part of types.LocationRequest — every zonal
// resource carries it on its own *PropertiesRequest under JSON tag "dataCenter".
// This mixin therefore only owns the value; each wrapper's toRequest() reads it
// via Zone() (for required Zone wire fields) or zonePtr() (for *Zone omitempty
// fields) and places it itself.
type zonalMixin struct {
	regionalMixin
	zone *Zone
}

func (m *zonalMixin) inZone(z Zone) { m.zone = &z }

// Zone returns the configured zone, or "" if InZone was never called.
func (m *zonalMixin) Zone() Zone {
	if m.zone == nil {
		return ""
	}
	return *m.zone
}

// zonePtr returns the underlying *Zone for resources whose wire field is
// *Zone with omitempty (e.g. BlockStorage, DBaaS). Returns nil if InZone
// was not called.
func (m *zonalMixin) zonePtr() *Zone { return m.zone }

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
		// Production URI uses "vpnTunnels" (camelCase); mixin/test form uses "vpn-tunnels".
		if v := parseURIIDs(parent.URI())["vpnTunnels"]; v != "" {
			tunnelID = v
			ok = true
		}
	}
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
		// Production URI uses "vpcPeerings" (camelCase); mixin/test form uses "peerings".
		if v := parseURIIDs(parent.URI())["vpcPeerings"]; v != "" {
			peeringID = v
			ok = true
		}
	}
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

// WaitOption configures WaitUntilActive / WaitUntilStates behaviour.
type WaitOption func(*waitOptions)

type waitOptions struct {
	retries   int
	baseDelay time.Duration
	timeout   time.Duration
}

func defaultWaitOptions() waitOptions {
	return waitOptions{
		retries:   async.DefaultRetries,
		baseDelay: async.DefaultBaseDelay,
		timeout:   async.DefaultTimeout,
	}
}

// WithRetries sets the maximum number of polling attempts (default: 60).
func WithRetries(n int) WaitOption { return func(o *waitOptions) { o.retries = n } }

// WithBaseDelay sets the fixed delay between polling attempts (default: 10s).
func WithBaseDelay(d time.Duration) WaitOption { return func(o *waitOptions) { o.baseDelay = d } }

// WithTimeout sets the overall deadline for the polling loop (default: 600s).
func WithTimeout(d time.Duration) WaitOption { return func(o *waitOptions) { o.timeout = d } }

func applyWaitOptions(opts []WaitOption) waitOptions {
	out := defaultWaitOptions()
	for _, opt := range opts {
		if opt != nil {
			opt(&out)
		}
	}
	return out
}

type statusMixin struct {
	status  *types.ResourceStatus
	refresh func(ctx context.Context) error
}

func (m *statusMixin) setStatus(s *types.ResourceStatus) { m.status = s }

func (m *statusMixin) setRefresh(fn func(context.Context) error) { m.refresh = fn }

// State returns the current lifecycle state, or the zero State ("").
func (m *statusMixin) State() types.State {
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

// PreviousState returns the previous lifecycle state, or the zero State ("").
func (m *statusMixin) PreviousState() types.State {
	if m.status == nil || m.status.PreviousStatus == nil || m.status.PreviousStatus.State == nil {
		return ""
	}
	return *m.status.PreviousStatus.State
}

// WaitUntilActive blocks until the resource reaches the "Active" state.
// Equivalent to WaitUntilStates(ctx, []State{StateActive}, opts...).
func (m *statusMixin) WaitUntilActive(ctx context.Context, opts ...WaitOption) error {
	return m.WaitUntilStates(ctx, []types.State{types.StateActive}, opts...)
}

// WaitUntilReady blocks until the resource reaches any healthy settled state.
// Use this when a caller does not care which steady state the resource lands in
// — only that it is no longer transitioning. Succeeds on Active, Running,
// Stopped, NotUsed, Reserved, InUse, or Used.
func (m *statusMixin) WaitUntilReady(ctx context.Context, opts ...WaitOption) error {
	return m.WaitUntilStates(ctx, []types.State{
		types.StateActive,
		types.StateRunning,
		types.StateStopped,
		types.StateNotUsed,
		types.StateReserved,
		types.StateInUse,
		types.StateUsed,
	}, opts...)
}

// WaitUntilStates blocks until the resource reaches any of the given target states.
//
// The check applies four rules in order:
//  1. state ∈ targets → success.
//  2. state.IsFailure() → terminal error.
//  3. state == "" || state.IsTransitory() → keep polling.
//  4. otherwise (settled, non-target) → terminal error.
//
// Rule 4 makes wait semantics context-dependent: a resource that settles in
// "Reserved" succeeds for a waiter that lists Reserved as a target and fails
// fast for one that does not. Returns a descriptive error if the refresh
// callback was not set (resource not produced by an adapter).
func (m *statusMixin) WaitUntilStates(ctx context.Context, targets []types.State, opts ...WaitOption) error {
	if m.refresh == nil {
		return errors.New("WaitUntilStates: refresh callback not set; resource must be produced by an adapter (Create/Get/Update/List) to support polling")
	}
	cfg := applyWaitOptions(opts)
	call := func(ctx context.Context) (*types.Response[any], error) {
		if err := m.refresh(ctx); err != nil {
			return nil, err
		}
		return &types.Response[any]{}, nil
	}
	var terminalErr error
	check := func(_ *types.Response[any]) (bool, error) {
		state := m.State()
		for _, t := range targets {
			if state == t {
				return true, nil
			}
		}
		if state.IsFailure() {
			terminalErr = fmt.Errorf("resource entered failure state %q (targets %v)", state, targets)
			return true, terminalErr
		}
		if state == "" || state.IsTransitory() {
			return false, nil
		}
		// settled, non-target, non-failure
		terminalErr = fmt.Errorf("resource settled in state %q which is not a wait target %v", state, targets)
		return true, terminalErr
	}
	_, err := async.WaitFor[any](ctx, cfg.retries, cfg.baseDelay, cfg.timeout, call, check).Await(ctx)
	if terminalErr != nil {
		return terminalErr
	}
	return err
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
