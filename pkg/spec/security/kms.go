package security

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// ListKMSKeys retrieves all KMS keys for a project
func (s *Service) ListKMSKeys(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KmsList], error) {
	s.client.Logger().Debugf("Listing KMS keys for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeysPath, project)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &KMSListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.KmsList](httpResp)
}

// GetKMSKey retrieves a specific KMS key by ID
func (s *Service) GetKMSKey(ctx context.Context, project string, kmsKeyId string, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error) {
	s.client.Logger().Debugf("Getting KMS key: %s in project: %s", kmsKeyId, project)

	if err := schema.ValidateProjectAndResource(project, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeyPath, project, kmsKeyId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &KMSReadAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSReadAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.KmsResponse](httpResp)
}

// CreateKMSKey creates a new KMS key
func (s *Service) CreateKMSKey(ctx context.Context, project string, body schema.KmsRequest, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error) {
	s.client.Logger().Debugf("Creating KMS key in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeysPath, project)

	if params == nil {
		params = &schema.RequestParameters{
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

	httpResp, err := s.client.DoRequest(ctx, http.MethodPost, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.KmsResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.KmsResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// UpdateKMSKey updates an existing KMS key
func (s *Service) UpdateKMSKey(ctx context.Context, project string, kmsKeyId string, body schema.KmsRequest, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error) {
	s.client.Logger().Debugf("Updating KMS key: %s in project: %s", kmsKeyId, project)

	if err := schema.ValidateProjectAndResource(project, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeyPath, project, kmsKeyId)

	if params == nil {
		params = &schema.RequestParameters{
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

	httpResp, err := s.client.DoRequest(ctx, http.MethodPut, path, bytes.NewReader(bodyBytes), queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	response := &schema.Response[schema.KmsResponse]{
		HTTPResponse: httpResp,
		StatusCode:   httpResp.StatusCode,
		Headers:      httpResp.Header,
		RawBody:      respBytes,
	}

	if response.IsSuccess() {
		var data schema.KmsResponse
		if err := json.Unmarshal(respBytes, &data); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		response.Data = &data
	} else if response.IsError() && len(respBytes) > 0 {
		var errorResp schema.ErrorResponse
		if err := json.Unmarshal(respBytes, &errorResp); err == nil {
			response.Error = &errorResp
		}
	}

	return response, nil
}

// DeleteKMSKey deletes a KMS key by ID
func (s *Service) DeleteKMSKey(ctx context.Context, projectId string, kmsKeyId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting KMS key: %s in project: %s", kmsKeyId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KMSKeyPath, projectId, kmsKeyId)

	if params == nil {
		params = &schema.RequestParameters{
			APIVersion: &KMSDeleteAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &KMSDeleteAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[any](httpResp)
}
