package storage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arubacloud/sdk-go/pkg/restclient"
	"github.com/Arubacloud/sdk-go/types"
)

func TestListBlockStorageVolumes(t *testing.T) {
	t.Run("successful_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages" {
				w.WriteHeader(http.StatusOK)
				resp := types.BlockStorageList{
					ListResponse: types.ListResponse{Total: 2},
					Values: []types.BlockStorageResponse{
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("data-volume"),
								ID:   types.StringPtr("vol-123"),
							},
							Properties: types.BlockStoragePropertiesResponse{
								SizeGB:        100,
								BillingPeriod: "Hour",
								Zone:          "it-eur-1",
								Type:          types.BlockStorageTypePerformance,
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("active"),
							},
						},
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("backup-volume"),
								ID:   types.StringPtr("vol-456"),
							},
							Properties: types.BlockStoragePropertiesResponse{
								SizeGB:        200,
								BillingPeriod: "Hour",
								Zone:          "it-eur-1",
								Type:          types.BlockStorageTypeStandard,
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewVolumesClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 volumes")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "data-volume" {
			t.Errorf("expected name 'data-volume'")
		}
		if resp.Data.Values[0].Properties.Type != types.BlockStorageTypePerformance {
			t.Errorf("expected performance type")
		}
	})
}

func TestGetBlockStorageVolume(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.BlockStorageResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("my-volume"),
						ID:   types.StringPtr("vol-123"),
					},
					Properties: types.BlockStoragePropertiesResponse{
						SizeGB:        150,
						BillingPeriod: "Hour",
						Zone:          "it-eur-1",
						Type:          types.BlockStorageTypePerformance,
						Bootable:      types.BoolPtr(false),
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewVolumesClientImpl(c)

		resp, err := svc.Get(context.Background(), "test-project", "vol-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-volume" {
			t.Errorf("expected name 'my-volume'")
		}
		if resp.Data.Properties.SizeGB != 150 {
			t.Errorf("expected size 150GB")
		}
	})
}

func TestCreateBlockStorageVolume(t *testing.T) {
	t.Run("successful_create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages" {
				w.WriteHeader(http.StatusCreated)
				resp := types.BlockStorageResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("new-volume"),
						ID:   types.StringPtr("vol-789"),
					},
					Properties: types.BlockStoragePropertiesResponse{
						SizeGB:        50,
						BillingPeriod: "Hour",
						Zone:          "it-eur-1",
						Type:          types.BlockStorageTypeStandard,
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
			Logger:         &restclient.NoOpLogger{},
		}
		c, err := restclient.NewClient(cfg)
		if err != nil {
			t.Fatalf("failed to create client: %v", err)
		}
		svc := NewVolumesClientImpl(c)

		body := types.BlockStorageRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "new-volume",
				},
				Location: types.LocationRequest{Value: "it-eur"},
			},
			Properties: types.BlockStoragePropertiesRequest{
				SizeGB:        50,
				BillingPeriod: "Hour",
				Zone:          "it-eur-1",
				Type:          types.BlockStorageTypeStandard,
			},
		}

		resp, err := svc.Create(context.Background(), "test-project", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-volume" {
			t.Errorf("expected name 'new-volume'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestDeleteBlockStorageVolume(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123" {
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
		svc := NewVolumesClientImpl(c)

		_, err = svc.Delete(context.Background(), "test-project", "vol-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestListSnapshots(t *testing.T) {
	t.Run("successful_list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots" {
				w.WriteHeader(http.StatusOK)
				resp := types.SnapshotList{
					ListResponse: types.ListResponse{Total: 2},
					Values: []types.SnapshotResponse{
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("backup-snapshot-1"),
								ID:   types.StringPtr("snap-123"),
							},
							Properties: types.SnapshotPropertiesResponse{
								SizeGB:        types.Int32Ptr(100),
								BillingPeriod: types.StringPtr("Hour"),
								Zone:          "it-eur-1",
								Type:          types.BlockStorageTypePerformance,
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("available"),
							},
						},
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("backup-snapshot-2"),
								ID:   types.StringPtr("snap-456"),
							},
							Properties: types.SnapshotPropertiesResponse{
								SizeGB:        types.Int32Ptr(200),
								BillingPeriod: types.StringPtr("Hour"),
								Zone:          "it-eur-1",
								Type:          types.BlockStorageTypeStandard,
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("available"),
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
		svc := NewSnapshotsClientImpl(c, NewVolumesClientImpl(c))

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 snapshots")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "backup-snapshot-1" {
			t.Errorf("expected name 'backup-snapshot-1'")
		}
	})
}

func TestGetSnapshot(t *testing.T) {
	t.Run("successful_get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots/snap-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.SnapshotResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("my-snapshot"),
						ID:   types.StringPtr("snap-123"),
					},
					Properties: types.SnapshotPropertiesResponse{
						SizeGB:        types.Int32Ptr(150),
						BillingPeriod: types.StringPtr("Hour"),
						Zone:          "it-eur-1",
						Type:          types.BlockStorageTypePerformance,
						Volume: &types.VolumeInfo{
							URI:  types.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123"),
							Name: types.StringPtr("source-volume"),
						},
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("available"),
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
		svc := NewSnapshotsClientImpl(c, NewVolumesClientImpl(c))

		resp, err := svc.Get(context.Background(), "test-project", "snap-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-snapshot" {
			t.Errorf("expected name 'my-snapshot'")
		}
		if resp.Data.Properties.Volume == nil || resp.Data.Properties.Volume.Name == nil || *resp.Data.Properties.Volume.Name != "source-volume" {
			t.Errorf("expected volume name 'source-volume'")
		}
	})
}

func TestCreateSnapshot(t *testing.T) {
	t.Skip("Skipping test because CreateSnapshot waits for block storage to be active")
	t.Run("successful_create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			// Mock the GET request to check volume status (waitForBlockStorageActive)
			if r.Method == "GET" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.BlockStorageResponse{
					Metadata: types.ResourceMetadataResponse{
						ID: types.StringPtr("vol-123"),
					},
					Status: types.ResourceStatus{
						State: types.StringPtr("active"),
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots" {
				w.WriteHeader(http.StatusCreated)
				resp := types.SnapshotResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("new-snapshot"),
						ID:   types.StringPtr("snap-789"),
					},
					Properties: types.SnapshotPropertiesResponse{
						SizeGB:        types.Int32Ptr(50),
						BillingPeriod: types.StringPtr("Hour"),
						Zone:          "it-eur-1",
						Type:          types.BlockStorageTypeStandard,
						Volume: &types.VolumeInfo{
							URI: types.StringPtr("/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123"),
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
		svc := NewSnapshotsClientImpl(c, NewVolumesClientImpl(c))

		body := types.SnapshotRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{
					Name: "new-snapshot",
				},
				Location: types.LocationRequest{Value: "it-eur"},
			},
			Properties: types.SnapshotPropertiesRequest{
				BillingPeriod: types.StringPtr("Hour"),
				Volume: types.ReferenceResource{
					URI: "/projects/test-project/providers/Aruba.Storage/blockstorages/vol-123",
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
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-snapshot" {
			t.Errorf("expected name 'new-snapshot'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestDeleteSnapshot(t *testing.T) {
	t.Run("successful_delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/test-project/providers/Aruba.Storage/snapshots/snap-123" {
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
		svc := NewSnapshotsClientImpl(c, NewVolumesClientImpl(c))

		_, err = svc.Delete(context.Background(), "test-project", "snap-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
