package compute

import (
	"context"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// CloudServerAPI defines the interface for CloudServer operations
type CloudServerAPI interface {
	ListCloudServers(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetCloudServer(ctx context.Context, project string, cloudServerId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateCloudServer(ctx context.Context, project string, body schema.CloudServerRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, params *schema.RequestParameters) (*http.Response, error)
}

// Additional interfaces for other compute resources can be defined here
