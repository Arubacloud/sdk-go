package security

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/types"
)

type kmsKeysClientImpl struct {
	client *restclient.Client
}

// NewKMSKeysClientImpl creates a new unified Security service
func NewKMSKeysClientImpl(client *restclient.Client) *kmsKeysClientImpl {
	return &kmsKeysClientImpl{
		client: client,
	}
}

// List retrieves all KMS keys for a project
func (c *kmsKeysClientImpl) List(ctx context.Context, projectID string, params *types.RequestParameters) (*types.Response[types.KmsList], error) {
	c.client.Logger().Debugf("Listing KMS keys for project: %s", projectID)

	if err := types.ValidateProject(projectID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeysPath, projectID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &KMSListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.KmsList](httpResp)
}

// Get retrieves a specific KMS key by ID
func (c *kmsKeysClientImpl) Get(ctx context.Context, projectID string, kmsKeyId string, params *types.RequestParameters) (*types.Response[types.KmsResponse], error) {
	c.client.Logger().Debugf("Getting KMS key: %s in project: %s", kmsKeyId, projectID)

	if err := types.ValidateProjectAndResource(projectID, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeyPath, projectID, kmsKeyId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &KMSReadAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSReadAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.KmsResponse](httpResp)
}

// Create creates a new KMS key
func (c *kmsKeysClientImpl) Create(ctx context.Context, projectID string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error) {
	c.client.Logger().Debugf("Creating KMS key in project: %s", projectID)

	if err := types.ValidateProject(projectID); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeysPath, projectID)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &KMSCreateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSCreateAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := c.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &types.Response[types.KmsResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.KmsResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// Update updates an existing KMS key
func (c *kmsKeysClientImpl) Update(ctx context.Context, projectID string, kmsKeyId string, body types.KmsRequest, params *types.RequestParameters) (*types.Response[types.KmsResponse], error) {
	c.client.Logger().Debugf("Updating KMS key: %s in project: %s", kmsKeyId, projectID)

	if err := types.ValidateProjectAndResource(projectID, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeyPath, projectID, kmsKeyId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &KMSUpdateAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSUpdateAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := c.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &types.Response[types.KmsResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data types.KmsResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp types.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// Delete deletes a KMS key by ID
func (c *kmsKeysClientImpl) Delete(ctx context.Context, projectID string, kmsKeyId string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting KMS key: %s in project: %s", kmsKeyId, projectID)

	if err := types.ValidateProjectAndResource(projectID, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeyPath, projectID, kmsKeyId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &KMSDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSDeleteAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[any](httpResp)
}
