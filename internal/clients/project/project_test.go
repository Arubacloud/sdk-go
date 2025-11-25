package project

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

func TestListProjects(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects" {
				w.WriteHeader(http.StatusOK)
				resp := types.ProjectList{
					ListResponse: types.ListResponse{Total: 2},
					Values: []types.ProjectResponse{
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("default-project"),
								ID:   types.StringPtr("project-123"),
							},
							Properties: types.ProjectPropertiesResponse{
								Description:     types.StringPtr("Default project"),
								Default:         true,
								ResourcesNumber: 10,
							},
							Status: types.ResourceStatus{
								State: types.StringPtr("active"),
							},
						},
						{
							Metadata: types.ResourceMetadataResponse{
								Name: types.StringPtr("test-project"),
								ID:   types.StringPtr("project-456"),
							},
							Properties: types.ProjectPropertiesResponse{
								Description:     types.StringPtr("Test project"),
								Default:         false,
								ResourcesNumber: 5,
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
		svc := NewProjectsClientImpl(c)

		resp, err := svc.List(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil || len(resp.Data.Values) != 2 {
			t.Errorf("expected 2 projects")
		}
		if resp.Data.Values[0].Metadata.Name == nil || *resp.Data.Values[0].Metadata.Name != "default-project" {
			t.Errorf("expected name 'default-project'")
		}
		if !resp.Data.Values[0].Properties.Default {
			t.Errorf("expected first project to be default")
		}
		if resp.Data.Values[0].Properties.ResourcesNumber != 10 {
			t.Errorf("expected 10 resources, got %d", resp.Data.Values[0].Properties.ResourcesNumber)
		}
	})
}

func TestGetProject(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "GET" && r.URL.Path == "/projects/project-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.ProjectResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("my-project"),
						ID:   types.StringPtr("project-123"),
					},
					Properties: types.ProjectPropertiesResponse{
						Description:     types.StringPtr("My test project"),
						Default:         false,
						ResourcesNumber: 15,
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
		svc := NewProjectsClientImpl(c)

		resp, err := svc.Get(context.Background(), "project-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-project" {
			t.Errorf("expected name 'my-project'")
		}
		if resp.Data.Properties.ResourcesNumber != 15 {
			t.Errorf("expected 15 resources, got %d", resp.Data.Properties.ResourcesNumber)
		}
		if resp.Data.Properties.Default {
			t.Errorf("expected project to not be default")
		}
	})
}

func TestCreateProject(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/projects" {
				w.WriteHeader(http.StatusCreated)
				resp := types.ProjectResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("new-project"),
						ID:   types.StringPtr("project-789"),
					},
					Properties: types.ProjectPropertiesResponse{
						Description:     types.StringPtr("A new project"),
						Default:         false,
						ResourcesNumber: 0,
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
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{
				Name: "new-project",
			},
			Properties: types.ProjectPropertiesRequest{
				Description: types.StringPtr("A new project"),
				Default:     false,
			},
		}

		resp, err := svc.Create(context.Background(), body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "new-project" {
			t.Errorf("expected name 'new-project'")
		}
		if resp.Data.Status.State == nil || *resp.Data.Status.State != "creating" {
			t.Errorf("expected state 'creating'")
		}
	})
}

func TestUpdateProject(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "PUT" && r.URL.Path == "/projects/project-123" {
				w.WriteHeader(http.StatusOK)
				resp := types.ProjectResponse{
					Metadata: types.ResourceMetadataResponse{
						Name: types.StringPtr("updated-project"),
						ID:   types.StringPtr("project-123"),
					},
					Properties: types.ProjectPropertiesResponse{
						Description:     types.StringPtr("Updated description"),
						Default:         false,
						ResourcesNumber: 15,
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
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{
				Name: "updated-project",
			},
			Properties: types.ProjectPropertiesRequest{
				Description: types.StringPtr("Updated description"),
				Default:     false,
			},
		}

		resp, err := svc.Update(context.Background(), "project-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "updated-project" {
			t.Errorf("expected name 'updated-project'")
		}
		if resp.Data.Properties.Description == nil || *resp.Data.Properties.Description != "Updated description" {
			t.Errorf("expected description 'Updated description'")
		}
	})
}

func TestDeleteProject(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method == "DELETE" && r.URL.Path == "/projects/project-123" {
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
		svc := NewProjectsClientImpl(c)

		_, err = svc.Delete(context.Background(), "project-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}
