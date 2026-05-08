package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/security"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// KMS is the wrapper for an Aruba Cloud Key Management Service instance
// (a direct child of a Project). Construct with aruba.NewKMS() and bind it
// via IntoProject(project).
//
// Family A: regional, Metadata/Properties envelope, location-aware.
// Supports full CRUD. Create and Update share the same request type (KmsRequest).
//
// The KMS client also exposes Keys() and Kmips() to access nested raw clients
// for cryptographic keys and KMIP services.
//
// Path: /projects/{projectID}/providers/Aruba.Security/kms[/{kmsID}]
type KMS struct {
	errMixin
	metadataMixin
	regionalMixin
	projectScopedMixin
	responseMetadataMixin
	statusMixin
	httpEnvelopeMixin

	billingPeriod *BillingPeriod

	response *types.KmsResponse
}

// ---------------------------------------------------------------------------
// Standard setters
// ---------------------------------------------------------------------------

func (k *KMS) IntoProject(p Ref) *KMS                 { k.intoProject(p); return k }
func (k *KMS) WithName(n string) *KMS                 { k.withName(n); return k }
func (k *KMS) AddTag(t string) *KMS                   { k.addTag(t); return k }
func (k *KMS) RemoveTag(t string) *KMS                { k.removeTag(t); return k }
func (k *KMS) ReplaceTags(ts ...string) *KMS          { k.replaceTags(ts...); return k }
func (k *KMS) InRegion(region Region) *KMS            { k.inRegion(region); return k }
func (k *KMS) WithBillingPeriod(p BillingPeriod) *KMS { k.billingPeriod = &p; return k }

// ---------------------------------------------------------------------------
// Ref + ID accessors
// ---------------------------------------------------------------------------

func (k *KMS) URI() string   { return k.RespURI() }
func (k *KMS) KMSID() string { return k.ID() }

// ---------------------------------------------------------------------------
// Raw accessors
// ---------------------------------------------------------------------------

func (k *KMS) Raw() *types.KmsResponse      { return k.response }
func (k *KMS) RawRequest() types.KmsRequest { return k.toRequest() }

// ---------------------------------------------------------------------------
// Response-preferring accessors
// ---------------------------------------------------------------------------

func (k *KMS) BillingPeriod() BillingPeriod {
	if k.response != nil && k.response.Properties.BillingPeriod != "" {
		return k.response.Properties.BillingPeriod
	}
	if k.billingPeriod != nil {
		return *k.billingPeriod
	}
	return ""
}

// ---------------------------------------------------------------------------
// Wire conversions
// ---------------------------------------------------------------------------

func (k *KMS) toRequest() types.KmsRequest {
	props := types.KmsPropertiesRequest{}
	if k.billingPeriod != nil {
		props.BillingPeriod = *k.billingPeriod
	}
	return types.KmsRequest{
		Metadata: types.RegionalResourceMetadataRequest{
			ResourceMetadataRequest: k.toMetadata(),
			Location:                k.toLocation(),
		},
		Properties: props,
	}
}

func (k *KMS) fromResponse(resp *types.KmsResponse) {
	if resp == nil {
		return
	}
	k.response = resp
	k.setMeta(&resp.Metadata)
	k.withName(kmsDeref(resp.Metadata.Name))
	if len(resp.Metadata.Tags) > 0 {
		k.replaceTags(resp.Metadata.Tags...)
	}
	if resp.Metadata.LocationResponse != nil {
		k.inRegion(resp.Metadata.LocationResponse.Value)
	}
	k.setStatus(&resp.Status)
	k.setTerminalStates(kmsTerminalStates)

	if resp.Properties.BillingPeriod != "" {
		v := resp.Properties.BillingPeriod
		k.billingPeriod = &v
	}

	if resp.Metadata.ProjectResponseMetadata != nil && resp.Metadata.ProjectResponseMetadata.ID != "" {
		k.projectID = resp.Metadata.ProjectResponseMetadata.ID
	}
	if k.projectID == "" && k.RespURI() != "" {
		if pid := parseURIIDs(k.RespURI())["projects"]; pid != "" {
			k.projectID = pid
		}
	}
}

func kmsDeref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// kmsIDsFromRef
// ---------------------------------------------------------------------------

func kmsIDsFromRef(ref Ref) (projectID, kmsID string, err error) {
	kid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withKMSID); ok {
			return w.KMSID(), true
		}
		return "", false
	}, "kms")
	if !ok || kid == "" {
		return "", "", fmt.Errorf("cannot determine KMS ID from Ref %q", ref.URI())
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
	return pid, kid, nil
}

var kmsTerminalStates = map[string]bool{
	"Active": true,
	"Error":  false,
	"Failed": false,
}

// ---------------------------------------------------------------------------
// Low-level interface + adapter
// ---------------------------------------------------------------------------

type kmsLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KmsList], error)
	Get(ctx context.Context, projectID, kmsID string, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	Create(ctx context.Context, projectID string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	Update(ctx context.Context, projectID, kmsID string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error)
	Delete(ctx context.Context, projectID, kmsID string, params *types.RequestParameters) (*types.Response[any], error)
}

type kmsClientAdapter struct {
	low kmsLowLevelClient
}

func newKMSClientAdapter(rest *restclient.Client) *kmsClientAdapter {
	if rest == nil {
		return &kmsClientAdapter{}
	}
	return &kmsClientAdapter{
		low: security.NewKMSClientImpl(rest),
	}
}

func (a *kmsClientAdapter) Create(ctx context.Context, k *KMS, opts ...CallOption) (*KMS, error) {
	if err := k.Err(); err != nil {
		return k, err
	}
	if k.ProjectID() == "" {
		return k, fmt.Errorf("Create: KMS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, k.ProjectID(), k.toRequest(), rp)
	populateHTTPEnvelope(&k.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		k.fromResponse(resp.Data)
		k.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, k)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				k.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return k, err
	}
	if resp != nil && !resp.IsSuccess() {
		return k, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return k, nil
}

func (a *kmsClientAdapter) Update(ctx context.Context, k *KMS, opts ...CallOption) (*KMS, error) {
	if err := k.Err(); err != nil {
		return k, err
	}
	if k.KMSID() == "" {
		return k, fmt.Errorf("Update: KMS has no ID")
	}
	if k.ProjectID() == "" {
		return k, fmt.Errorf("Update: KMS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, k.ProjectID(), k.KMSID(), k.toRequest(), rp)
	populateHTTPEnvelope(&k.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		k.fromResponse(resp.Data)
		k.setRefresh(func(ctx context.Context) error {
			fresh, err := a.Get(ctx, k)
			if err != nil {
				return err
			}
			if fresh != nil && fresh.Raw() != nil {
				k.fromResponse(fresh.Raw())
			}
			return nil
		})
	}
	if err != nil {
		return k, err
	}
	if resp != nil && !resp.IsSuccess() {
		return k, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return k, nil
}

func (a *kmsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*KMS, error) {
	projectID, kmsID, err := kmsIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, kmsID, rp)
	out := &KMS{}
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

func (a *kmsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, kmsID, err := kmsIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, kmsID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *kmsClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*KMS], error) {
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
	var items []*KMS
	if resp != nil && resp.Data != nil {
		items = make([]*KMS, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			k := &KMS{}
			k.projectID = projectID
			k.fromResponse(&resp.Data.Values[i])
			k.setRefresh(func(ctx context.Context) error {
				fresh, err := a.Get(ctx, k)
				if err != nil {
					return err
				}
				if fresh != nil && fresh.Raw() != nil {
					k.fromResponse(fresh.Raw())
				}
				return nil
			})
			if k.projectID == "" {
				k.projectID = projectID
			}
			items = append(items, k)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*KMS], error) {
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
