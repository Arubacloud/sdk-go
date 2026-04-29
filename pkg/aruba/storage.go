package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/storage"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

type StorageClient interface {
	Snapshots() SnapshotsClient
	Volumes() VolumesClient
	Backups() StorageBackupsClient
	Restores() StorageRestoreClient
}

type storageClientImpl struct {
	snapshotsClient SnapshotsClient
	volumesClient   VolumesClient
	backupsClient   StorageBackupsClient
	restoresClient  StorageRestoreClient
}

var _ StorageClient = (*storageClientImpl)(nil)

func (c *storageClientImpl) Snapshots() SnapshotsClient {
	return c.snapshotsClient
}

func (c *storageClientImpl) Volumes() VolumesClient {
	return c.volumesClient
}

func (c *storageClientImpl) Backups() StorageBackupsClient {
	return c.backupsClient
}

func (c *storageClientImpl) Restores() StorageRestoreClient {
	return c.restoresClient
}

type SnapshotsClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*Snapshot], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*Snapshot, error)
	Create(ctx context.Context, snap *Snapshot, opts ...CallOption) (*Snapshot, error)
	Update(ctx context.Context, snap *Snapshot, opts ...CallOption) (*Snapshot, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

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

type VolumesClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*BlockStorage], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*BlockStorage, error)
	Create(ctx context.Context, vol *BlockStorage, opts ...CallOption) (*BlockStorage, error)
	Update(ctx context.Context, vol *BlockStorage, opts ...CallOption) (*BlockStorage, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type volumeLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.BlockStorageList], error)
	Get(ctx context.Context, projectID, volumeID string, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Create(ctx context.Context, projectID string, body types.BlockStorageRequest, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Update(ctx context.Context, projectID, volumeID string, body types.BlockStorageRequest, params *types.RequestParameters) (*types.Response[types.BlockStorageResponse], error)
	Delete(ctx context.Context, projectID, volumeID string, params *types.RequestParameters) (*types.Response[any], error)
}

type volumesClientAdapter struct{ low volumeLowLevelClient }

func newVolumesClientAdapter(rest *restclient.Client) *volumesClientAdapter {
	if rest == nil {
		return &volumesClientAdapter{}
	}
	return &volumesClientAdapter{low: storage.NewVolumesClientImpl(rest)}
}

func (a *volumesClientAdapter) Create(ctx context.Context, vol *BlockStorage, opts ...CallOption) (*BlockStorage, error) {
	if err := vol.Err(); err != nil {
		return vol, err
	}
	if vol.ProjectID() == "" {
		return vol, fmt.Errorf("Create: BlockStorage has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, vol.ProjectID(), vol.toRequest(), rp)
	populateHTTPEnvelope(&vol.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		vol.fromResponse(resp.Data)
	}
	if err != nil {
		return vol, err
	}
	if resp != nil && !resp.IsSuccess() {
		return vol, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return vol, nil
}

func (a *volumesClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*BlockStorage, error) {
	projectID, blockStorageID, err := blockStorageIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, blockStorageID, rp)
	out := &BlockStorage{}
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

func (a *volumesClientAdapter) Update(ctx context.Context, vol *BlockStorage, opts ...CallOption) (*BlockStorage, error) {
	if err := vol.Err(); err != nil {
		return vol, err
	}
	if vol.ID() == "" {
		return vol, fmt.Errorf("Update: BlockStorage has no ID — call Get first or seed from response metadata")
	}
	if vol.ProjectID() == "" {
		return vol, fmt.Errorf("Update: BlockStorage has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, vol.ProjectID(), vol.ID(), vol.toRequest(), rp)
	populateHTTPEnvelope(&vol.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		vol.fromResponse(resp.Data)
	}
	if err != nil {
		return vol, err
	}
	if resp != nil && !resp.IsSuccess() {
		return vol, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return vol, nil
}

func (a *volumesClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, blockStorageID, err := blockStorageIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, blockStorageID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *volumesClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*BlockStorage], error) {
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
	var items []*BlockStorage
	if resp != nil && resp.Data != nil {
		items = make([]*BlockStorage, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			bs := &BlockStorage{}
			bs.fromResponse(&resp.Data.Values[i])
			if bs.projectID == "" {
				bs.projectID = projectID
			}
			items = append(items, bs)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*BlockStorage], error) {
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

// blockStorageIDsFromRef extracts (projectID, blockStorageID) from a Ref.
func blockStorageIDsFromRef(ref Ref) (projectID, blockStorageID string, err error) {
	bid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withBlockStorageID); ok {
			return w.BlockStorageID(), true
		}
		return "", false
	}, "blockstorages")
	if !ok || bid == "" {
		return "", "", fmt.Errorf("cannot determine BlockStorage ID from Ref %q", ref.URI())
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
	return pid, bid, nil
}

type StorageBackupsClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.StorageBackupList], error)
	Get(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Update(ctx context.Context, projectID string, backupID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Create(ctx context.Context, projectID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Delete(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type StorageRestoreClient interface {
	List(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[types.RestoreList], error)
	Get(ctx context.Context, projectID string, backupID string, restoreID string, params *types.RequestParameters) (*types.Response[types.RestoreResponse], error)
	Update(ctx context.Context, projectID string, backupID string, restoreID string, body types.RestoreRequest, params *types.RequestParameters) (*types.Response[types.RestoreResponse], error)
	Create(ctx context.Context, projectID string, backupID string, body types.RestoreRequest, params *types.RequestParameters) (*types.Response[types.RestoreResponse], error)
	Delete(ctx context.Context, projectID string, backupID string, restoreID string, params *types.RequestParameters) (*types.Response[any], error)
}
