package container

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/pkg/client"
	"github.com/Arubacloud/sdk-go/pkg/spec/schema"
)

func TestListKaaS(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/kaas" {
				w.WriteHeader(http.StatusOK)
				resp := schema.KaaSList{
					ListResponse: schema.ListResponse{Total: 1},
					Values: []schema.KaaSResponse{
						{
							Metadata: schema.ResourceMetadataResponse{
								Name: schema.StringPtr("test-kaas"),
							},
							Properties: schema.KaaSPropertiesResponse{
								Preset: false,
								Ha:     true,
								KubernetesVersion: schema.KubernetesVersionInfoResponse{
									KubernetesVersionInfo: schema.KubernetesVersionInfo{
										Value: "1.28.0",
									},
									Recommended: true,
								},
							},
							Status: schema.ResourceStatus{
								State: schema.StringPtr("active"),
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

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		resp, err := svc.ListKaaS(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 KaaS cluster")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "test-kaas" {
			t.Errorf("expected name 'test-kaas'")
		}
		if !resp.Data.Values[0].Properties.Ha {
			t.Errorf("expected HA to be true")
		}
	})
}

func TestGetKaaS(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/kaas/kaas-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.KaaSResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("test-kaas"),
						Id:   schema.StringPtr("kaas-123"),
					},
					Properties: schema.KaaSPropertiesResponse{
						Preset: false,
						Ha:     true,
						KubernetesVersion: schema.KubernetesVersionInfoResponse{
							KubernetesVersionInfo: schema.KubernetesVersionInfo{
								Value: "1.28.0",
							},
							Recommended: true,
						},
						NodePools: []schema.NodePoolPropertiesResponse{
							{
								NodePoolProperties: schema.NodePoolProperties{
									Name:     "default-pool",
									Nodes:    3,
									Instance: "small",
									Zone:     "dc-01",
								},
								Autoscaling: false,
							},
						},
						ManagementIp: schema.StringPtr("10.0.0.100"),
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("active"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		resp, err := svc.GetKaaS(context.Background(), "test-project", "kaas-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "test-kaas" {
			t.Errorf("expected name 'test-kaas'")
		}
		if resp.Data.Properties.KubernetesVersion.Value != "1.28.0" {
			t.Errorf("expected Kubernetes version '1.28.0', got %s", resp.Data.Properties.KubernetesVersion.Value)
		}
		if len(resp.Data.Properties.NodePools) != 1 {
			t.Errorf("expected 1 node pool")
		}
		if resp.Data.Properties.ManagementIp == nil || *resp.Data.Properties.ManagementIp != "10.0.0.100" {
			t.Errorf("expected management IP '10.0.0.100'")
		}
	})
}

func TestCreateKaaS(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/kaas" {
				w.WriteHeader(http.StatusCreated)
				resp := schema.KaaSResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("new-kaas"),
						Id:   schema.StringPtr("kaas-456"),
					},
					Properties: schema.KaaSPropertiesResponse{
						Preset: false,
						Ha:     true,
						KubernetesVersion: schema.KubernetesVersionInfoResponse{
							KubernetesVersionInfo: schema.KubernetesVersionInfo{
								Value: "1.28.0",
							},
						},
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("creating"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		body := schema.KaaSRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "new-kaas",
				},
			},
			Properties: schema.KaaSPropertiesRequest{
				Preset: false,
				Ha:     true,
				KubernetesVersion: schema.KubernetesVersionInfo{
					Value: "1.28.0",
				},
			},
		}

		resp, err := svc.CreateKaaS(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-kaas" {
			t.Errorf("expected name 'new-kaas'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestUpdateKaaS(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "PUT" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/kaas/kaas-123" {
				w.WriteHeader(http.StatusOK)
				resp := schema.KaaSResponse{
					Metadata: schema.ResourceMetadataResponse{
						Name: schema.StringPtr("updated-kaas"),
						Id:   schema.StringPtr("kaas-123"),
					},
					Properties: schema.KaaSPropertiesResponse{
						Preset: false,
						Ha:     true,
						KubernetesVersion: schema.KubernetesVersionInfoResponse{
							KubernetesVersionInfo: schema.KubernetesVersionInfo{
								Value: "1.29.0",
							},
						},
					},
					Status: schema.ResourceStatus{
						State: schema.StringPtr("updating"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		body := schema.KaaSRequest{
			Metadata: schema.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: schema.ResourceMetadataRequest{
					Name: "updated-kaas",
				},
			},
			Properties: schema.KaaSPropertiesRequest{
				Preset: false,
				Ha:     true,
				KubernetesVersion: schema.KubernetesVersionInfo{
					Value: "1.29.0",
				},
			},
		}

		resp, err := svc.UpdateKaaS(context.Background(), "test-project", "kaas-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-kaas" {
			t.Errorf("expected name 'updated-kaas'")
		}
		if resp.Data.Properties.KubernetesVersion.Value != "1.29.0" {
			t.Errorf("expected Kubernetes version '1.29.0'")
		}
	})
}

func TestDeleteKaaS(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Container/kaas/kaas-123" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			http.NotFound(w, r)
		}))
		defer server.Close()

		cfg := &client.Config{
			BaseURL:        server.URL,
			HTTPClient:     http.DefaultClient,
			TokenIssuerURL: server.URL + "/token",
			ClientID:       "test-client",
			ClientSecret:   "test-secret",
			Logger:         &client.NoOpLogger{},
		}
		c, err := client.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewService(c)

		_, err = svc.DeleteKaaS(context.Background(), "test-project", "kaas-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
