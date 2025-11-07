package compute

import (
	"context"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// CloudServerAPI defines the interface for CloudServer operations
type CloudServerAPI interface {
	ListCloudServers(ctx context.Context, project string, params *schema.ListCloudServerParams, reqEditors ...schema.RequestEditorFn) (*http.Response, error)
	GetCloudServer(ctx context.Context, project string, ncloudServerId string, reqEditors ...schema.RequestEditorFn) (*http.Response, error)
	CreateOrUpdateCloudServer(ctx context.Context, project string, name schema.ResourcePathParam, params *schema.CreateOrUpdateParams, body schema.CloudServerRequest, reqEditors ...schema.RequestEditorFn) (*http.Response, error)
	CreateOrUpdateCloudServerWithBody(ctx context.Context, projectId string, cloudServerId string, params *schema.CreateOrUpdateParams, contentType string, body io.Reader, reqEditors ...schema.RequestEditorFn) (*http.Response, error)
	DeleteCloudServer(ctx context.Context, projectId string, cloudServerId string, reqEditors ...schema.RequestEditorFn) (*http.Response, error)
}
