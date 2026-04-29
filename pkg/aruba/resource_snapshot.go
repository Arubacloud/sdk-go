package aruba

import (
	"fmt"

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
