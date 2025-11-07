package database

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// DBaaSAPI defines the interface for DBaaS operations
type DBaaSAPI interface {
	ListDBaaS(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetDBaaS(ctx context.Context, project string, databaseId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateDBaaS(ctx context.Context, project string, body schema.DBaaSRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteDBaaS(ctx context.Context, projectId string, databaseId string, params *schema.RequestParameters) (*http.Response, error)
}

// DatabaseAPI defines the interface for Database operations
type DatabaseAPI interface {
	ListDatabases(ctx context.Context, project string, dbaasId string, params *schema.RequestParameters) (*http.Response, error)
	GetDatabase(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateDatabase(ctx context.Context, project string, dbaasId string, body schema.DatabaseRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteDatabase(ctx context.Context, projectId string, dbaasId string, databaseId string, params *schema.RequestParameters) (*http.Response, error)
}

// BackupAPI defines the interface for Backup operations
type BackupAPI interface {
	ListBackups(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetBackup(ctx context.Context, project string, backupId string, params *schema.RequestParameters) (*http.Response, error)
	CreateBackup(ctx context.Context, project string, body schema.BackupRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteBackup(ctx context.Context, projectId string, backupId string, params *schema.RequestParameters) (*http.Response, error)
}

// UserAPI defines the interface for User operations
type UserAPI interface {
	ListUsers(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetUser(ctx context.Context, project string, userId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateUser(ctx context.Context, project string, body schema.UserRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteUser(ctx context.Context, projectId string, userId string, params *schema.RequestParameters) (*http.Response, error)
}

// GrantAPI defines the interface for Grant operations
type GrantAPI interface {
	ListGrants(ctx context.Context, project string, dbaasId string, databaseId string, params *schema.RequestParameters) (*http.Response, error)
	GetGrant(ctx context.Context, project string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateGrant(ctx context.Context, project string, dbaasId string, databaseId string, body schema.GrantRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteGrant(ctx context.Context, projectId string, dbaasId string, databaseId string, grantId string, params *schema.RequestParameters) (*http.Response, error)
}
