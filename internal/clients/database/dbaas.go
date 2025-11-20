package database

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

type dbaasClientImpl struct {
	client *restclient.Client
}

// NewService creates a new unified Database service
func NewDBaaSClientImpl(client *restclient.Client) *dbaasClientImpl {
	return &dbaasClientImpl{
		client: client,
	}
}

// List retrieves all DBaaS instances for a project
func (c *dbaasClientImpl) List(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.DBaaSList], error) {
	c.client.Logger().Debugf("Listing DBaaS instances for project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSPath, project)
	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseDBaaSListVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSListVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.DBaaSList](httpResp)
}

// Get retrieves a specific DBaaS instance by ID
func (c *dbaasClientImpl) Get(ctx context.Context, project string, dbaasId string, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error) {
	c.client.Logger().Debugf("Getting DBaaS instance: %s in project: %s", dbaasId, project)

	if err := types.ValidateProjectAndResource(project, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSItemPath, project, dbaasId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseDBaaSGetVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSGetVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := c.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.DBaaSResponse](httpResp)
}

// Create creates a new DBaaS instance
func (c *dbaasClientImpl) Create(ctx context.Context, project string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error) {
	c.client.Logger().Debugf("Creating DBaaS instance in project: %s", project)

	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseDBaaSCreateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSCreateVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := c.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read the response body
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &types.Response[types.DBaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.DBaaSResponse
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

// Update updates an existing DBaaS instance
func (c *dbaasClientImpl) Update(ctx context.Context, project string, databaseId string, body types.DBaaSRequest, params *types.RequestParameters) (*types.Response[types.DBaaSResponse], error) {
	c.client.Logger().Debugf("Updating DBaaS instance: %s in project: %s", databaseId, project)

	if err := types.ValidateProjectAndResource(project, databaseId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSItemPath, project, databaseId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseDBaaSUpdateVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSUpdateVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	// Marshal the request body to JSON
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	httpResp, err := c.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	// Read the response body
	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Create the response wrapper
	response := &types.Response[types.DBaaSResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	// Parse the response body if successful
	if response.IsSuccess() {
		var data types.DBaaSResponse
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

// Delete deletes a DBaaS instance by ID
func (c *dbaasClientImpl) Delete(ctx context.Context, projectId string, dbaasId string, params *types.RequestParameters) (*types.Response[any], error) {
	c.client.Logger().Debugf("Deleting DBaaS instance: %s in project: %s", dbaasId, projectId)

	if err := types.ValidateProjectAndResource(projectId, dbaasId, "DBaaS ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(DBaaSItemPath, projectId, dbaasId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &DatabaseDBaaSDeleteVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &DatabaseDBaaSDeleteVersion
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
