package aruba

import (
	"context"
	"fmt"
	"time"

	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

// Grant is the wrapper for an Aruba Cloud DBaaS grant (a child of a Database).
// Construct with aruba.NewGrant() and bind via IntoDatabase(parent).
//
// Family B: flat request (no Metadata/Properties boxing, no metadataMixin,
// no tags, no location).
//
// Identity quirk: GrantResponse carries no `id` field, but the path segment
// /grants/<grantID> uses a server-supplied opaque ID. Practical consequences:
//
//   - After Create, ID() stays empty — the wire response cannot reveal it.
//     A subsequent Update/Delete on the same wrapper will fail pre-flight.
//     Discover the new grant via List and use the typed *Grant from there,
//     or call Get with a URI Ref carrying /grants/<id>.
//   - List items have empty ID() and empty URI() for the same reason.
//   - Get with URI(".../grants/<id>") populates ID() from the URL segment.
//   - Update requires the wrapper to already carry ID() — typically from Get.
type Grant struct {
	errMixin
	databaseScopedMixin
	responseMetadataMixin
	httpEnvelopeMixin

	id       *string
	userName *string
	roleName *string
	response *types.GrantResponse
}

// Setters.

func (g *Grant) IntoDatabase(parent Ref) *Grant  { g.intoDatabase(parent); return g }
func (g *Grant) WithUserName(name string) *Grant { g.userName = &name; return g }
func (g *Grant) WithRoleName(name string) *Grant { g.roleName = &name; return g }

// Ref + ID accessors.

// ID returns the opaque server-supplied grantID. See type docstring for when
// this is and isn't populated. Shadows responseMetadataMixin.ID() since the
// response has no id field.
func (g *Grant) ID() string { return grantDerefString(g.id) }

// URI constructs the canonical URI from (projectID, dbaasID, databaseID, ID).
// Returns "" if any component is missing.
func (g *Grant) URI() string {
	pid, did, dbid, gid := g.ProjectID(), g.DBaaSID(), g.DatabaseID(), g.ID()
	if pid == "" || did == "" || dbid == "" || gid == "" {
		return ""
	}
	return fmt.Sprintf(
		"/projects/%s/providers/Aruba.Database/dbaas/%s/databases/%s/grants/%s",
		pid, did, dbid, gid,
	)
}

// Read accessors.

// UserName returns the username from the response if available, else from the
// locally-set value.
func (g *Grant) UserName() string {
	if g.response != nil && g.response.User.Username != "" {
		return g.response.User.Username
	}
	return grantDerefString(g.userName)
}

// RoleName returns the role name from the response if available, else from the
// locally-set value.
func (g *Grant) RoleName() string {
	if g.response != nil && g.response.Role.Name != "" {
		return g.response.Role.Name
	}
	return grantDerefString(g.roleName)
}

// DatabaseName returns the database name from the response.
func (g *Grant) DatabaseName() string {
	if g.response != nil {
		return g.response.Database.Name
	}
	return ""
}

// CreatedAt returns the grant creation time from the response.
func (g *Grant) CreatedAt() time.Time {
	if g.response != nil && g.response.CreationDate != nil {
		return *g.response.CreationDate
	}
	return time.Time{}
}

// CreatedBy returns the identity that created this grant.
func (g *Grant) CreatedBy() string {
	if g.response != nil && g.response.CreatedBy != nil {
		return *g.response.CreatedBy
	}
	return ""
}

func (g *Grant) Raw() *types.GrantResponse      { return g.response }
func (g *Grant) RawRequest() types.GrantRequest { return g.toRequest() }

// Wire conversions.

func (g *Grant) toRequest() types.GrantRequest {
	return types.GrantRequest{
		User: types.GrantUser{Username: grantDerefString(g.userName)},
		Role: types.GrantRole{Name: grantDerefString(g.roleName)},
	}
}

func (g *Grant) fromResponse(resp *types.GrantResponse) {
	if resp == nil {
		return
	}
	g.response = resp
	if resp.User.Username != "" {
		v := resp.User.Username
		g.userName = &v
	}
	if resp.Role.Name != "" {
		v := resp.Role.Name
		g.roleName = &v
	}
	// Do not touch g.id — GrantResponse has no id field. The opaque grantID
	// is set by the adapter Get from a URI Ref's path segment, never here.
}

func grantDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// ---------------------------------------------------------------------------
// Grants low-level client, adapter, and helpers
// ---------------------------------------------------------------------------

// grantIDsFromRef extracts (projectID, dbaasID, databaseID, grantID) from a Ref.
func grantIDsFromRef(ref Ref) (projectID, dbaasID, databaseID, grantID string, err error) {
	gid, ok := extractID(ref, func(r Ref) (string, bool) {
		return "", false // no withGrantID interface; rely on URI segment fallback
	}, "grants")
	if !ok || gid == "" {
		return "", "", "", "", fmt.Errorf("cannot determine grant ID from Ref %q", ref.URI())
	}
	dbid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDatabaseID); ok {
			return w.DatabaseID(), true
		}
		return "", false
	}, "databases")
	if !ok || dbid == "" {
		return "", "", "", "", fmt.Errorf("cannot determine database ID from Ref %q", ref.URI())
	}
	did, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDBaaSID); ok {
			return w.DBaaSID(), true
		}
		return "", false
	}, "dbaas")
	if !ok || did == "" {
		return "", "", "", "", fmt.Errorf("cannot determine DBaaS ID from Ref %q", ref.URI())
	}
	pid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withProjectID); ok {
			return w.ProjectID(), true
		}
		return "", false
	}, "projects")
	if !ok || pid == "" {
		return "", "", "", "", fmt.Errorf("cannot determine project ID from Ref %q", ref.URI())
	}
	return pid, did, dbid, gid, nil
}

type grantsLowLevelClient interface {
	List(ctx context.Context, projectID, dbaasID, databaseID string, params *types.RequestParameters) (*types.Response[types.GrantList], error)
	Get(ctx context.Context, projectID, dbaasID, databaseID, grantID string, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Create(ctx context.Context, projectID, dbaasID, databaseID string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Update(ctx context.Context, projectID, dbaasID, databaseID, grantID string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Delete(ctx context.Context, projectID, dbaasID, databaseID, grantID string, params *types.RequestParameters) (*types.Response[any], error)
}

type grantsClientAdapter struct{ low grantsLowLevelClient }

func newGrantsClientAdapter(rest *restclient.Client) *grantsClientAdapter {
	if rest == nil {
		return &grantsClientAdapter{}
	}
	return &grantsClientAdapter{low: database.NewGrantsClientImpl(rest)}
}

func (a *grantsClientAdapter) Create(ctx context.Context, g *Grant, opts ...CallOption) (*Grant, error) {
	if err := g.Err(); err != nil {
		return g, err
	}
	if g.ProjectID() == "" {
		return g, fmt.Errorf("Create: Grant has no parent project — call IntoDatabase first")
	}
	if g.DBaaSID() == "" {
		return g, fmt.Errorf("Create: Grant has no parent DBaaS — call IntoDatabase first")
	}
	if g.DatabaseID() == "" {
		return g, fmt.Errorf("Create: Grant has no parent database — call IntoDatabase first")
	}
	if g.userName == nil {
		return g, fmt.Errorf("Create: Grant has no username — call WithUserName first")
	}
	if g.roleName == nil {
		return g, fmt.Errorf("Create: Grant has no role — call WithRoleName first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, g.ProjectID(), g.DBaaSID(), g.DatabaseID(), g.toRequest(), rp)
	populateHTTPEnvelope(&g.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		g.fromResponse(resp.Data)
	}
	if err != nil {
		return g, err
	}
	if resp != nil && !resp.IsSuccess() {
		return g, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return g, nil
}

func (a *grantsClientAdapter) Update(ctx context.Context, g *Grant, opts ...CallOption) (*Grant, error) {
	if err := g.Err(); err != nil {
		return g, err
	}
	if g.ID() == "" {
		return g, fmt.Errorf("Update: Grant has no ID — get the grant via Get first to obtain the opaque ID")
	}
	if g.DatabaseID() == "" {
		return g, fmt.Errorf("Update: Grant has no parent database — call IntoDatabase first")
	}
	if g.DBaaSID() == "" {
		return g, fmt.Errorf("Update: Grant has no parent DBaaS — call IntoDatabase first")
	}
	if g.ProjectID() == "" {
		return g, fmt.Errorf("Update: Grant has no parent project — call IntoDatabase first")
	}
	if g.userName == nil {
		return g, fmt.Errorf("Update: Grant has no username — call WithUserName first")
	}
	if g.roleName == nil {
		return g, fmt.Errorf("Update: Grant has no role — call WithRoleName first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, g.ProjectID(), g.DBaaSID(), g.DatabaseID(), g.ID(), g.toRequest(), rp)
	populateHTTPEnvelope(&g.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		g.fromResponse(resp.Data)
	}
	if err != nil {
		return g, err
	}
	if resp != nil && !resp.IsSuccess() {
		return g, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return g, nil
}

func (a *grantsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Grant, error) {
	projectID, dbaasID, databaseID, grantID, err := grantIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, dbaasID, databaseID, grantID, rp)
	out := &Grant{}
	out.databaseID = databaseID
	out.dbaasID = dbaasID
	out.projectID = projectID
	out.id = &grantID
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

func (a *grantsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, dbaasID, databaseID, grantID, err := grantIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, dbaasID, databaseID, grantID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *grantsClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*Grant], error) {
	projectID, dbaasID, databaseID, err := databaseIDsFromRef(parent)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.List(ctx, projectID, dbaasID, databaseID, rp)
	if err != nil {
		return nil, err
	}
	if resp != nil && !resp.IsSuccess() {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	var items []*Grant
	if resp != nil && resp.Data != nil {
		items = make([]*Grant, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			g := &Grant{}
			g.databaseID = databaseID
			g.dbaasID = dbaasID
			g.projectID = projectID
			g.fromResponse(&resp.Data.Values[i])
			items = append(items, g)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*Grant], error) {
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
