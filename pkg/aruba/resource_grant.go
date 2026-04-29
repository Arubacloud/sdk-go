package aruba

import (
	"fmt"
	"time"

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
