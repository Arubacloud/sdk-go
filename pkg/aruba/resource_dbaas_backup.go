package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// DBaaSBackup is the wrapper for an Aruba Cloud DBaaS Backup (a direct child
// of a Project; the source DBaaS and Database are body-refs, not path-parents).
// Construct with aruba.NewDBaaSBackup() and bind it via IntoProject(project),
// WithDBaaS(d), and WithDatabase(db).
//
// Family A: regional, Metadata/Properties envelope, location-aware. No Update
// operation — the underlying API exposes only Create / Get / List / Delete.
//
// Path: /projects/{projectID}/providers/Aruba.Database/backups[/{backupID}]
type DBaaSBackup struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	billingPeriod *string
	dbaasRef      *string // body URI
	databaseRef   *string // body URI

	response *types.BackupResponse
}

func (b *DBaaSBackup) IntoProject(p Ref) *DBaaSBackup        { b.intoProject(p); return b }
func (b *DBaaSBackup) WithName(n string) *DBaaSBackup        { b.withName(n); return b }
func (b *DBaaSBackup) AddTag(t string) *DBaaSBackup          { b.addTag(t); return b }
func (b *DBaaSBackup) RemoveTag(t string) *DBaaSBackup       { b.removeTag(t); return b }
func (b *DBaaSBackup) ReplaceTags(ts ...string) *DBaaSBackup { b.replaceTags(ts...); return b }
func (b *DBaaSBackup) WithLocation(loc string) *DBaaSBackup  { b.withLocation(loc); return b }
func (b *DBaaSBackup) InRegion(region string) *DBaaSBackup   { b.withLocation(region); return b }

// WithDBaaS binds the source DBaaS via its URI. Empty URIs are recorded on the
// error sink and the field remains unset.
func (b *DBaaSBackup) WithDBaaS(d Ref) *DBaaSBackup {
	uri := d.URI()
	if uri == "" {
		b.addErr(fmt.Errorf("WithDBaaS: empty URI"))
		return b
	}
	b.dbaasRef = &uri
	return b
}

// WithDatabase binds the source Database via its URI. Empty URIs are recorded
// on the error sink and the field remains unset.
func (b *DBaaSBackup) WithDatabase(db Ref) *DBaaSBackup {
	uri := db.URI()
	if uri == "" {
		b.addErr(fmt.Errorf("WithDatabase: empty URI"))
		return b
	}
	b.databaseRef = &uri
	return b
}

func (b *DBaaSBackup) WithBillingPeriod(p string) *DBaaSBackup {
	b.billingPeriod = &p
	return b
}

// URI satisfies Ref.
func (b *DBaaSBackup) URI() string { return b.RespURI() }

// DBaaSBackupID satisfies withDBaaSBackupID.
func (b *DBaaSBackup) DBaaSBackupID() string { return b.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed response.
func (b *DBaaSBackup) Raw() *types.BackupResponse { return b.response }

// RawRequest returns what toRequest() would emit right now.
func (b *DBaaSBackup) RawRequest() types.BackupRequest { return b.toRequest() }

func (b *DBaaSBackup) BillingPeriod() string { return dbaasBackupDerefString(b.billingPeriod) }
func (b *DBaaSBackup) DBaaSURI() string      { return dbaasBackupDerefString(b.dbaasRef) }
func (b *DBaaSBackup) DatabaseURI() string   { return dbaasBackupDerefString(b.databaseRef) }

func (b *DBaaSBackup) SizeGB() int {
	if b.response == nil {
		return 0
	}
	return int(b.response.Properties.Storage.Size)
}

func (b *DBaaSBackup) Zone() string {
	if b.response == nil {
		return ""
	}
	return b.response.Properties.Zone
}

func (b *DBaaSBackup) toRequest() types.BackupRequest {
	props := types.BackupPropertiesRequest{
		Zone: b.Region(), // auto-derive from regionalMixin (no separate setter)
	}
	if b.dbaasRef != nil {
		props.DBaaS = types.ReferenceResource{URI: *b.dbaasRef}
	}
	if b.databaseRef != nil {
		props.Database = types.ReferenceResource{URI: *b.databaseRef}
	}
	if b.billingPeriod != nil {
		props.BillingPlan = types.BillingPeriodResource{BillingPeriod: *b.billingPeriod}
	}
	return types.BackupRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: b.toMetadata(),
			Location:                b.toLocation(),
		},
		Properties: props,
	}
}

func (b *DBaaSBackup) fromResponse(resp *types.BackupResponse) {
	if resp == nil {
		return
	}
	b.response = resp
	b.setMeta(&resp.Metadata)
	b.withName(dbaasBackupDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		b.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		b.withLocation(resp.Metadata.LocationResponse.Value)
	}
	b.setStatus(&resp.Status)
	b.setTerminalStates(dbaasBackupTerminalStates)

	if resp.Properties.DBaaS.URI != "" {
		v := resp.Properties.DBaaS.URI
		b.dbaasRef = &v
	}
	if resp.Properties.Database.URI != "" {
		v := resp.Properties.Database.URI
		b.databaseRef = &v
	}
	if resp.Properties.BillingPlan.BillingPeriod != "" {
		v := resp.Properties.BillingPlan.BillingPeriod
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

func dbaasBackupDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

var dbaasBackupTerminalStates = map[string]bool{
	"Available": true,
	"Error":     false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// DBaaS Backups low-level client, adapter, and helpers
// ---------------------------------------------------------------------------

// dbaasBackupIDsFromRef extracts (projectID, backupID) from a Ref. Uses the
// dedicated withDBaaSBackupID interface for typed extraction so a typed
// *StorageBackup Ref does not silently route to the DBaaS endpoint. URI
// fallback (segment "backups") remains inherently ambiguous between the two
// backup scopes — callers must pass URIs from the correct domain.
func dbaasBackupIDsFromRef(ref Ref) (projectID, backupID string, err error) {
	bid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDBaaSBackupID); ok {
			return w.DBaaSBackupID(), true
		}
		return "", false
	}, "backups")
	if !ok || bid == "" {
		return "", "", fmt.Errorf("cannot determine backup ID from Ref %q", ref.URI())
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

type dbaasBackupsLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.BackupList], error)
	Get(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Create(ctx context.Context, projectID string, body types.BackupRequest, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Delete(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type dbaasBackupsClientAdapter struct{ low dbaasBackupsLowLevelClient }

func newDBaaSBackupsClientAdapter(rest *restclient.Client) *dbaasBackupsClientAdapter {
	if rest == nil {
		return &dbaasBackupsClientAdapter{}
	}
	return &dbaasBackupsClientAdapter{low: database.NewBackupsClientImpl(rest)}
}

func (a *dbaasBackupsClientAdapter) Create(ctx context.Context, b *DBaaSBackup, opts ...CallOption) (*DBaaSBackup, error) {
	if err := b.Err(); err != nil {
		return b, err
	}
	if b.ProjectID() == "" {
		return b, fmt.Errorf("Create: DBaaSBackup has no parent project — call IntoProject first")
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

func (a *dbaasBackupsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*DBaaSBackup, error) {
	projectID, backupID, err := dbaasBackupIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, backupID, rp)
	out := &DBaaSBackup{}
	out.projectID = projectID
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

func (a *dbaasBackupsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, backupID, err := dbaasBackupIDsFromRef(ref)
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

func (a *dbaasBackupsClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*DBaaSBackup], error) {
	projectID, err := projectIDFromRef(parent)
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
	var items []*DBaaSBackup
	if resp != nil && resp.Data != nil {
		items = make([]*DBaaSBackup, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			b := &DBaaSBackup{}
			b.projectID = projectID
			b.fromResponse(&resp.Data.Values[i])
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
			if b.projectID == "" {
				b.projectID = projectID
			}
			items = append(items, b)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*DBaaSBackup], error) {
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
