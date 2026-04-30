package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/security"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// Key is the wrapper for a cryptographic key nested inside a KMS instance.
// Construct with aruba.NewKey() and bind via IntoKMS(parent).
//
// Family B: flat request (no Metadata/Properties boxing, no metadataMixin,
// no tags, no location).
//
// No Update operation. CRUD: Create / Get / Delete / List.
//
// Identity: KeyResponse carries no ResourceMetadataResponse; ID() and KeyID()
// read from KeyResponse.KeyID, and URI() is constructed from (projectID, kmsID, keyID).
type Key struct {
	errMixin
	kmsScopedMixin
	responseMetadataMixin // present but never populated; ID/URI shadowed below
	httpEnvelopeMixin

	name      *string
	algorithm *types.KeyAlgorithm
	response  *types.KeyResponse
}

// ---------------------------------------------------------------------------
// Standard setters
// ---------------------------------------------------------------------------

func (k *Key) IntoKMS(parent Ref) *Key     { k.intoKMS(parent); return k }
func (k *Key) WithName(n string) *Key      { k.name = &n; return k }
func (k *Key) WithAlgorithm(a string) *Key { v := types.KeyAlgorithm(a); k.algorithm = &v; return k }

// ---------------------------------------------------------------------------
// Ref + ID accessors (shadow responseMetadataMixin)
// ---------------------------------------------------------------------------

// ID returns the key's unique ID from the response, or "" before a Create/Get.
func (k *Key) ID() string {
	if k.response != nil && k.response.KeyID != nil {
		return *k.response.KeyID
	}
	return ""
}

// KeyID is an alias for ID() and satisfies withKeyID for future child wrappers.
func (k *Key) KeyID() string { return k.ID() }

// URI constructs the canonical path for this key.
// Returns "" if any of projectID, kmsID, or keyID is missing.
func (k *Key) URI() string {
	pid, kid, keyID := k.ProjectID(), k.KMSID(), k.ID()
	if pid == "" || kid == "" || keyID == "" {
		return ""
	}
	return fmt.Sprintf("/projects/%s/providers/Aruba.Security/kms/%s/keys/%s", pid, kid, keyID)
}

// ---------------------------------------------------------------------------
// Raw accessors
// ---------------------------------------------------------------------------

func (k *Key) Raw() *types.KeyResponse      { return k.response }
func (k *Key) RawRequest() types.KeyRequest { return k.toRequest() }

// ---------------------------------------------------------------------------
// Response-preferring read accessors
// ---------------------------------------------------------------------------

func (k *Key) Name() string {
	if k.response != nil && k.response.Name != nil {
		return *k.response.Name
	}
	return keyDeref(k.name)
}

func (k *Key) Algorithm() string {
	if k.response != nil && k.response.Algorithm != nil {
		return string(*k.response.Algorithm)
	}
	if k.algorithm != nil {
		return string(*k.algorithm)
	}
	return ""
}

func (k *Key) Type() string {
	if k.response != nil && k.response.Type != nil {
		return string(*k.response.Type)
	}
	return ""
}

func (k *Key) KeyStatus() string {
	if k.response != nil && k.response.Status != nil {
		return string(*k.response.Status)
	}
	return ""
}

func (k *Key) CreationSource() string {
	if k.response != nil && k.response.CreationSource != nil {
		return string(*k.response.CreationSource)
	}
	return ""
}

func (k *Key) PrivateKeyID() string {
	if k.response != nil && k.response.PrivateKeyID != nil {
		return *k.response.PrivateKeyID
	}
	return ""
}

// ---------------------------------------------------------------------------
// Wire conversions
// ---------------------------------------------------------------------------

func (k *Key) toRequest() types.KeyRequest {
	req := types.KeyRequest{}
	if k.name != nil {
		req.Name = *k.name
	}
	if k.algorithm != nil {
		req.Algorithm = *k.algorithm
	}
	return req
}

func (k *Key) fromResponse(resp *types.KeyResponse) {
	if resp == nil {
		return
	}
	k.response = resp
	if resp.Name != nil {
		v := *resp.Name
		k.name = &v
	}
	if resp.Algorithm != nil {
		v := *resp.Algorithm
		k.algorithm = &v
	}
}

func keyDeref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// keyIDsFromRef
// ---------------------------------------------------------------------------

func keyIDsFromRef(ref Ref) (projectID, kmsID, keyID string, err error) {
	keyID, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withKeyID); ok {
			return w.KeyID(), true
		}
		return "", false
	}, "keys")
	if !ok || keyID == "" {
		return "", "", "", fmt.Errorf("cannot determine key ID from Ref %q", ref.URI())
	}
	kmsID, ok = extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withKMSID); ok {
			return w.KMSID(), true
		}
		return "", false
	}, "kms")
	if !ok || kmsID == "" {
		return "", "", "", fmt.Errorf("cannot determine KMS ID from Ref %q", ref.URI())
	}
	projectID, ok = extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || projectID == "" {
		return "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return projectID, kmsID, keyID, nil
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type keysLowLevelClient interface {
	List(ctx context.Context, projectID, kmsID string, params *types.RequestParameters) (*types.Response[types.KeyList], error)
	Get(ctx context.Context, projectID, kmsID, keyID string, params *types.RequestParameters) (*types.Response[types.KeyResponse], error)
	Create(ctx context.Context, projectID, kmsID string, body types.KeyRequest, params *types.RequestParameters) (*types.Response[types.KeyResponse], error)
	Delete(ctx context.Context, projectID, kmsID, keyID string, params *types.RequestParameters) (*types.Response[any], error)
}

type keysClientAdapter struct{ low keysLowLevelClient }

func newKeysClientAdapter(rest *restclient.Client) *keysClientAdapter {
	if rest == nil {
		return &keysClientAdapter{}
	}
	return &keysClientAdapter{low: security.NewKeyClientImpl(rest)}
}

func (a *keysClientAdapter) Create(ctx context.Context, k *Key, opts ...CallOption) (*Key, error) {
	if err := k.Err(); err != nil {
		return k, err
	}
	if k.ProjectID() == "" {
		return k, fmt.Errorf("Create: Key has no parent project — call IntoKMS first")
	}
	if k.KMSID() == "" {
		return k, fmt.Errorf("Create: Key has no parent KMS — call IntoKMS first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, k.ProjectID(), k.KMSID(), k.toRequest(), rp)
	populateHTTPEnvelope(&k.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		k.fromResponse(resp.Data)
	}
	if err != nil {
		return k, err
	}
	if resp != nil && !resp.IsSuccess() {
		return k, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return k, nil
}

func (a *keysClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Key, error) {
	projectID, kmsID, keyID, err := keyIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, kmsID, keyID, rp)
	out := &Key{}
	out.projectID = projectID
	out.kmsID = kmsID
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
	}
	if err != nil {
		return out, err
	}
	if resp != nil && !resp.IsSuccess() {
		return out, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return out, nil
}

func (a *keysClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, kmsID, keyID, err := keyIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, kmsID, keyID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *keysClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*Key], error) {
	projectID, kmsID, err := kmsIDsFromRef(parent)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, kmsID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*Key
	if resp != nil && resp.Data != nil {
		items = make([]*Key, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			k := &Key{}
			k.projectID = projectID
			k.kmsID = kmsID
			k.fromResponse(&resp.Data.Values[i])
			items = append(items, k)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*Key], error) {
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
