package network

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Arubacloud/sdk-go/types"
)

// ListLoadBalancers retrieves all load balancers for a project
func (s *Service) ListLoadBalancers(ctx context.Context, project string, params *types.RequestParameters) (*types.Response[types.LoadBalancerList], error) {
	if err := types.ValidateProject(project); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(LoadBalancersPath, project)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &LoadBalancerListAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &LoadBalancerListAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.LoadBalancerList](httpResp)
}

// GetLoadBalancer retrieves a specific load balancer by ID
func (s *Service) GetLoadBalancer(ctx context.Context, project string, loadBalancerId string, params *types.RequestParameters) (*types.Response[types.LoadBalancerResponse], error) {
	if err := types.ValidateProjectAndResource(project, loadBalancerId, "load balancer ID"); err != nil {
		return nil, err
	}

	path := fmt.Sprintf(LoadBalancerPath, project, loadBalancerId)

	if params == nil {
		params = &types.RequestParameters{
			APIVersion: &LoadBalancerGetAPIVersion,
		}
	} else if params.APIVersion == nil {
		params.APIVersion = &LoadBalancerGetAPIVersion
	}

	queryParams := params.ToQueryParams()
	headers := params.ToHeaders()

	httpResp, err := s.client.DoRequest(ctx, http.MethodGet, path, nil, queryParams, headers)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	return types.ParseResponseBody[types.LoadBalancerResponse](httpResp)
}
