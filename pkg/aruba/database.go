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
	List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.BackupList], error)
	Get(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Create(ctx context.Context, projectID string, body types.BackupRequest, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Delete(ctx context.Context, projectID string, backupID string, params *types.RequestParameters) (*types.Response[any], error)
}

type UsersClient interface {
	List(ctx context.Context, projectID string, dbaasID string, params *types.RequestParameters) (*types.Response[types.UserList], error)
	Get(ctx context.Context, projectID string, dbaasID string, userID string, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Create(ctx context.Context, projectID string, dbaasID string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Update(ctx context.Context, projectID string, dbaasID string, userID string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Delete(ctx context.Context, projectID string, dbaasID string, userID string, params *types.RequestParameters) (*types.Response[any], error)
}

type GrantsClient interface {
	List(ctx context.Context, projectID string, dbaasID string, databaseID string, params *types.RequestParameters) (*types.Response[types.GrantList], error)
	Get(ctx context.Context, projectID string, dbaasID string, databaseID string, grantID string, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Create(ctx context.Context, projectID string, dbaasID string, databaseID string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Update(ctx context.Context, projectID string, dbaasID string, databaseID string, grantID string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Delete(ctx context.Context, projectID string, dbaasID string, databaseID string, grantID string, params *types.RequestParameters) (*types.Response[any], error)
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
