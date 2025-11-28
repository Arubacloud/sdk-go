package container

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/impl/interceptor/standard"
	"github.com/Arubacloud/sdk-go/internal/impl/logger/noop"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListContainerRegistry(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Request: %s %s", r.Method, r.URL.Path)

			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/registries" {
				w.WriteHeader(http.StatusOK)
				resp := types.ContainerRegistryList{
					ListResponse: types.ListResponse{Total: 1},
					Values: []types.ContainerRegistryResponse{
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("test-registry"),
							},
							Properties: types.ContainerRegistryPropertiesResult{
								VPC: types.ReferenceResource{
									URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1"),
								},
								Subnet: types.ReferenceResource{
									URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124"),
								},
								SecurityGroup: types.ReferenceResource{
									URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890"),
								},
								PublicIp: types.ReferenceResource{
									URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/elasticips/eip-12345"),
								},
								BlockStorage: types.ReferenceResource{
									URI: *types.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321"),
								},
								BillingPlan: &types.BillingPeriodResource{
									BillingPeriod: *types.StringPtr("Hour"),
								},
								AdminUser: &types.UserCredential{
									Username: "admin",
								},
								ConcurrentUsers: types.IntPtr(100),
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("active"),
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewContainerRegistryClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatalf("resp is nil")
		}
		if resp.Data == nil {
			t.Fatalf("resp.Data is nil")
		}
		if resp.Data.Values == nil {
			t.Fatalf("resp.Data.Values is nil")
		}
		if len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 Container Registry")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "test-registry" {
			t.Errorf("expected name 'test-registry'")
		}
		if resp.Data.Values[0].Properties.PublicIp.URI != "/projects/test-project/providers/Aruba.Network/elasticips/eip-12345" {
			t.Errorf("expected PublicIp URI")
		}
		if resp.Data.Values[0].Properties.VPC.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1" {
			t.Errorf("expected VPC URI")
		}
		if resp.Data.Values[0].Properties.Subnet.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124" {
			t.Errorf("expected Subnet URI")
		}
		if resp.Data.Values[0].Properties.SecurityGroup.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890" {
			t.Errorf("expected SecurityGroup URI")
		}
		if resp.Data.Values[0].Properties.BlockStorage.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321" {
			t.Errorf("expected BlockStorage URI")
		}
		if resp.Data.Values[0].Properties.BillingPlan == nil || resp.Data.Values[0].Properties.BillingPlan.BillingPeriod != "Hour" {
			t.Errorf("expected BillingPlan Hour")
		}
		if resp.Data.Values[0].Properties.AdminUser == nil || resp.Data.Values[0].Properties.AdminUser.Username != "admin" {
			t.Errorf("expected AdminUser username")
		}
		if resp.Data.Values[0].Properties.ConcurrentUsers == nil || *resp.Data.Values[0].Properties.ConcurrentUsers != 100 {
			t.Errorf("expected ConcurrentUsers 100")
		}
	})
}

func TestGetContainerRegistry(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Request: %s %s", r.Method, r.URL.Path)

			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/registries/registry-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.ContainerRegistryResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("test-registry"),
						ID:   types.StringPtr("registry-123"),
					},
					Properties: types.ContainerRegistryPropertiesResult{
						VPC: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1"),
						},
						Subnet: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124"),
						},
						SecurityGroup: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890"),
						},
						PublicIp: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/elasticips/eip-12345"),
						},
						BlockStorage: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321"),
						},
						BillingPlan: &types.BillingPeriodResource{
							BillingPeriod: *types.StringPtr("Hour"),
						},
						AdminUser: &types.UserCredential{
							Username: "admin",
						},
						ConcurrentUsers: types.IntPtr(100),
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("active"),
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewContainerRegistryClientImpl(c)

		resp, err := svc.Get(context.Background(), "test-project", "registry-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "test-registry" {
			t.Errorf("expected name 'test-registry'")
		}
		if resp.Data.Properties.PublicIp.URI != "/projects/test-project/providers/Aruba.Network/elasticips/eip-12345" {
			t.Errorf("expected PublicIp URI")
		}
		if resp.Data.Properties.VPC.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1" {
			t.Errorf("expected VPC URI")
		}
		if resp.Data.Properties.Subnet.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124" {
			t.Errorf("expected Subnet URI")
		}
		if resp.Data.Properties.SecurityGroup.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890" {
			t.Errorf("expected SecurityGroup URI")
		}
		if resp.Data.Properties.BlockStorage.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321" {
			t.Errorf("expected BlockStorage URI")
		}
		if resp.Data.Properties.BillingPlan == nil || resp.Data.Properties.BillingPlan.BillingPeriod != "Hour" {
			t.Errorf("expected BillingPlan Hour")
		}
		if resp.Data.Properties.AdminUser == nil || resp.Data.Properties.AdminUser.Username != "admin" {
			t.Errorf("expected AdminUser username")
		}
		if resp.Data.Properties.ConcurrentUsers == nil || *resp.Data.Properties.ConcurrentUsers != 100 {
			t.Errorf("expected ConcurrentUsers 100")
		}
	})
}

func TestCreateContainerRegistry(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Request: %s %s", r.Method, r.URL.Path)

			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/registries" {
				w.WriteHeader(http.StatusCreated)
				resp := types.ContainerRegistryResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("new-registry"),
						ID:   types.StringPtr("registry-456"),
					},
					Properties: types.ContainerRegistryPropertiesResult{
						VPC: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1"),
						},
						Subnet: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124"),
						},
						SecurityGroup: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890"),
						},
						PublicIp: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/elasticips/eip-12345"),
						},
						BlockStorage: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321"),
						},
						BillingPlan: &types.BillingPeriodResource{
							BillingPeriod: *types.StringPtr("Hour"),
						},
						AdminUser: &types.UserCredential{
							Username: "admin",
						},
						ConcurrentUsers: types.IntPtr(100),
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("creating"),
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewContainerRegistryClientImpl(c)

		body := types.ContainerRegistryRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "new-registry",
				},
			},
			Properties: types.ContainerRegistryPropertiesRequest{
				PublicIp:        types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/elasticips/eip-12345"},
				VPC:             types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1"},
				Subnet:          types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124"},
				SecurityGroup:   types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890"},
				BlockStorage:    types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321"},
				BillingPlan:     &types.BillingPeriodResource{BillingPeriod: "Hour"},
				AdminUser:       &types.UserCredential{Username: "admin"},
				ConcurrentUsers: types.IntPtr(100),
			},
		}

		resp, err := svc.Create(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-registry" {
			t.Errorf("expected name 'new-registry'")
		}
		if resp.Data.Properties.PublicIp.URI != "/projects/test-project/providers/Aruba.Network/elasticips/eip-12345" {
			t.Errorf("expected PublicIp URI")
		}
		if resp.Data.Properties.VPC.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1" {
			t.Errorf("expected VPC URI")
		}
		if resp.Data.Properties.Subnet.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124" {
			t.Errorf("expected Subnet URI")
		}
		if resp.Data.Properties.SecurityGroup.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890" {
			t.Errorf("expected SecurityGroup URI")
		}
		if resp.Data.Properties.BlockStorage.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321" {
			t.Errorf("expected BlockStorage URI")
		}
		if resp.Data.Properties.BillingPlan == nil || resp.Data.Properties.BillingPlan.BillingPeriod != "Hour" {
			t.Errorf("expected BillingPlan Hour")
		}
		if resp.Data.Properties.AdminUser == nil || resp.Data.Properties.AdminUser.Username != "admin" {
			t.Errorf("expected AdminUser username")
		}
		if resp.Data.Properties.ConcurrentUsers == nil || *resp.Data.Properties.ConcurrentUsers != 100 {
			t.Errorf("expected ConcurrentUsers 100")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestUpdateContainerRegistry(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Request: %s %s", r.Method, r.URL.Path)

			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "PUT" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/registries/registry-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.ContainerRegistryResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("updated-registry"),
						ID:   types.StringPtr("registry-123"),
					},
					Properties: types.ContainerRegistryPropertiesResult{
						VPC: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1"),
						},
						Subnet: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124"),
						},
						SecurityGroup: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890"),
						},
						PublicIp: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Network/elasticips/eip-12345"),
						},
						BlockStorage: types.ReferenceResource{
							URI: *types.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321"),
						},
						BillingPlan: &types.BillingPeriodResource{
							BillingPeriod: *types.StringPtr("Hour"),
						},
						AdminUser: &types.UserCredential{
							Username: "admin",
						},
						ConcurrentUsers: types.IntPtr(100),
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("updating"),
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewContainerRegistryClientImpl(c)

		body := types.ContainerRegistryRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "updated-registry",
				},
			},
			Properties: types.ContainerRegistryPropertiesRequest{
				PublicIp:        types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/elasticips/eip-12345"},
				VPC:             types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1"},
				Subnet:          types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124"},
				SecurityGroup:   types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890"},
				BlockStorage:    types.ReferenceResource{URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321"},
				BillingPlan:     &types.BillingPeriodResource{BillingPeriod: "Hour"},
				AdminUser:       &types.UserCredential{Username: "admin"},
				ConcurrentUsers: types.IntPtr(100),
			},
		}

		resp, err := svc.Update(context.Background(), "test-project", "registry-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-registry" {
			t.Errorf("expected name 'updated-registry'")
		}
		if resp.Data.Properties.PublicIp.URI != "/projects/test-project/providers/Aruba.Network/elasticips/eip-12345" {
			t.Errorf("expected PublicIp URI")
		}
		if resp.Data.Properties.VPC.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1" {
			t.Errorf("expected VPC URI")
		}
		if resp.Data.Properties.Subnet.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/subnets/subnet-124" {
			t.Errorf("expected Subnet URI")
		}
		if resp.Data.Properties.SecurityGroup.URI != "/projects/test-project/providers/Aruba.Network/vpcs/vpc-1/securitygroups/sg-67890" {
			t.Errorf("expected SecurityGroup URI")
		}
		if resp.Data.Properties.BlockStorage.URI != "/projects/test-project/providers/Aruba.Storage/blockstorages/bs-54321" {
			t.Errorf("expected BlockStorage URI")
		}
		if resp.Data.Properties.BillingPlan == nil || resp.Data.Properties.BillingPlan.BillingPeriod != "Hour" {
			t.Errorf("expected BillingPlan Hour")
		}
		if resp.Data.Properties.AdminUser == nil || resp.Data.Properties.AdminUser.Username != "admin" {
			t.Errorf("expected AdminUser username")
		}
		if resp.Data.Properties.ConcurrentUsers == nil || *resp.Data.Properties.ConcurrentUsers != 100 {
			t.Errorf("expected ConcurrentUsers 100")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "updating" {
			t.Errorf("expected state 'updating'")
		}
	})
}

func TestDeleteContainerRegistry(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Request: %s %s", r.Method, r.URL.Path)

			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/registries/registry-123" {
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
			Logger:         &noop.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg, cfg.Logger, standard.NewInterceptor())
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewContainerRegistryClientImpl(c)

		_, err = svc.Delete(context.Background(), "test-project", "registry-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
