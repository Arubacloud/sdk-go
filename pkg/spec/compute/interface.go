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

// KeyPairAPI defines the interface for KeyPair operations
type KeyPairAPI interface {
	ListKeyPairs(ctx context.Context, project string, params *schema.RequestParameters) (*http.Response, error)
	GetKeyPair(ctx context.Context, project string, keyPairId string, params *schema.RequestParameters) (*http.Response, error)
	CreateOrUpdateKeyPair(ctx context.Context, project string, body schema.KeyPairRequest, params *schema.RequestParameters) (*http.Response, error)
	DeleteKeyPair(ctx context.Context, projectId string, keyPairId string, params *schema.RequestParameters) (*http.Response, error)
}
