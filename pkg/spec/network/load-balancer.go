package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

// LoadBalancerService implements the LoadBalancerAPI interface
type LoadBalancerService struct {
	client *client.Client
}

// NewLoadBalancerService creates a new LoadBalancerService
func NewLoadBalancerService(client *client.Client) *LoadBalancerService {
	return &LoadBalancerService{
		client: client,
	}
}

// ListLoadBalancers retrieves all load balancers for a project
func (s *LoadBalancerService) ListLoadBalancers(ctx context.Context, project string, params *schema.RequestParameters) (*schema.Response[schema.LoadBalancerList], error) {
	if err := schema.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(LoadBalancersPath, project)

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

	return schema.ParseResponseBody[schema.LoadBalancerList](httpResp)
}

// GetLoadBalancer retrieves a specific load balancer by ID
func (s *LoadBalancerService) GetLoadBalancer(ctx context.Context, project string, loadBalancerId string, params *schema.RequestParameters) (*schema.Response[schema.LoadBalancerResponse], error) {
	if err := schema.ValidateProjectAndResource(project, loadBalancerId, "load balancer ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(LoadBalancerPath, project, loadBalancerId)

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

	return schema.ParseResponseBody[schema.LoadBalancerResponse](httpResp)
}
