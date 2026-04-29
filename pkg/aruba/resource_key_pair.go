package aruba

import (
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// KeyPair is the wrapper for an Aruba Cloud Compute SSH Key Pair (a direct
// child of a Project). Construct with aruba.NewKeyPair() and bind it via
// IntoProject(project) and WithPublicKey(key).
type KeyPair struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	linkedMixin
	httpEnvelopeMixin

	publicKey *string

	response *types.KeyPairResponse
}

// Setters (chainable).

func (k *KeyPair) IntoProject(p Ref) *KeyPair          { k.intoProject(p); return k }
func (k *KeyPair) WithName(n string) *KeyPair          { k.withName(n); return k }
func (k *KeyPair) AddTag(t string) *KeyPair            { k.addTag(t); return k }
func (k *KeyPair) RemoveTag(t string) *KeyPair         { k.removeTag(t); return k }
func (k *KeyPair) ReplaceTags(ts ...string) *KeyPair   { k.replaceTags(ts...); return k }
func (k *KeyPair) WithLocation(loc string) *KeyPair    { k.withLocation(loc); return k }
func (k *KeyPair) InRegion(region string) *KeyPair     { k.withLocation(region); return k }

// WithPublicKey sets the SSH public key (mapped to wire field "value").
func (k *KeyPair) WithPublicKey(key string) *KeyPair {
	k.publicKey = &key
	return k
}

// URI satisfies Ref.
func (k *KeyPair) URI() string { return k.RespURI() }

// KeyPairID satisfies withKeyPairID.
func (k *KeyPair) KeyPairID() string { return k.ID() }

// Raw shadows responseMetadataMixin.Raw() with the typed key-pair response.
func (k *KeyPair) Raw() *types.KeyPairResponse { return k.response }

// RawRequest returns what toRequest() would emit right now.
func (k *KeyPair) RawRequest() types.KeyPairRequest { return k.toRequest() }

// PublicKey returns the SSH public key value ("" if unset). On a hydrated
// response wrapper this surfaces the response's Properties.Value.
func (k *KeyPair) PublicKey() string { return keyPairDerefString(k.publicKey) }

func (k *KeyPair) toRequest() types.KeyPairRequest {
	props := types.KeyPairPropertiesRequest{}
	if k.publicKey != nil {
		props.Value = *k.publicKey
	}
	return types.KeyPairRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: k.toMetadata(),
			Location:                k.toLocation(),
		},
		Properties: props,
	}
}

func (k *KeyPair) fromResponse(resp *types.KeyPairResponse) {
	if resp == nil {
		return
	}
	k.response = resp
	k.setMeta(&resp.Metadata)
	k.withName(keyPairDerefString(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		k.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		k.withLocation(resp.Metadata.LocationResponse.Value)
	}
	k.setLinked(resp.Properties.LinkedResources)

	if resp.Properties.Value != "" {
		v := resp.Properties.Value
		k.publicKey = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		k.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if k.projectID == "" && k.RespURI() != "" {
		ids := parseURIIDs(k.RespURI())
		k.projectID = ids["projects"]
	}
}

func keyPairDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
