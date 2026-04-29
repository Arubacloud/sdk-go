package aruba

import (
	"fmt"
	"time"

	"github.com/Arubacloud/sdk-go/pkg/types"
)

// Database is the wrapper for an Aruba Cloud database (a child of a DBaaS
// instance). Construct with aruba.NewDatabase() and bind via IntoDBaaS(parent).
//
// Family B: flat request (no Metadata/Properties boxing, no metadataMixin,
// no tags, no location).
//
// Identity: DatabaseResponse carries no `id` field; the Name IS the path
// identifier. ID() and DatabaseID() return the locally-set name, and URI()
// is constructed from (projectID, dbaasID, name).
type Database struct {
	errMixin
	dbaasScopedMixin
	responseMetadataMixin
	httpEnvelopeMixin

	name     *string
	response *types.DatabaseResponse
}

// Setters.

func (d *Database) IntoDBaaS(parent Ref) *Database { d.intoDBaaS(parent); return d }
func (d *Database) WithName(name string) *Database { d.name = &name; return d }

// Ref + ID accessors.

// ID returns the database's name (which serves as its path identifier).
// Shadows responseMetadataMixin.ID() since the response has no separate id field.
func (d *Database) ID() string { return dbDerefString(d.name) }

// DatabaseID is an alias for ID() and satisfies withDatabaseID for child wrappers.
func (d *Database) DatabaseID() string { return d.ID() }

// URI constructs the canonical URI from (projectID, dbaasID, name).
// Returns "" if any component is missing.
func (d *Database) URI() string {
	pid, did, name := d.ProjectID(), d.DBaaSID(), d.ID()
	if pid == "" || did == "" || name == "" {
		return ""
	}
	return fmt.Sprintf("/projects/%s/providers/Aruba.Database/dbaas/%s/databases/%s", pid, did, name)
}

// Read accessors.

// Name returns the database name from the response if available, else from the
// locally-set value.
func (d *Database) Name() string {
	if d.response != nil && d.response.Name != "" {
		return d.response.Name
	}
	return dbDerefString(d.name)
}

// CreatedAt returns the database creation time from the response.
func (d *Database) CreatedAt() time.Time {
	if d.response != nil && d.response.CreationDate != nil {
		return *d.response.CreationDate
	}
	return time.Time{}
}

// CreatedBy returns the identity that created this database.
func (d *Database) CreatedBy() string {
	if d.response != nil && d.response.CreatedBy != nil {
		return *d.response.CreatedBy
	}
	return ""
}

func (d *Database) Raw() *types.DatabaseResponse      { return d.response }
func (d *Database) RawRequest() types.DatabaseRequest { return d.toRequest() }

// Wire conversions.

func (d *Database) toRequest() types.DatabaseRequest {
	return types.DatabaseRequest{Name: dbDerefString(d.name)}
}

func (d *Database) fromResponse(resp *types.DatabaseResponse) {
	if resp == nil {
		return
	}
	d.response = resp
	if resp.Name != "" {
		v := resp.Name
		d.name = &v
	}
}

func dbDerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
