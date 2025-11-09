package security

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// KmsKeyService implements the KMSAPI interface
type KmsKeyService struct {
	client *client.Client
}

// NewKmsKeyService creates a new KmsKeyService
func NewKmsKeyService(client *client.Client) *KmsKeyService {
	return &KmsKeyService{
		client: client,
	}
}

// ListKMSKeys retrieves all KMS keys for a project
func (s *KmsKeyService) ListKMSKeys(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.KmsList], error) {
	s.client.Logger().Debugf("Listing KMS keys for project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KmsKeysPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.KmsList](httpResp)
}

// GetKMSKey retrieves a specific KMS key by ID
func (s *KmsKeyService) GetKMSKey(ctx context.Context, project string, kmsKeyId string, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error) {
	s.client.Logger().Debugf("Getting KMS key: %s in project: %s", kmsKeyId, project)

	if err := schema.ValidateProjectAndResource(project, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KmsKeyPath, project, kmsKeyId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[schema.KmsResponse](httpResp)
}

// CreateKMSKey creates a new KMS key
func (s *KmsKeyService) CreateKMSKey(ctx context.Context, project string, body schema.KmsRequest, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error) {
	s.client.Logger().Debugf("Creating KMS key in project: %s", project)

	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KmsKeysPath, project)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

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
func (s *KmsKeyService) UpdateKMSKey(ctx context.Context, project string, kmsKeyId string, body schema.KmsRequest, params *schema.RequestParameters) (*schema.Response[schema.KmsResponse], error) {
	s.client.Logger().Debugf("Updating KMS key: %s in project: %s", kmsKeyId, project)

	if err := schema.ValidateProjectAndResource(project, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KmsKeyPath, project, kmsKeyId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

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
func (s *KmsKeyService) DeleteKMSKey(ctx context.Context, projectId string, kmsKeyId string, params *schema.RequestParameters) (*schema.Response[any], error) {
	s.client.Logger().Debugf("Deleting KMS key: %s in project: %s", kmsKeyId, projectId)

	if err := schema.ValidateProjectAndResource(projectId, kmsKeyId, "KMS key ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(KmsKeyPath, projectId, kmsKeyId)

	var queryParams map[string]string
	var headers map[string]string

	if params != nil {
		queryParams = params.ToQueryParams()
		headers = params.ToHeaders()
	}

	httpResp, err := s.client.DoRequest(ctx, http.MethodDelete, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return schema.ParseResponseBody[any](httpResp)
}
