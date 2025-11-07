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
