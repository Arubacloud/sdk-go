package aruba

import (
	"fmt"

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
