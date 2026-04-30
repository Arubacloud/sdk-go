package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/storage"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// Snapshot is the wrapper for an Aruba Cloud Snapshot (a direct child of a
// Project, derived from a BlockStorage volume). Construct with
// aruba.NewSnapshot() and bind it via IntoProject(project) and OfVolume(bs).
type Snapshot struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	billingPeriod *string
	volumeRef     *string // body URI

	// Read-only fields hydrated from response.
	sizeGB      *int32
	storageType *types.BlockStorageType
	zone        *string
	bootable    *bool

	response *types.SnapshotResponse
}

func (s *Snapshot) IntoProject(p Ref) *Snapshot          { s.intoProject(p); return s }
func (s *Snapshot) WithName(n string) *Snapshot          { s.withName(n); return s }
func (s *Snapshot) AddTag(t string) *Snapshot            { s.addTag(t); return s }
func (s *Snapshot) RemoveTag(t string) *Snapshot         { s.removeTag(t); return s }
func (s *Snapshot) ReplaceTags(ts ...string) *Snapshot   { s.replaceTags(ts...); return s }
func (s *Snapshot) WithLocation(loc string) *Snapshot    { s.withLocation(loc); return s }
func (s *Snapshot) InRegion(region string) *Snapshot     { s.withLocation(region); return s }
func (s *Snapshot) WithBillingPeriod(p string) *Snapshot { s.billingPeriod = &p; return s }

// OfVolume binds the source BlockStorage via its URI. Pass any Ref (typed or
// aruba.URI(...)). Empty URIs are recorded on the error sink and the field
// remains unset.
func (s *Snapshot) OfVolume(vol Ref) *Snapshot {
	uri := vol.URI()
	if uri == "" {
		s.addErr(fmt.Errorf("OfVolume: empty URI"))
		return s
	}
	s.volumeRef = &uri
	return s
}

// URI satisfies Ref.
func (s *Snapshot) URI() string { return s.RespURI() }

// SnapshotID satisfies withSnapshotID.
func (s *Snapshot) SnapshotID() string { return s.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed Snapshot response.
func (s *Snapshot) Raw() *types.SnapshotResponse { return s.response }

// RawRequest returns what toRequest() would emit right now.
func (s *Snapshot) RawRequest() types.SnapshotRequest { return s.toRequest() }

func (s *Snapshot) BillingPeriod() string { return snapshotDerefString(s.billingPeriod) }
func (s *Snapshot) VolumeURI() string     { return snapshotDerefString(s.volumeRef) }

// Read-only accessors hydrated from response.
func (s *Snapshot) Size() int32 {
	if s.sizeGB == nil {
		return 0
	}
	return *s.sizeGB
}
func (s *Snapshot) Type() types.BlockStorageType {
	if s.storageType == nil {
		return ""
	}
	return *s.storageType
}
func (s *Snapshot) Zone() string { return snapshotDerefString(s.zone) }
func (s *Snapshot) Bootable() bool {
	if s.bootable == nil {
		return false
	}
	return *s.bootable
}

func (s *Snapshot) toRequest() types.SnapshotRequest {
	props := types.SnapshotPropertiesRequest{
		BillingPeriod: s.billingPeriod,
	}
	if s.volumeRef != nil {
		props.Volume = types.ReferenceResource{URI: *s.volumeRef}
	}
	return types.SnapshotRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: s.toMetadata(),
			Location:                s.toLocation(),
		},
		Properties: props,
	}
}

func (s *Snapshot) fromResponse(resp *types.SnapshotResponse) {
	if resp == nil {
		return
	}
	s.response = resp
	s.setMeta(&resp.Metadata)
	s.withName(snapshotDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		s.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		s.withLocation(resp.Metadata.LocationResponse.Value)
	}
	s.setStatus(&resp.Status)

	if resp.Properties.SizeGB != nil {
		v := *resp.Properties.SizeGB
		s.sizeGB = &v
	}
	if resp.Properties.BillingPeriod != nil && *resp.Properties.BillingPeriod != "" {
		v := *resp.Properties.BillingPeriod
		s.billingPeriod = &v
	}
	if resp.Properties.Type != "" {
		v := resp.Properties.Type
		s.storageType = &v
	}
	if resp.Properties.Zone != "" {
		v := resp.Properties.Zone
		s.zone = &v
	}
	if resp.Properties.Bootable != nil {
		v := *resp.Properties.Bootable
		s.bootable = &v
	}
	if resp.Properties.Volume != nil && resp.Properties.Volume.URI != nil && *resp.Properties.Volume.URI != "" {
		v := *resp.Properties.Volume.URI
		s.volumeRef = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		s.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if s.projectID == "" && s.RespURI() != "" {
		if pid := parseURIIDs(s.RespURI())["projects"]; pid != "" {
			s.projectID = pid
		}
	}
}

func snapshotDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// Low-level client interface, adapter, constructor, and methods
// ---------------------------------------------------------------------------

type snapshotLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.SnapshotList], error)
	Get(ctx context.Context, projectID, snapshotID string, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Create(ctx context.Context, projectID string, body types.SnapshotRequest, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Update(ctx context.Context, projectID, snapshotID string, body types.SnapshotRequest, params *types.RequestParameters) (*types.Response[types.SnapshotResponse], error)
	Delete(ctx context.Context, projectID, snapshotID string, params *types.RequestParameters) (*types.Response[any], error)
}

type snapshotsClientAdapter struct{ low snapshotLowLevelClient }

func newSnapshotsClientAdapter(rest *restclient.Client) *snapshotsClientAdapter {
	if rest == nil {
		return &snapshotsClientAdapter{}
	}
	return &snapshotsClientAdapter{
		low: storage.NewSnapshotsClientImpl(rest, storage.NewVolumesClientImpl(rest)),
	}
}

func (a *snapshotsClientAdapter) Create(ctx context.Context, snap *Snapshot, opts ...CallOption) (*Snapshot, error) {
	if err := snap.Err(); err != nil {
		return snap, err
	}
	if snap.ProjectID() == "" {
		return snap, fmt.Errorf("Create: Snapshot has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, snap.ProjectID(), snap.toRequest(), rp)
	populateHTTPEnvelope(&snap.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		snap.fromResponse(resp.Data)
	}
	if err != nil {
		return snap, err
	}
	if resp != nil && !resp.IsSuccess() {
		return snap, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return snap, nil
}

func (a *snapshotsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Snapshot, error) {
	projectID, snapshotID, err := snapshotIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, snapshotID, rp)
	out := &Snapshot{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if out.projectID == "" {
		out.projectID = projectID
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *snapshotsClientAdapter) Update(ctx context.Context, snap *Snapshot, opts ...CallOption) (*Snapshot, error) {
	if err := snap.Err(); err != nil {
		return snap, err
	}
	if snap.ID() == "" {
		return snap, fmt.Errorf("Update: Snapshot has no ID — call Get first or seed from response metadata")
	}
	if snap.ProjectID() == "" {
		return snap, fmt.Errorf("Update: Snapshot has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, snap.ProjectID(), snap.ID(), snap.toRequest(), rp)
	populateHTTPEnvelope(&snap.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		snap.fromResponse(resp.Data)
	}
	if err != nil {
		return snap, err
	}
	if resp != nil && !resp.IsSuccess() {
		return snap, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return snap, nil
}

func (a *snapshotsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, snapshotID, err := snapshotIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, snapshotID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *snapshotsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*Snapshot], error) {
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
	var items []*Snapshot
	if resp != nil && resp.Data != nil {
		items = make([]*Snapshot, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			snap := &Snapshot{}
			snap.fromResponse(&resp.Data.Values[i])
			if snap.projectID == "" {
				snap.projectID = projectID
			}
			items = append(items, snap)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*Snapshot], error) {
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

// snapshotIDsFromRef extracts (projectID, snapshotID) from a Ref.
func snapshotIDsFromRef(ref Ref) (projectID, snapshotID string, err error) {
	sid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withSnapshotID); ok {
			return w.SnapshotID(), true
		}
		return "", false
	}, "snapshots")
	if !ok || sid == "" {
		return "", "", fmt.Errorf("cannot determine Snapshot ID from Ref %q", ref.URI())
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
	return pid, sid, nil
}
