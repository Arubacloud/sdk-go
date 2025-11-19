package database

import (
	"context"

	"github.com/Arubacloud/sdk-go/types"
)

// DatabaseAPI defines the unified interface for all Database operations
type DatabaseAPI interface {
	// DBaaS operations
	ListDBaaS(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.DBaaSList], error)
	GetDBaaS(ctx context.Context, project string, databaseId string, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	CreateDBaaS(ctx context.Context, project string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	UpdateDBaaS(ctx context.Context, project string, databaseId string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error)
	DeleteDBaaS(ctx context.Context, projectId string, databaseId string, params *types.RequestParameters) (*types.Response[any], error)

	// Database operations
	ListDatabases(ctx context.Context, project string, dbaasId string, params *types.RequestParameters) (*types.Response[types.DatabaseList], error)
	GetDatabase(ctx context.Context, project string, dbaasId string, databaseId string, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	CreateDatabase(ctx context.Context, project string, dbaasId string, body types.DatabaseRequest, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	UpdateDatabase(ctx context.Context, project string, dbaasId string, databaseId string, body types.DatabaseRequest, params *types.RequestParameters) (*types.Response[types.DatabaseResponse], error)
	DeleteDatabase(ctx context.Context, projectId string, dbaasId string, databaseId string, params *types.RequestParameters) (*types.Response[any], error)

	// Backup operations
	ListBackups(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.BackupList], error)
	GetBackup(ctx context.Context, project string, backupId string, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	CreateBackup(ctx context.Context, project string, body types.BackupRequest, params *types.RequestParameters) (*types.Response[types.BackupResponse], error)
	DeleteBackup(ctx context.Context, projectId string, backupId string, params *types.RequestParameters) (*types.Response[any], error)

	// User operations
	ListUsers(ctx context.Context, project string, dbaasId string, params *types.RequestParameters) (*types.Response[types.UserList], error)
	GetUser(ctx context.Context, project string, dbaasId string, userId string, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	CreateUser(ctx context.Context, project string, dbaasId string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	UpdateUser(ctx context.Context, project string, dbaasId string, userId string, body types.UserRequest, params *types.RequestParameters) (*types.Response[types.UserResponse], error)
	DeleteUser(ctx context.Context, projectId string, dbaasId string, userId string, params *types.RequestParameters) (*types.Response[any], error)

	// Grant operations
	ListGrants(ctx context.Context, project string, dbaasId string, databaseId string, params *types.RequestParameters) (*types.Response[types.GrantList], error)
	GetGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	CreateGrant(ctx context.Context, project string, dbaasId string, databaseId string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	UpdateGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, body types.GrantRequest, params *types.RequestParameters) (*types.Response[types.GrantResponse], error)
	DeleteGrant(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *types.RequestParameters) (*types.Response[any], error)
}
