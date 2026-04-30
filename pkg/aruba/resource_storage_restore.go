package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/storage"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// StorageRestore is the wrapper for an Aruba Cloud Storage Restore (a direct
// child of a StorageBackup, grandchild of a Project). Construct with
// aruba.NewStorageRestore() and bind it via IntoBackup(backup) and WithTarget(volume).
type StorageRestore struct {
	errMixin
	metadataMixin
	regionalMixin
	backupScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	targetRef *string // body URI (request "Target" / response "Destination")

	response *types.StorageRestoreResponse
}

func (r *StorageRestore) IntoBackup(b Ref) *StorageRestore         { r.intoBackup(b); return r }
func (r *StorageRestore) WithName(n string) *StorageRestore        { r.withName(n); return r }
func (r *StorageRestore) AddTag(t string) *StorageRestore          { r.addTag(t); return r }
func (r *StorageRestore) RemoveTag(t string) *StorageRestore       { r.removeTag(t); return r }
func (r *StorageRestore) ReplaceTags(ts ...string) *StorageRestore { r.replaceTags(ts...); return r }
func (r *StorageRestore) WithLocation(loc string) *StorageRestore  { r.withLocation(loc); return r }
func (r *StorageRestore) InRegion(region string) *StorageRestore   { r.withLocation(region); return r }

// WithTarget binds the destination volume (where the backup will be restored to)
// via its URI. Pass any Ref (typed or aruba.URI(...)). Empty URIs are recorded
// on the error sink and the field remains unset.
func (r *StorageRestore) WithTarget(vol Ref) *StorageRestore {
	uri := vol.URI()
	if uri == "" {
		r.addErr(fmt.Errorf("WithTarget: empty URI"))
		return r
	}
	r.targetRef = &uri
	return r
}

// URI satisfies Ref.
func (r *StorageRestore) URI() string { return r.RespURI() }

// RestoreID satisfies withRestoreID.
func (r *StorageRestore) RestoreID() string { return r.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed restore response.
func (r *StorageRestore) Raw() *types.StorageRestoreResponse { return r.response }

// RawRequest returns what toRequest() would emit right now.
func (r *StorageRestore) RawRequest() types.StorageRestoreRequest { return r.toRequest() }

// TargetURI returns the destination volume URI ("" if unset).
// On a hydrated response wrapper this surfaces the response's "Destination" field.
func (r *StorageRestore) TargetURI() string { return storageRestoreDerefString(r.targetRef) }

func (r *StorageRestore) toRequest() types.StorageRestoreRequest {
	props := types.StorageRestorePropertiesRequest{}
	if r.targetRef != nil {
		props.Target = types.ReferenceResource{URI: *r.targetRef}
	}
	return types.StorageRestoreRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: r.toMetadata(),
			Location:                r.toLocation(),
		},
		Properties: props,
	}
}

func (r *StorageRestore) fromResponse(resp *types.StorageRestoreResponse) {
	if resp == nil {
		return
	}
	r.response = resp
	r.setMeta(&resp.Metadata)
	r.withName(storageRestoreDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		r.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		r.withLocation(resp.Metadata.LocationResponse.Value)
	}
	r.setStatus(&resp.Status)

	// Response shape uses Destination (not Target).
	if resp.Properties.Destination.URI != "" {
		v := resp.Properties.Destination.URI
		r.targetRef = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		r.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if (r.projectID == "" || r.backupID == "") && r.RespURI() != "" {
		ids := parseURIIDs(r.RespURI())
		if r.projectID == "" {
			r.projectID = ids["projects"]
		}
		if r.backupID == "" {
			r.backupID = ids["backups"]
		}
	}
}

func storageRestoreDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// Low-level client interface, adapter, constructor, and methods
// ---------------------------------------------------------------------------

type storageRestoreLowLevelClient interface {
	List(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[types.StorageRestoreList], error)
	Get(ctx context.Context, projectID, backupID, restoreID string, params *types.RequestParameters) (*types.Response[types.StorageRestoreResponse], error)
	Create(ctx context.Context, projectID, backupID string, body types.StorageRestoreRequest, params *types.RequestParameters) (*types.Response[types.StorageRestoreResponse], error)
	Update(ctx context.Context, projectID, backupID, restoreID string, body types.StorageRestoreRequest, params *types.RequestParameters) (*types.Response[types.StorageRestoreResponse], error)
	Delete(ctx context.Context, projectID, backupID, restoreID string, params *types.RequestParameters) (*types.Response[any], error)
}

type storageRestoresClientAdapter struct{ low storageRestoreLowLevelClient }

func newStorageRestoresClientAdapter(rest *restclient.Client) *storageRestoresClientAdapter {
	if rest == nil {
		return &storageRestoresClientAdapter{}
	}
	// NewRestoreClientImpl panics if backupClient is nil — instantiate one internally
	// so the public adapter constructor stays single-arg.
	return &storageRestoresClientAdapter{
		low: storage.NewRestoreClientImpl(rest, storage.NewBackupClientImpl(rest)),
	}
}

func (a *storageRestoresClientAdapter) Create(ctx context.Context, r *StorageRestore, opts ...CallOption) (*StorageRestore, error) {
	if err := r.Err(); err != nil {
		return r, err
	}
	if r.BackupID() == "" || r.ProjectID() == "" {
		return r, fmt.Errorf("Create: StorageRestore has no parent backup — call IntoBackup first")
	}
	if r.targetRef == nil {
		return r, fmt.Errorf("Create: StorageRestore has no target — call WithTarget first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, r.ProjectID(), r.BackupID(), r.toRequest(), rp)
	populateHTTPEnvelope(&r.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		r.fromResponse(resp.Data)
	}
	if err != nil {
		return r, err
	}
	if resp != nil && !resp.IsSuccess() {
		return r, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return r, nil
}

func (a *storageRestoresClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*StorageRestore, error) {
	projectID, backupID, restoreID, err := restoreIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, backupID, restoreID, rp)
	out := &StorageRestore{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if out.projectID == "" {
		out.projectID = projectID
	}
	if out.backupID == "" {
		out.backupID = backupID
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *storageRestoresClientAdapter) Update(ctx context.Context, r *StorageRestore, opts ...CallOption) (*StorageRestore, error) {
	if err := r.Err(); err != nil {
		return r, err
	}
	if r.ID() == "" {
		return r, fmt.Errorf("Update: StorageRestore has no ID — call Get first or seed from response metadata")
	}
	if r.BackupID() == "" || r.ProjectID() == "" {
		return r, fmt.Errorf("Update: StorageRestore has no parent backup — call IntoBackup first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, r.ProjectID(), r.BackupID(), r.ID(), r.toRequest(), rp)
	populateHTTPEnvelope(&r.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		r.fromResponse(resp.Data)
	}
	if err != nil {
		return r, err
	}
	if resp != nil && !resp.IsSuccess() {
		return r, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return r, nil
}

func (a *storageRestoresClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, backupID, restoreID, err := restoreIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, backupID, restoreID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *storageRestoresClientAdapter) List(ctx context.Context, backup Ref, opts ...CallOption) (*List[*StorageRestore], error) {
	projectID, backupID, err := backupIDsFromRef(backup)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, backupID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*StorageRestore
	if resp != nil && resp.Data != nil {
		items = make([]*StorageRestore, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			v := &StorageRestore{}
			v.fromResponse(&resp.Data.Values[i])
			if v.projectID == "" {
				v.projectID = projectID
			}
			if v.backupID == "" {
				v.backupID = backupID
			}
			items = append(items, v)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*StorageRestore], error) {
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

// restoreIDsFromRef extracts (projectID, backupID, restoreID) from a Ref.
func restoreIDsFromRef(ref Ref) (projectID, backupID, restoreID string, err error) {
	rid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withRestoreID); ok {
			return w.RestoreID(), true
		}
		return "", false
	}, "restores")
	if !ok || rid == "" {
		return "", "", "", fmt.Errorf("cannot determine StorageRestore ID from Ref %q", ref.URI())
	}
	bid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withBackupID); ok {
			return w.BackupID(), true
		}
		return "", false
	}, "backups")
	if !ok || bid == "" {
		return "", "", "", fmt.Errorf("cannot determine backup ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, bid, rid, nil
}
