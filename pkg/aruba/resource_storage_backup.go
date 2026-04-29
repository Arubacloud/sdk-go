package aruba

import (
	"fmt"

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
	billingPeriod *string

	response *types.StorageBackupResponse
}

func (b *StorageBackup) IntoProject(p Ref) *StorageBackup        { b.intoProject(p); return b }
func (b *StorageBackup) WithName(n string) *StorageBackup        { b.withName(n); return b }
func (b *StorageBackup) AddTag(t string) *StorageBackup          { b.addTag(t); return b }
func (b *StorageBackup) RemoveTag(t string) *StorageBackup       { b.removeTag(t); return b }
func (b *StorageBackup) ReplaceTags(ts ...string) *StorageBackup { b.replaceTags(ts...); return b }
func (b *StorageBackup) WithLocation(loc string) *StorageBackup  { b.withLocation(loc); return b }
func (b *StorageBackup) InRegion(region string) *StorageBackup   { b.withLocation(region); return b }

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

func (b *StorageBackup) WithBillingPeriod(p string) *StorageBackup {
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

func (b *StorageBackup) OriginURI() string     { return storageBackupDerefString(b.originRef) }
func (b *StorageBackup) BillingPeriod() string { return storageBackupDerefString(b.billingPeriod) }

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
