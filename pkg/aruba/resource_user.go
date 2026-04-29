package aruba

import (
	"fmt"
	"time"

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
