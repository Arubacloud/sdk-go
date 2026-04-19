package project

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListProjects(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/projects" {
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
						},
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		})
		c := testutil.NewClient(t, server.URL)
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

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "project list not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		resp, err := svc.List(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 response")
		}
		if resp.Error == nil {
			t.Fatalf("expected resp.Error to be populated")
		}
		if resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected title 'Not Found', got %v", resp.Error.Title)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		resp, err := svc.List(context.Background(), nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusBadGateway {
			t.Fatalf("expected 502 response")
		}
		if resp.Error != nil {
			t.Errorf("expected resp.Error to be nil for non-JSON body, got %+v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected RawBody 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		_, err := svc.List(context.Background(), nil)
		if err == nil {
			t.Fatal("expected a network error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"total":0,"values":[]}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		if _, err := svc.List(context.Background(), nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestGetProject(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet && r.URL.Path == "/projects/project-123" {
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
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		})
		c := testutil.NewClient(t, server.URL)
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

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "project not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		resp, err := svc.Get(context.Background(), "project-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 response")
		}
		if resp.Error == nil {
			t.Fatalf("expected resp.Error to be populated")
		}
		if resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected title 'Not Found', got %v", resp.Error.Title)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		resp, err := svc.Get(context.Background(), "project-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusBadGateway {
			t.Fatalf("expected 502 response")
		}
		if resp.Error != nil {
			t.Errorf("expected resp.Error to be nil for non-JSON body, got %+v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected RawBody 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		_, err := svc.Get(context.Background(), "project-123", nil)
		if err == nil {
			t.Fatal("expected a network error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
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
			}
			json.NewEncoder(w).Encode(resp)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		if _, err := svc.Get(context.Background(), "project-123", nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty project ID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		_, err := svc.Get(context.Background(), "", nil)
		if err == nil {
			t.Fatal("expected an error for empty project ID")
		}
	})
}

func TestCreateProject(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost && r.URL.Path == "/projects" {
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
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		})
		c := testutil.NewClient(t, server.URL)
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
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "project not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "new-project"},
		}

		resp, err := svc.Create(context.Background(), body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 response")
		}
		if resp.Error == nil {
			t.Fatalf("expected resp.Error to be populated")
		}
		if resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected title 'Not Found', got %v", resp.Error.Title)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "new-project"},
		}

		resp, err := svc.Create(context.Background(), body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusBadGateway {
			t.Fatalf("expected 502 response")
		}
		// TD-010: Create swallows unmarshal failures for non-JSON error bodies; RawBody is the contract here.
		if resp.Error != nil {
			t.Errorf("expected resp.Error to be nil for non-JSON body, got %+v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected RawBody 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "new-project"},
		}

		_, err := svc.Create(context.Background(), body, nil)
		if err == nil {
			t.Fatal("expected a network error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
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
			}
			json.NewEncoder(w).Encode(resp)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "new-project"},
			Properties: types.ProjectPropertiesRequest{
				Description: types.StringPtr("A new project"),
				Default:     false,
			},
		}

		if _, err := svc.Create(context.Background(), body, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestUpdateProject(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut && r.URL.Path == "/projects/project-123" {
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
				}
				json.NewEncoder(w).Encode(resp)
				return
			}
			http.NotFound(w, r)
		})
		c := testutil.NewClient(t, server.URL)
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

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "project not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "updated-project"},
		}

		resp, err := svc.Update(context.Background(), "project-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 response")
		}
		if resp.Error == nil {
			t.Fatalf("expected resp.Error to be populated")
		}
		if resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected title 'Not Found', got %v", resp.Error.Title)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "updated-project"},
		}

		resp, err := svc.Update(context.Background(), "project-123", body, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusBadGateway {
			t.Fatalf("expected 502 response")
		}
		// TD-010: Update swallows unmarshal failures for non-JSON error bodies; RawBody is the contract here.
		if resp.Error != nil {
			t.Errorf("expected resp.Error to be nil for non-JSON body, got %+v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected RawBody 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "updated-project"},
		}

		_, err := svc.Update(context.Background(), "project-123", body, nil)
		if err == nil {
			t.Fatal("expected a network error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
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
			}
			json.NewEncoder(w).Encode(resp)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "updated-project"},
			Properties: types.ProjectPropertiesRequest{
				Description: types.StringPtr("Updated description"),
				Default:     false,
			},
		}

		if _, err := svc.Update(context.Background(), "project-123", body, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty project ID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		body := types.ProjectRequest{
			Metadata: types.ResourceMetadataRequest{Name: "updated-project"},
		}

		_, err := svc.Update(context.Background(), "", body, nil)
		if err == nil {
			t.Fatal("expected an error for empty project ID")
		}
	})
}

func TestDeleteProject(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete && r.URL.Path == "/projects/project-123" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			http.NotFound(w, r)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		_, err := svc.Delete(context.Background(), "project-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "project not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		resp, err := svc.Delete(context.Background(), "project-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected 404 response")
		}
		if resp.Error != nil {
			t.Errorf("expected resp.Error to be nil because Delete does not parse error bodies, got %+v", resp.Error)
		}
		if string(resp.RawBody) != testutil.ErrorBodyJSON("Not Found", "project not found", 404) {
			t.Errorf("expected RawBody to preserve the JSON error payload, got %q", string(resp.RawBody))
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		resp, err := svc.Delete(context.Background(), "project-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.StatusCode != http.StatusBadGateway {
			t.Fatalf("expected 502 response")
		}
		if resp.Error != nil {
			t.Errorf("expected resp.Error to be nil because Delete does not parse error bodies, got %+v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected RawBody 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		_, err := svc.Delete(context.Background(), "project-123", nil)
		if err == nil {
			t.Fatal("expected a network error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewProjectsClientImpl(c)

		if _, err := svc.Delete(context.Background(), "project-123", nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty project ID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewProjectsClientImpl(c)

		_, err := svc.Delete(context.Background(), "", nil)
		if err == nil {
			t.Fatal("expected an error for empty project ID")
		}
	})
}
