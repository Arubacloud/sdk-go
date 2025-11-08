package database

import (
	"context"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// DBaaSAPI defines the interface for DBaaS operations
type DBaaSAPI interface {
	ListDBaaS(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.DBaaSList], error)
	GetDBaaS(ctx context.Context, project string, databaseId string, params *schema.RequestParameters) (*schema.Response[schema.DBaaSResponse], error)
	CreateDBaaS(ctx context.Context, project string, body schema.DBaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.DBaaSResponse], error)
	UpdateDBaaS(ctx context.Context, project string, databaseId string, body schema.DBaaSRequest, params *schema.RequestParameters) (*schema.Response[schema.DBaaSResponse], error)
	DeleteDBaaS(ctx context.Context, projectId string, databaseId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// DatabaseAPI defines the interface for Database operations
type DatabaseAPI interface {
	ListDatabases(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*schema.Response[schema.DatabaseList], error)
	GetDatabase(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error)
	CreateDatabase(ctx context.Context, project string, dbaasId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error)
	UpdateDatabase(ctx context.Context, project string, dbaasId string, databaseId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*schema.Response[schema.DatabaseResponse], error)
	DeleteDatabase(ctx context.Context, projectId string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// BackupAPI defines the interface for Backup operations
type BackupAPI interface {
	ListBackups(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.BackupList], error)
	GetBackup(ctx context.Context, project string, backupId string, params *schema.RequestParameters) (*schema.Response[schema.BackupResponse], error)
	CreateBackup(ctx context.Context, project string, body schema.BackupRequest, params *schema.RequestParameters) (*schema.Response[schema.BackupResponse], error)
	DeleteBackup(ctx context.Context, projectId string, backupId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// UserAPI defines the interface for User operations
type UserAPI interface {
	ListUsers(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*schema.Response[schema.UserList], error)
	GetUser(ctx context.Context, project string, dbaasId string, userId string, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error)
	CreateUser(ctx context.Context, project string, dbaasId string, body schema.UserRequest, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error)
	UpdateUser(ctx context.Context, project string, dbaasId string, userId string, body schema.UserRequest, params *schema.RequestParameters) (*schema.Response[schema.UserResponse], error)
	DeleteUser(ctx context.Context, projectId string, dbaasId string, userId string, params *schema.RequestParameters) (*schema.Response[any], error)
}

// GrantAPI defines the interface for Grant operations
type GrantAPI interface {
	ListGrants(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*schema.Response[schema.GrantList], error)
	GetGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error)
	CreateGrant(ctx context.Context, project string, dbaasId string, databaseId string, body schema.GrantRequest, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error)
	UpdateGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, body schema.GrantRequest, params *schema.RequestParameters) (*schema.Response[schema.GrantResponse], error)
	DeleteGrant(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*schema.Response[any], error)
}
