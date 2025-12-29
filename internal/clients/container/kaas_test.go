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
				resp := types.KaaSList{
					ListResponse: types.ListResponse{Total: 1},
					Values: []types.KaaSResponse{
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("test-kaas"),
							},
							Properties: types.KaaSPropertiesResponse{
								Preset: false,
								HA:     types.BoolPtr(true),
								KubernetesVersion: types.KubernetesVersionInfoResponse{
									Value:       types.StringPtr("1.28.0"),
									Recommended: true,
								},
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

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewKaaSClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 1 {
			t.Errorf("expected 1 KaaS cluster")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "test-kaas" {
			t.Errorf("expected name 'test-kaas'")
		}
		if resp.Data.Values[0].Properties.HA == nil || !*resp.Data.Values[0].Properties.HA {
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
				resp := types.KaaSResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("test-kaas"),
						ID:   types.StringPtr("kaas-123"),
					},
					Properties: types.KaaSPropertiesResponse{
						Preset: false,
						HA:     types.BoolPtr(true),
						KubernetesVersion: types.KubernetesVersionInfoResponse{
							Value:       types.StringPtr("1.28.0"),
							Recommended: true,
						},
						NodePools: func() *[]types.NodePoolPropertiesResponse {
							pools := []types.NodePoolPropertiesResponse{
								{
									Name:        types.StringPtr("default-pool"),
									Nodes:       types.Int32Ptr(3),
									Instance:    &types.InstanceResponse{Name: types.StringPtr("small")},
									DataCenter:  &types.DataCenterResponse{Code: types.StringPtr("dc-01")},
									Autoscaling: false,
								},
							}
							return &pools
						}(),
						ManagementIP: types.StringPtr("10.0.0.100"),
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

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewKaaSClientImpl(c)

		resp, err := svc.Get(context.Background(), "test-project", "kaas-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "test-kaas" {
			t.Errorf("expected name 'test-kaas'")
		}
		if resp.Data.Properties.KubernetesVersion.Value == nil || *resp.Data.Properties.KubernetesVersion.Value != "1.28.0" {
			val := ""
			if resp.Data.Properties.KubernetesVersion.Value != nil {
				val = *resp.Data.Properties.KubernetesVersion.Value
			}
			t.Errorf("expected Kubernetes version '1.28.0', got %s", val)
		}
		if resp.Data.Properties.NodePools == nil || len(*resp.Data.Properties.NodePools) != 1 {
			t.Errorf("expected 1 node pool")
		}
		if resp.Data.Properties.ManagementIP == nil || *resp.Data.Properties.ManagementIP != "10.0.0.100" {
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
				resp := types.KaaSResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("new-kaas"),
						ID:   types.StringPtr("kaas-456"),
					},
					Properties: types.KaaSPropertiesResponse{
						Preset: false,
						HA:     types.BoolPtr(true),
						KubernetesVersion: types.KubernetesVersionInfoResponse{
							Value: types.StringPtr("1.28.0"),
						},
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

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewKaaSClientImpl(c)

		body := types.KaaSRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "new-kaas",
				},
			},
			Properties: types.KaaSPropertiesRequest{
				Preset: types.BoolPtr(false),
				HA:     types.BoolPtr(true),
				KubernetesVersion: types.KubernetesVersionInfo{
					Value: "1.28.0",
				},
			},
		}

		resp, err := svc.Create(context.Background(), "test-project", body, nil)
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
				resp := types.KaaSResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("updated-kaas"),
						ID:   types.StringPtr("kaas-123"),
					},
					Properties: types.KaaSPropertiesResponse{
						Preset: false,
						HA:     types.BoolPtr(true),
						KubernetesVersion: types.KubernetesVersionInfoResponse{
							Value: types.StringPtr("1.29.0"),
						},
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

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewKaaSClientImpl(c)

		body := types.KaaSUpdateRequest{
			Properties: types.KaaSPropertiesUpdateRequest{
				KubernetesVersion: types.KubernetesVersionInfoUpdate{
					Value: "1.29.0",
				},
				NodePools: []types.NodePoolProperties{
					{
						Name:     "default-pool",
						Nodes:    3,
						Instance: "K4A8",
						Zone:     "ITBG-1",
					},
				},
				HA: types.BoolPtr(true),
			},
		}

		resp, err := svc.Update(context.Background(), "test-project", "kaas-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-kaas" {
			t.Errorf("expected name 'updated-kaas'")
		}
		if resp.Data.Properties.KubernetesVersion.Value == nil || *resp.Data.Properties.KubernetesVersion.Value != "1.29.0" {
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

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewKaaSClientImpl(c)

		_, err := svc.Delete(context.Background(), "test-project", "kaas-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
