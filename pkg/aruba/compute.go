package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/compute"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

type ComputeClient interface {
	CloudServers() CloudServersClient
	KeyPairs() KeyPairsClient
}

type computeClientImpl struct {
	cloudServerClient CloudServersClient
	keyPairClient     KeyPairsClient
}

var _ ComputeClient = (*computeClientImpl)(nil)

func (c *computeClientImpl) CloudServers() CloudServersClient {
	return c.cloudServerClient
}

func (c *computeClientImpl) KeyPairs() KeyPairsClient {
	return c.keyPairClient
}

type CloudServersClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.CloudServerList], error)
	Get(ctx context.Context, projectID string, cloudServerID string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Create(ctx context.Context, projectID string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Update(ctx context.Context, projectID string, cloudServerID string, body types.CloudServerRequest, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	Delete(ctx context.Context, projectID string, cloudServerID string, params *types.RequestParameters) (*types.Response[any], error)
	PowerOn(ctx context.Context, projectID string, cloudServerID string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	PowerOff(ctx context.Context, projectID string, cloudServerID string, params *types.RequestParameters) (*types.Response[types.CloudServerResponse], error)
	SetPassword(ctx context.Context, projectID string, cloudServerID string, body types.CloudServerPasswordRequest, params *types.RequestParameters) (*types.Response[any], error)
}

type KeyPairsClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*KeyPair], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*KeyPair, error)
	Create(ctx context.Context, kp *KeyPair, opts ...CallOption) (*KeyPair, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type keyPairLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KeyPairListResponse], error)
	Get(ctx context.Context, projectID, keyPairID string, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error)
	Create(ctx context.Context, projectID string, body types.KeyPairRequest, params *types.RequestParameters) (*types.Response[types.KeyPairResponse], error)
	Delete(ctx context.Context, projectID, keyPairID string, params *types.RequestParameters) (*types.Response[any], error)
}

type keyPairsClientAdapter struct{ low keyPairLowLevelClient }

func newKeyPairsClientAdapter(rest *restclient.Client) *keyPairsClientAdapter {
	if rest == nil {
		return &keyPairsClientAdapter{}
	}
	return &keyPairsClientAdapter{low: compute.NewKeyPairsClientImpl(rest)}
}

func (a *keyPairsClientAdapter) Create(ctx context.Context, kp *KeyPair, opts ...CallOption) (*KeyPair, error) {
	if err := kp.Err(); err != nil {
		return kp, err
	}
	if kp.ProjectID() == "" {
		return kp, fmt.Errorf("Create: KeyPair has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, kp.ProjectID(), kp.toRequest(), rp)
	populateHTTPEnvelope(&kp.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		kp.fromResponse(resp.Data)
	}
	if err != nil {
		return kp, err
	}
	if resp != nil && !resp.IsSuccess() {
		return kp, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return kp, nil
}

func (a *keyPairsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*KeyPair, error) {
	projectID, keyPairID, err := keyPairIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, keyPairID, rp)
	out := &KeyPair{}
	populateHTTPEnvelope(&out.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		out.fromResponse(resp.Data)
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

func (a *keyPairsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, keyPairID, err := keyPairIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, keyPairID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *keyPairsClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*KeyPair], error) {
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
	var items []*KeyPair
	if resp != nil && resp.Data != nil {
		items = make([]*KeyPair, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			kp := &KeyPair{}
			kp.fromResponse(&resp.Data.Values[i])
			if kp.projectID == "" {
				kp.projectID = projectID
			}
			items = append(items, kp)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*KeyPair], error) {
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

// keyPairIDsFromRef extracts (projectID, keyPairID) from a Ref.
func keyPairIDsFromRef(ref Ref) (projectID, keyPairID string, err error) {
	kid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withKeyPairID); ok {
			return w.KeyPairID(), true
		}
		return "", false
	}, "keypairs")
	if !ok || kid == "" {
		return "", "", fmt.Errorf("cannot determine KeyPair ID from Ref %q", ref.URI())
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
