package aruba

import (
	"context"
	"fmt"
	"time"

	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// User is the wrapper for an Aruba Cloud DBaaS user (a child of a DBaaS
// instance). Construct with aruba.NewUser() and bind via IntoDBaaS(parent).
//
// Family B: flat request (no Metadata/Properties boxing, no metadataMixin,
// no tags, no location).
//
// Identity: UserResponse carries no `id` field; the Username IS the path
// identifier. ID() returns the locally-set username, and URI() is
// constructed from (projectID, dbaasID, username).
//
// Password is a write-only field: WithPassword stores it locally for use in
// Create/Update wire bodies, but the wrapper deliberately exposes no
// Password() accessor. The response struct (UserResponse) contains no
// Password field, so hydration cannot accidentally surface it either.
type User struct {
	errMixin
	dbaasScopedMixin
	responseMetadataMixin
	httpEnvelopeMixin

	username *string
	password *string
	response *types.UserResponse
}

// Setters.

func (u *User) IntoDBaaS(parent Ref) *User     { u.intoDBaaS(parent); return u }
func (u *User) WithUsername(name string) *User { u.username = &name; return u }
func (u *User) WithPassword(pw string) *User   { u.password = &pw; return u }

// Ref + ID accessors.

// ID returns the user's username (which serves as its path identifier).
// Shadows responseMetadataMixin.ID() since the response has no id field.
func (u *User) ID() string { return userDerefString(u.username) }

// URI constructs the canonical URI from (projectID, dbaasID, username).
// Returns "" if any component is missing.
func (u *User) URI() string {
	pid, did, name := u.ProjectID(), u.DBaaSID(), u.ID()
	if pid == "" || did == "" || name == "" {
		return ""
	}
	return fmt.Sprintf("/projects/%s/providers/Aruba.Database/dbaas/%s/users/%s", pid, did, name)
}

// Read accessors.

// Username returns the username from the response if available, else from the
// locally-set value.
func (u *User) Username() string {
	if u.response != nil && u.response.Username != "" {
		return u.response.Username
	}
	return userDerefString(u.username)
}

// CreatedAt returns the user creation time from the response.
func (u *User) CreatedAt() time.Time {
	if u.response != nil && u.response.CreationDate != nil {
		return *u.response.CreationDate
	}
	return time.Time{}
}

// CreatedBy returns the identity that created this user.
func (u *User) CreatedBy() string {
	if u.response != nil && u.response.CreatedBy != nil {
		return *u.response.CreatedBy
	}
	return ""
}

func (u *User) Raw() *types.UserResponse { return u.response }

// RawRequest returns the wire body that would be sent on Create/Update. It
// includes the locally-set password if WithPassword was called — by design,
// for parity with other wrappers' RawRequest debugging surface. There is no
// Password() accessor on *User; the password is intentionally not exposed
// through any read-only path other than this wire mirror.
func (u *User) RawRequest() types.UserRequest { return u.toRequest() }

// Wire conversions.

func (u *User) toRequest() types.UserRequest {
	return types.UserRequest{
		Username: userDerefString(u.username),
		Password: userDerefString(u.password),
	}
}

func (u *User) fromResponse(resp *types.UserResponse) {
	if resp == nil {
		return
	}
	u.response = resp
	if resp.Username != "" {
		v := resp.Username
		u.username = &v
	}
	// Do not touch u.password — UserResponse has no Password field, and
	// the locally-set password must survive hydration so a subsequent
	// Update can still send it on the wire.
}

func userDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// Users low-level client, adapter, and helpers
// ---------------------------------------------------------------------------

// userIDsFromRef extracts (projectID, dbaasID, userID) from a Ref.
func userIDsFromRef(ref Ref) (projectID, dbaasID, userID string, err error) {
	name, ok := extractID(ref, func(r Ref) (string, bool) {
		return "", false // no withUserID interface; rely on URI segment fallback
	}, "users")
	if !ok || name == "" {
		return "", "", "", fmt.Errorf("cannot determine user ID from Ref %q", ref.URI())
	}
	did, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDBaaSID); ok {
			return w.DBaaSID(), true
		}
		return "", false
	}, "dbaas")
	if !ok || did == "" {
		return "", "", "", fmt.Errorf("cannot determine DBaaS ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, did, name, nil
}

type usersLowLevelClient interface {
	List(ctx context.Context, projectID, dbaasID string, params *types.RequestParameters) (*types.Response[types.UserList], error)
	Get(ctx context.Context, projectID, dbaasID, userID string, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Create(ctx context.Context, projectID, dbaasID string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Update(ctx context.Context, projectID, dbaasID, userID string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Delete(ctx context.Context, projectID, dbaasID, userID string, params *types.RequestParameters) (*types.Response[any], error)
}

type usersClientAdapter struct{ low usersLowLevelClient }

func newUsersClientAdapter(rest *restclient.Client) *usersClientAdapter {
	if rest == nil {
		return &usersClientAdapter{}
	}
	return &usersClientAdapter{low: database.NewUsersClientImpl(rest)}
}

func (a *usersClientAdapter) Create(ctx context.Context, u *User, opts ...CallOption) (*User, error) {
	if err := u.Err(); err != nil {
		return u, err
	}
	if u.ProjectID() == "" {
		return u, fmt.Errorf("Create: User has no parent project — call IntoDBaaS first")
	}
	if u.DBaaSID() == "" {
		return u, fmt.Errorf("Create: User has no parent DBaaS — call IntoDBaaS first")
	}
	if u.Username() == "" {
		return u, fmt.Errorf("Create: User has no username — call WithUsername first")
	}
	if u.password == nil {
		return u, fmt.Errorf("Create: password is required — call WithPassword first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, u.ProjectID(), u.DBaaSID(), u.toRequest(), rp)
	populateHTTPEnvelope(&u.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		u.fromResponse(resp.Data)
	}
	if err != nil {
		return u, err
	}
	if resp != nil && !resp.IsSuccess() {
		return u, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return u, nil
}

func (a *usersClientAdapter) Update(ctx context.Context, u *User, opts ...CallOption) (*User, error) {
	if err := u.Err(); err != nil {
		return u, err
	}
	if u.ID() == "" {
		return u, fmt.Errorf("Update: User has no ID — call WithUsername first")
	}
	if u.DBaaSID() == "" {
		return u, fmt.Errorf("Update: User has no parent DBaaS — call IntoDBaaS first")
	}
	if u.ProjectID() == "" {
		return u, fmt.Errorf("Update: User has no parent project — call IntoDBaaS first")
	}
	if u.password == nil {
		return u, fmt.Errorf("Update: password is required — call WithPassword first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, u.ProjectID(), u.DBaaSID(), u.ID(), u.toRequest(), rp)
	populateHTTPEnvelope(&u.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		u.fromResponse(resp.Data)
	}
	if err != nil {
		return u, err
	}
	if resp != nil && !resp.IsSuccess() {
		return u, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return u, nil
}

func (a *usersClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*User, error) {
	projectID, dbaasID, userID, err := userIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, dbaasID, userID, rp)
	out := &User{}
	out.dbaasID = dbaasID
	out.projectID = projectID
	name := userID
	out.username = &name
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

func (a *usersClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, dbaasID, userID, err := userIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, dbaasID, userID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *usersClientAdapter) List(ctx context.Context, dbaas Ref, opts ...CallOption) (*List[*User], error) {
	projectID, dbaasID, err := dbaasIDsFromRef(dbaas)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, dbaasID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*User
	if resp != nil && resp.Data != nil {
		items = make([]*User, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			u := &User{}
			u.dbaasID = dbaasID
			u.projectID = projectID
			u.fromResponse(&resp.Data.Values[i])
			items = append(items, u)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*User], error) {
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
