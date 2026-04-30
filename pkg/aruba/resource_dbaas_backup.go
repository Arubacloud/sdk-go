package aruba

import (
	"fmt"

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

func (b *DBaaSBackup) Size() int32 {
	if b.response == nil {
		return 0
	}
	return b.response.Properties.Storage.Size
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
