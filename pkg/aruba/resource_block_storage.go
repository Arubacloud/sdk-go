package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/storage"
	"github.com/Arubacloud/sdk-go/internal/restclient"
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
func (b *BlockStorage) WithSizeGB(gb int) *BlockStorage        { b.sizeGB = gb; return b }
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

func (b *BlockStorage) SizeGB() int { return b.sizeGB }
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
	b.setTerminalStates(blockStorageTerminalStates)
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

var blockStorageTerminalStates = map[string]bool{
	"NotUsed": true,
	"InUse":   true,
	"Used":    true,
	"Error":   false,
	"Failed": false,
}

// WaitUntilNotUsed blocks until the BlockStorage reaches the "NotUsed" state —
// the steady terminal state for an unattached volume. Call this after Create
// and before passing the volume to a CloudServer.
func (b *BlockStorage) WaitUntilNotUsed(ctx context.Context, opts ...WaitOption) error {
	return b.WaitUntilStates(ctx, []string{"NotUsed"}, opts...)
}

// WaitUntilUsed blocks until the BlockStorage reaches the "InUse" or "Used"
// state — both signal that the volume has been attached to a CloudServer. The
// platform may emit either value; this method succeeds on whichever arrives.
func (b *BlockStorage) WaitUntilUsed(ctx context.Context, opts ...WaitOption) error {
	return b.WaitUntilStates(ctx, []string{"InUse", "Used"}, opts...)
}

// ---------------------------------------------------------------------------
// Low-level client interface, adapter, constructor, and methods
// ---------------------------------------------------------------------------

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
		vol.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, vol)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				vol.fromResponse(fresh.Raw())
			}
			return nil
		})
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
		vol.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, vol)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				vol.fromResponse(fresh.Raw())
			}
			return nil
		})
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
			bs.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, bs)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					bs.fromResponse(fresh.Raw())
				}
				return nil
			})
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
