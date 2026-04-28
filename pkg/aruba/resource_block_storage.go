package aruba

import (
	"fmt"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// BlockStorage is the wrapper for an Aruba Cloud BlockStorage volume
// (a direct child of a Project). Construct with aruba.NewBlockStorage()
// and bind it via IntoProject(project).
type BlockStorage struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	linkedMixin
	httpEnvelopeMixin

	sizeGB        int
	storageType   *types.BlockStorageType
	zone          *string
	billingPeriod *string
	snapshotRef   *string // body URI
	image         *string
	bootable      *bool

	response *types.BlockStorageResponse
}

func (b *BlockStorage) IntoProject(p Ref) *BlockStorage        { b.intoProject(p); return b }
func (b *BlockStorage) WithName(n string) *BlockStorage        { b.withName(n); return b }
func (b *BlockStorage) AddTag(t string) *BlockStorage          { b.addTag(t); return b }
func (b *BlockStorage) RemoveTag(t string) *BlockStorage       { b.removeTag(t); return b }
func (b *BlockStorage) ReplaceTags(ts ...string) *BlockStorage { b.replaceTags(ts...); return b }
func (b *BlockStorage) WithLocation(loc string) *BlockStorage  { b.withLocation(loc); return b }
func (b *BlockStorage) InRegion(region string) *BlockStorage   { b.withLocation(region); return b }
func (b *BlockStorage) InZone(zone string) *BlockStorage       { b.zone = &zone; return b }
func (b *BlockStorage) WithSize(gb int) *BlockStorage          { b.sizeGB = gb; return b }
func (b *BlockStorage) WithType(t types.BlockStorageType) *BlockStorage {
	b.storageType = &t
	return b
}
func (b *BlockStorage) WithBillingPeriod(p string) *BlockStorage { b.billingPeriod = &p; return b }
func (b *BlockStorage) WithImage(img string) *BlockStorage       { b.image = &img; return b }
func (b *BlockStorage) WithBootable(boot bool) *BlockStorage     { b.bootable = &boot; return b }

// FromSnapshot binds the source snapshot via its URI. Pass any Ref (typed or
// aruba.URI(...)). Empty URIs are recorded on the error sink and the field
// remains unset.
func (b *BlockStorage) FromSnapshot(snap Ref) *BlockStorage {
	uri := snap.URI()
	if uri == "" {
		b.addErr(fmt.Errorf("FromSnapshot: empty URI"))
		return b
	}
	b.snapshotRef = &uri
	return b
}

// URI satisfies Ref.
func (b *BlockStorage) URI() string { return b.RespURI() }

// BlockStorageID satisfies withBlockStorageID.
func (b *BlockStorage) BlockStorageID() string { return b.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed BlockStorage response.
func (b *BlockStorage) Raw() *types.BlockStorageResponse { return b.response }

// RawRequest returns what toRequest() would emit right now.
func (b *BlockStorage) RawRequest() types.BlockStorageRequest { return b.toRequest() }

func (b *BlockStorage) Size() int { return b.sizeGB }
func (b *BlockStorage) Type() types.BlockStorageType {
	if b.storageType == nil {
		return ""
	}
	return *b.storageType
}
func (b *BlockStorage) Zone() string          { return blockStorageDerefString(b.zone) }
func (b *BlockStorage) BillingPeriod() string { return blockStorageDerefString(b.billingPeriod) }
func (b *BlockStorage) Image() string         { return blockStorageDerefString(b.image) }
func (b *BlockStorage) Bootable() bool {
	if b.bootable == nil {
		return false
	}
	return *b.bootable
}
func (b *BlockStorage) SnapshotURI() string { return blockStorageDerefString(b.snapshotRef) }

func (b *BlockStorage) toRequest() types.BlockStorageRequest {
	var bp string
	if b.billingPeriod != nil {
		bp = *b.billingPeriod
	}
	var t types.BlockStorageType
	if b.storageType != nil {
		t = *b.storageType
	}
	props := types.BlockStoragePropertiesRequest{
		SizeGB:        b.sizeGB,
		BillingPeriod: bp,
		Zone:          b.zone,
		Type:          t,
		Bootable:      b.bootable,
		Image:         b.image,
	}
	if b.snapshotRef != nil {
		props.Snapshot = &types.ReferenceResource{URI: *b.snapshotRef}
	}
	return types.BlockStorageRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: b.toMetadata(),
			Location:                b.toLocation(),
		},
		Properties: props,
	}
}

func (b *BlockStorage) fromResponse(resp *types.BlockStorageResponse) {
	if resp == nil {
		return
	}
	b.response = resp
	b.setMeta(&resp.Metadata)
	b.withName(blockStorageDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		b.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		b.withLocation(resp.Metadata.LocationResponse.Value)
	}
	b.setStatus(&resp.Status)
	b.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.SizeGB != 0 {
		b.sizeGB = resp.Properties.SizeGB
	}
	if resp.Properties.Type != "" {
		v := resp.Properties.Type
		b.storageType = &v
	}
	if resp.Properties.Zone != "" {
		v := resp.Properties.Zone
		b.zone = &v
	}
	if resp.Properties.BillingPeriod != "" {
		v := resp.Properties.BillingPeriod
		b.billingPeriod = &v
	}
	if resp.Properties.Image != nil && *resp.Properties.Image != "" {
		v := *resp.Properties.Image
		b.image = &v
	}
	if resp.Properties.Bootable != nil {
		v := *resp.Properties.Bootable
		b.bootable = &v
	}
	if resp.Properties.Snapshot != nil && resp.Properties.Snapshot.URI != "" {
		v := resp.Properties.Snapshot.URI
		b.snapshotRef = &v
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

func blockStorageDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
