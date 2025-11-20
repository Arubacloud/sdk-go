package aruba

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
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
	List(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.DBaaSList], error)
	Get(ctx context.Context, project string, databaseId string, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Create(ctx context.Context, project string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Update(ctx context.Context, project string, databaseId string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	Delete(ctx context.Context, projectId string, databaseId string, params *types.RequestParameters) (*types.Response[any], error)
}

type DatabasesClient interface {
	List(ctx context.Context, project string, dbaasId string, params *types.RequestParameters) (*types.Response[types.DatabaseList], error)
	Get(ctx context.Context, project string, dbaasId string, databaseId string, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	Create(ctx context.Context, project string, dbaasId string, body types.DatabaseRequest, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	Update(ctx context.Context, project string, dbaasId string, databaseId string, body types.DatabaseRequest, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	Delete(ctx context.Context, projectId string, dbaasId string, databaseId string, params *types.RequestParameters) (*types.Response[any], error)
}

type BackupsClient interface {
	List(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.BackupList], error)
	Get(ctx context.Context, project string, backupId string, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Create(ctx context.Context, project string, body types.BackupRequest, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	Delete(ctx context.Context, projectId string, backupId string, params *types.RequestParameters) (*types.Response[any], error)
}

type UsersClient interface {
	List(ctx context.Context, project string, dbaasId string, params *types.RequestParameters) (*types.Response[types.UserList], error)
	Get(ctx context.Context, project string, dbaasId string, userId string, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Create(ctx context.Context, project string, dbaasId string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Update(ctx context.Context, project string, dbaasId string, userId string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	Delete(ctx context.Context, projectId string, dbaasId string, userId string, params *types.RequestParameters) (*types.Response[any], error)
}

type GrantsClient interface {
	List(ctx context.Context, project string, dbaasId string, databaseId string, params *types.RequestParameters) (*types.Response[types.GrantList], error)
	Get(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Create(ctx context.Context, project string, dbaasId string, databaseId string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Update(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	Delete(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *types.RequestParameters) (*types.Response[any], error)
}
