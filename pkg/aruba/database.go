package aruba

import (
	"context"
	"fmt"

	"github.com/Arubacloud/sdk-go/internal/clients/database"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

type DatabaseClient interface {
	DBaaS() DBaaSClient
	Databases() DatabasesClient
	Backups() BackupsClient
	Users() UsersClient
	Grants() GrantsClient
}

type databaseClientImpl struct {
	dbaasClient     DBaaSClient
	databasesClient DatabasesClient
	backupsClient   BackupsClient
	usersClient     UsersClient
	grantsClient    GrantsClient
}

var _ DatabaseClient = (*databaseClientImpl)(nil)

func (c databaseClientImpl) DBaaS() DBaaSClient {
	return c.dbaasClient
}

func (c databaseClientImpl) Databases() DatabasesClient {
	return c.databasesClient
}

func (c databaseClientImpl) Backups() BackupsClient {
	return c.backupsClient
}

func (c databaseClientImpl) Users() UsersClient {
	return c.usersClient
}

func (c databaseClientImpl) Grants() GrantsClient {
	return c.grantsClient
}

type DBaaSClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*DBaaS], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*DBaaS, error)
	Create(ctx context.Context, dbaas *DBaaS, opts ...CallOption) (*DBaaS, error)
	Update(ctx context.Context, dbaas *DBaaS, opts ...CallOption) (*DBaaS, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type dbaasLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.DBaaSList], error)
	Get(ctx context.Context, projectID, dbaasID string, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Create(ctx context.Context, projectID string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Update(ctx context.Context, projectID, dbaasID string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Delete(ctx context.Context, projectID, dbaasID string, params *types.RequestParameters) (*types.Response[any], error)
}

type dbaasClientAdapter struct{ low dbaasLowLevelClient }

func newDBaaSClientAdapter(rest *restclient.Client) *dbaasClientAdapter {
	if rest == nil {
		return &dbaasClientAdapter{}
	}
	return &dbaasClientAdapter{low: database.NewDBaaSClientImpl(rest)}
}

func (a *dbaasClientAdapter) Create(ctx context.Context, d *DBaaS, opts ...CallOption) (*DBaaS, error) {
	if err := d.Err(); err != nil {
		return d, err
	}
	if d.ProjectID() == "" {
		return d, fmt.Errorf("Create: DBaaS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, d.ProjectID(), d.toRequest(), rp)
	populateHTTPEnvelope(&d.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		d.fromResponse(resp.Data)
	}
	if err != nil {
		return d, err
	}
	if resp != nil && !resp.IsSuccess() {
		return d, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return d, nil
}

func (a *dbaasClientAdapter) Update(ctx context.Context, d *DBaaS, opts ...CallOption) (*DBaaS, error) {
	if err := d.Err(); err != nil {
		return d, err
	}
	if d.DBaaSID() == "" {
		return d, fmt.Errorf("Update: DBaaS has no ID")
	}
	if d.ProjectID() == "" {
		return d, fmt.Errorf("Update: DBaaS has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, d.ProjectID(), d.DBaaSID(), d.toRequest(), rp)
	populateHTTPEnvelope(&d.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		d.fromResponse(resp.Data)
	}
	if err != nil {
		return d, err
	}
	if resp != nil && !resp.IsSuccess() {
		return d, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return d, nil
}

func (a *dbaasClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*DBaaS, error) {
	projectID, dbaasID, err := dbaasIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, dbaasID, rp)
	out := &DBaaS{}
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

func (a *dbaasClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, dbaasID, err := dbaasIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, dbaasID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *dbaasClientAdapter) List(ctx context.Context, project Ref, opts ...CallOption) (*List[*DBaaS], error) {
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
	var items []*DBaaS
	if resp != nil && resp.Data != nil {
		items = make([]*DBaaS, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			d := &DBaaS{}
			d.fromResponse(&resp.Data.Values[i])
			if d.projectID == "" {
				d.projectID = projectID
			}
			items = append(items, d)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*DBaaS], error) {
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

// dbaasIDsFromRef extracts (projectID, dbaasID) from a Ref.
func dbaasIDsFromRef(ref Ref) (projectID, dbaasID string, err error) {
	did, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDBaaSID); ok {
			return w.DBaaSID(), true
		}
		return "", false
	}, "dbaas")
	if !ok || did == "" {
		return "", "", fmt.Errorf("cannot determine DBaaS ID from Ref %q", ref.URI())
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
	return pid, did, nil
}

type DatabasesClient interface {
	List(ctx context.Context, dbaas Ref, opts ...CallOption) (*List[*Database], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*Database, error)
	Create(ctx context.Context, db *Database, opts ...CallOption) (*Database, error)
	Update(ctx context.Context, db *Database, opts ...CallOption) (*Database, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type BackupsClient interface {
	List(ctx context.Context, project Ref, opts ...CallOption) (*List[*DBaaSBackup], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*DBaaSBackup, error)
	Create(ctx context.Context, b *DBaaSBackup, opts ...CallOption) (*DBaaSBackup, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type UsersClient interface {
	List(ctx context.Context, dbaas Ref, opts ...CallOption) (*List[*User], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*User, error)
	Create(ctx context.Context, u *User, opts ...CallOption) (*User, error)
	Update(ctx context.Context, u *User, opts ...CallOption) (*User, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type GrantsClient interface {
	List(ctx context.Context, database Ref, opts ...CallOption) (*List[*Grant], error)
	Get(ctx context.Context, ref Ref, opts ...CallOption) (*Grant, error)
	Create(ctx context.Context, g *Grant, opts ...CallOption) (*Grant, error)
	Update(ctx context.Context, g *Grant, opts ...CallOption) (*Grant, error)
	Delete(ctx context.Context, ref Ref, opts ...CallOption) error
}

type databasesLowLevelClient interface {
	List(ctx context.Context, projectID, dbaasID string, params *types.RequestParameters) (*types.Response[types.DatabaseList], error)
	Get(ctx context.Context, projectID, dbaasID, databaseID string, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	Create(ctx context.Context, projectID, dbaasID string, body types.DatabaseRequest, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	Update(ctx context.Context, projectID, dbaasID, databaseID string, body types.DatabaseRequest, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	Delete(ctx context.Context, projectID, dbaasID, databaseID string, params *types.RequestParameters) (*types.Response[any], error)
}

type databasesClientAdapter struct{ low databasesLowLevelClient }

func newDatabasesClientAdapter(rest *restclient.Client) *databasesClientAdapter {
	if rest == nil {
		return &databasesClientAdapter{}
	}
	return &databasesClientAdapter{low: database.NewDatabasesClientImpl(rest)}
}

func (a *databasesClientAdapter) Create(ctx context.Context, db *Database, opts ...CallOption) (*Database, error) {
	if err := db.Err(); err != nil {
		return db, err
	}
	if db.ProjectID() == "" {
		return db, fmt.Errorf("Create: Database has no parent project — call IntoDBaaS first")
	}
	if db.DBaaSID() == "" {
		return db, fmt.Errorf("Create: Database has no parent DBaaS — call IntoDBaaS first")
	}
	if db.Name() == "" {
		return db, fmt.Errorf("Create: Database has no name — call WithName first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, db.ProjectID(), db.DBaaSID(), db.toRequest(), rp)
	populateHTTPEnvelope(&db.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		db.fromResponse(resp.Data)
	}
	if err != nil {
		return db, err
	}
	if resp != nil && !resp.IsSuccess() {
		return db, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return db, nil
}

func (a *databasesClientAdapter) Update(ctx context.Context, db *Database, opts ...CallOption) (*Database, error) {
	if err := db.Err(); err != nil {
		return db, err
	}
	if db.DatabaseID() == "" {
		return db, fmt.Errorf("Update: Database has no ID — call WithName first")
	}
	if db.DBaaSID() == "" {
		return db, fmt.Errorf("Update: Database has no parent DBaaS — call IntoDBaaS first")
	}
	if db.ProjectID() == "" {
		return db, fmt.Errorf("Update: Database has no parent project — call IntoDBaaS first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Update(ctx, db.ProjectID(), db.DBaaSID(), db.DatabaseID(), db.toRequest(), rp)
	populateHTTPEnvelope(&db.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		db.fromResponse(resp.Data)
	}
	if err != nil {
		return db, err
	}
	if resp != nil && !resp.IsSuccess() {
		return db, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return db, nil
}

func (a *databasesClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*Database, error) {
	projectID, dbaasID, databaseID, err := databaseIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, dbaasID, databaseID, rp)
	out := &Database{}
	out.dbaasID = dbaasID
	out.projectID = projectID
	name := databaseID
	out.name = &name
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

func (a *databasesClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, dbaasID, databaseID, err := databaseIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, dbaasID, databaseID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *databasesClientAdapter) List(ctx context.Context, dbaas Ref, opts ...CallOption) (*List[*Database], error) {
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
	var items []*Database
	if resp != nil && resp.Data != nil {
		items = make([]*Database, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			db := &Database{}
			db.dbaasID = dbaasID
			db.projectID = projectID
			db.fromResponse(&resp.Data.Values[i])
			items = append(items, db)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*Database], error) {
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

// dbaasBackupIDsFromRef extracts (projectID, backupID) from a Ref. Uses the
// dedicated withDBaaSBackupID interface for typed extraction so a typed
// *StorageBackup Ref does not silently route to the DBaaS endpoint. URI
// fallback (segment "backups") remains inherently ambiguous between the two
// backup scopes — callers must pass URIs from the correct domain.
func dbaasBackupIDsFromRef(ref Ref) (projectID, backupID string, err error) {
	bid, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDBaaSBackupID); ok {
			return w.DBaaSBackupID(), true
		}
		return "", false
	}, "backups")
	if !ok || bid == "" {
		return "", "", fmt.Errorf("cannot determine backup ID from Ref %q", ref.URI())
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

type dbaasBackupsLowLevelClient interface {
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.BackupList], error)
	Get(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Create(ctx context.Context, projectID string, body types.BackupRequest, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Delete(ctx context.Context, projectID, backupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type dbaasBackupsClientAdapter struct{ low dbaasBackupsLowLevelClient }

func newDBaaSBackupsClientAdapter(rest *restclient.Client) *dbaasBackupsClientAdapter {
	if rest == nil {
		return &dbaasBackupsClientAdapter{}
	}
	return &dbaasBackupsClientAdapter{low: database.NewBackupsClientImpl(rest)}
}

func (a *dbaasBackupsClientAdapter) Create(ctx context.Context, b *DBaaSBackup, opts ...CallOption) (*DBaaSBackup, error) {
	if err := b.Err(); err != nil {
		return b, err
	}
	if b.ProjectID() == "" {
		return b, fmt.Errorf("Create: DBaaSBackup has no parent project — call IntoProject first")
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Create(ctx, b.ProjectID(), b.toRequest(), rp)
	populateHTTPEnvelope(&b.httpEnvelopeMixin, resp)
	if resp != nil && resp.Data != nil {
		b.fromResponse(resp.Data)
	}
	if err != nil {
		return b, err
	}
	if resp != nil && !resp.IsSuccess() {
		return b, &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return b, nil
}

func (a *dbaasBackupsClientAdapter) Get(ctx context.Context, ref Ref, opts ...CallOption) (*DBaaSBackup, error) {
	projectID, backupID, err := dbaasBackupIDsFromRef(ref)
	if err != nil {
		return nil, err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Get(ctx, projectID, backupID, rp)
	out := &DBaaSBackup{}
	out.projectID = projectID
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

func (a *dbaasBackupsClientAdapter) Delete(ctx context.Context, ref Ref, opts ...CallOption) error {
	projectID, backupID, err := dbaasBackupIDsFromRef(ref)
	if err != nil {
		return err
	}
	co := applyCallOptions(opts)
	rp := co.toRequestParameters()
	resp, err := a.low.Delete(ctx, projectID, backupID, rp)
	if err != nil {
		return err
	}
	if resp != nil && !resp.IsSuccess() {
		return &HTTPError{StatusCode: resp.StatusCode, Body: resp.RawBody, ErrResp: resp.Error}
	}
	return nil
}

func (a *dbaasBackupsClientAdapter) List(ctx context.Context, parent Ref, opts ...CallOption) (*List[*DBaaSBackup], error) {
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
	var items []*DBaaSBackup
	if resp != nil && resp.Data != nil {
		items = make([]*DBaaSBackup, 0, len(resp.Data.Values))
		for i := range resp.Data.Values {
			b := &DBaaSBackup{}
			b.projectID = projectID
			b.fromResponse(&resp.Data.Values[i])
			if b.projectID == "" {
				b.projectID = projectID
			}
			items = append(items, b)
		}
	}
	refetch := func(_ context.Context, _ string) (*List[*DBaaSBackup], error) {
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

// databaseIDsFromRef extracts (projectID, dbaasID, databaseID) from a Ref.
func databaseIDsFromRef(ref Ref) (projectID, dbaasID, databaseID string, err error) {
	name, ok := extractID(ref, func(r Ref) (string, bool) {
		if w, ok := r.(withDatabaseID); ok {
			return w.DatabaseID(), true
		}
		return "", false
	}, "databases")
	if !ok || name == "" {
		return "", "", "", fmt.Errorf("cannot determine database ID from Ref %q", ref.URI())
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
