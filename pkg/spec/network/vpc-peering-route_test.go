package network

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func TestListVpcPeeringRoutes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1/routes" {
			w.WriteHeader(http.StatusOK)
			resp := schema.VPCPeeringRouteList{
				ListResponse: schema.ListResponse{Total: 1},
				Values: []schema.VPCPeeringRouteResponse{
					{
						Metadata: schema.RegionalResourceMetadataRequest{
							ResourceMetadataRequest: schema.ResourceMetadataRequest{
								Name: "route-1",
							},
						},
						Properties: schema.VPCPeeringRoutePropertiesResponse{
							LocalNetworkAddress:  "10.0.0.0/16",
							RemoteNetworkAddress: "10.1.0.0/16",
						},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &restclient.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	resp, err := svc.ListVpcPeeringRoutes(context.Background(), "test-project", "vpc-123", "peering-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
		t.Errorf("expected 1 vpc peering route")
	}
}

func TestGetVpcPeeringRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1/routes/route-1" {
			w.WriteHeader(http.StatusOK)
			resp := schema.VPCPeeringRouteResponse{
				Metadata: schema.RegionalResourceMetadataRequest{
					ResourceMetadataRequest: schema.ResourceMetadataRequest{
						Name: "route-1",
					},
				},
				Properties: schema.VPCPeeringRoutePropertiesResponse{
					LocalNetworkAddress:  "10.0.0.0/16",
					RemoteNetworkAddress: "10.1.0.0/16",
				},
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &restclient.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	resp, err := svc.GetVpcPeeringRoute(context.Background(), "test-project", "vpc-123", "peering-1", "route-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp == nil || resp.Data == nil || resp.Data.Metadata.Name != "route-1" {
		t.Errorf("expected route name 'route-1'")
	}
}

func TestDeleteVpcPeeringRoute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
			return
		}

		if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123/vpcPeerings/peering-1/routes/route-1" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		http.NotFound(w, r)
	}))
	defer server.Close()

	cfg := &restclient.Config{
		BaseURL:        server.URL,
		HTTPClient:     http.DefaultClient,
		TokenIssuerURL: server.URL + "/token",
		ClientID:       "test-client",
		ClientSecret:   "test-secret",
		Logger:         &restclient.NoOpLogger{},
	}
	c, err := restclient.NewClient(cfg)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	svc := NewService(c)

	_, err = svc.DeleteVpcPeeringRoute(context.Background(), "test-project", "vpc-123", "peering-1", "route-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
