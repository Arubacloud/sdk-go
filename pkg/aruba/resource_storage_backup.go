package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/storage"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// StorageBackup is the wrapper for an Aruba Cloud Storage Backup (a direct
// child of a Project, derived from a BlockStorage volume). Construct with
// aruba.NewStorageBackup() and bind it via IntoProject(project) and
// WithOrigin(bs).
type StorageBackup struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	backupType    *types.StorageBackupType
	originRef     *string // body URI
	retentionDays *int
	billingPeriod *BillingPeriod

	response *types.StorageBackupResponse
}

func (b *StorageBackup) IntoProject(p Ref) *StorageBackup        { b.intoProject(p); return b }
func (b *StorageBackup) WithName(n string) *StorageBackup        { b.withName(n); return b }
func (b *StorageBackup) AddTag(t string) *StorageBackup          { b.addTag(t); return b }
func (b *StorageBackup) RemoveTag(t string) *StorageBackup       { b.removeTag(t); return b }
func (b *StorageBackup) ReplaceTags(ts ...string) *StorageBackup { b.replaceTags(ts...); return b }
func (b *StorageBackup) WithLocation(loc Region) *StorageBackup  { b.withLocation(loc); return b }
func (b *StorageBackup) InRegion(region Region) *StorageBackup   { b.withLocation(region); return b }

// WithType sets the backup type (Full or Incremental).
func (b *StorageBackup) WithType(t types.StorageBackupType) *StorageBackup {
	v := t
	b.backupType = &v
	return b
}

// WithRetentionDays sets the number of days the backup is retained.
func (b *StorageBackup) WithRetentionDays(days int) *StorageBackup {
	v := days
	b.retentionDays = &v
	return b
}

func (b *StorageBackup) WithBillingPeriod(p BillingPeriod) *StorageBackup {
	b.billingPeriod = &p
	return b
}

// WithOrigin binds the source BlockStorage via its URI. Pass any Ref (typed
// or aruba.URI(...)). Empty URIs are recorded on the error sink and the field
// remains unset.
func (b *StorageBackup) WithOrigin(vol Ref) *StorageBackup {
	uri := vol.URI()
	if uri == "" {
		b.addErr(fmt.Errorf("WithOrigin: empty URI"))
		return b
	}
	b.originRef = &uri
	return b
}

// URI satisfies Ref.
func (b *StorageBackup) URI() string { return b.RespURI() }

// BackupID satisfies withBackupID.
func (b *StorageBackup) BackupID() string { return b.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed response.
func (b *StorageBackup) Raw() *types.StorageBackupResponse { return b.response }

// RawRequest returns what toRequest() would emit right now.
func (b *StorageBackup) RawRequest() types.StorageBackupRequest { return b.toRequest() }

func (b *StorageBackup) Type() types.StorageBackupType {
	if b.backupType == nil {
		return ""
	}
	return *b.backupType
}

func (b *StorageBackup) OriginURI() string { return storageBackupDerefString(b.originRef) }
func (b *StorageBackup) BillingPeriod() BillingPeriod {
	if b.billingPeriod == nil {
		return ""
	}
	return *b.billingPeriod
}

func (b *StorageBackup) RetentionDays() int {
	if b.retentionDays == nil {
		return 0
	}
	return *b.retentionDays
}

func (b *StorageBackup) toRequest() types.StorageBackupRequest {
	props := types.StorageBackupPropertiesRequest{
		RetentionDays: b.retentionDays,
		BillingPeriod: b.billingPeriod,
	}
	if b.backupType != nil {
		props.StorageBackupType = *b.backupType
	}
	if b.originRef != nil {
		props.Origin = types.ReferenceResource{URI: *b.originRef}
	}
	return types.StorageBackupRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: b.toMetadata(),
			Location:                b.toLocation(),
		},
		Properties: props,
	}
}

func (b *StorageBackup) fromResponse(resp *types.StorageBackupResponse) {
	if resp == nil {
		return
	}
	b.response = resp
	b.setMeta(&resp.Metadata)
	b.withName(storageBackupDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		b.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		b.withLocation(resp.Metadata.LocationResponse.Value)
	}
	b.setStatus(&resp.Status)
	b.setTerminalStates(storageBackupTerminalStates)

	if resp.Properties.Type != "" {
		v := resp.Properties.Type
		b.backupType = &v
	}
	if resp.Properties.Origin.URI != "" {
		v := resp.Properties.Origin.URI
		b.originRef = &v
	}
	if resp.Properties.RetentionDays != nil {
		v := *resp.Properties.RetentionDays
		b.retentionDays = &v
	}
	if resp.Properties.BillingPeriod != nil && *resp.Properties.BillingPeriod != "" {
		v := *resp.Properties.BillingPeriod
		b.billingPeriod = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		b.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if b.projectID == "" && b.RespURI() != "" {
		if pid := parseURIIDs(b.RespURI())["projects"]; pid != "" {
			b.projectID = pid
		}
	}
}

func storageBackupDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var storageBackupTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// Low-level client interface, adapter, constructor, and methods
// ---------------------------------------------------------------------------

type storageBackupLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.StorageBackupList], error)
	Get(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Create(ctx context.Context, projectID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Update(ctx context.Context, projectID, backupID string, body types.StorageBackupRequest, params *types.RequestParameters) (*types.Response[types.StorageBackupResponse], error)
	Delete(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type storageBackupsClientAdapter struct{ low storageBackupLowLevelClient }

func newStorageBackupsClientAdapter(rest *restclient.Client) *storageBackupsClientAdapter {
	if rest == nil {
		return &storageBackupsClientAdapter{}
	}
	return &storageBackupsClientAdapter{low: storage.NewBackupClientImpl(rest)}
}

func (a *storageBackupsClientAdapter) Create(ctx context.Context, b *StorageBackup, opts ...CallOption) (*StorageBackup, error) {
	if err := b.Err(); err != nil {
		return b, err
	}
	if b.ProjectID() == "" {
		return b, fmt.Errorf("Create: StorageBackup has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, b.ProjectID(), b.toRequest(), rp)
	populateHTTPEnvelope(&b.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		b.fromResponse(resp.Data)
		b.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, b)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				b.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return b, err
	}
	if resp != nil && !resp.IsSuccess() {
		return b, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return b, nil
}

func (a *storageBackupsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*StorageBackup, error) {
	projectID, backupID, err := backupIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, backupID, rp)
	out := &StorageBackup{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
		out.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, out)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				out.fromResponse(fresh.Raw())
			}
			return nil
		})
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

func (a *storageBackupsClientAdapter) Update(ctx context.Context, b *StorageBackup, opts ...CallOption) (*StorageBackup, error) {
	if err := b.Err(); err != nil {
		return b, err
	}
	if b.ID() == "" {
		return b, fmt.Errorf("Update: StorageBackup has no ID — call Get first or seed from response metadata")
	}
	if b.ProjectID() == "" {
		return b, fmt.Errorf("Update: StorageBackup has no project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, b.ProjectID(), b.ID(), b.toRequest(), rp)
	populateHTTPEnvelope(&b.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		b.fromResponse(resp.Data)
		b.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, b)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				b.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return b, err
	}
	if resp != nil && !resp.IsSuccess() {
		return b, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return b, nil
}

func (a *storageBackupsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, backupID, err := backupIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, backupID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *storageBackupsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*StorageBackup], error) {
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
	var items []*StorageBackup
	if resp != nil && resp.Data != nil {
		items = make([]*StorageBackup, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			bkp := &StorageBackup{}
			bkp.fromResponse(&resp.Data.Values[i])
			bkp.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, bkp)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					bkp.fromResponse(fresh.Raw())
				}
				return nil
			})
			if bkp.projectID == "" {
				bkp.projectID = projectID
			}
			items = append(items, bkp)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*StorageBackup], error) {
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

// backupIDsFromRef extracts (projectID, backupID) from a Ref.
func backupIDsFromRef(ref Ref) (projectID, backupID string, err error) {
	bid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withBackupID); ok {
			return w.BackupID(), true
		}
		return "", false
	}, "backups")
	if !ok || bid == "" {
		return "", "", fmt.Errorf("cannot determine StorageBackup ID from Ref %q", ref.URI())
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
