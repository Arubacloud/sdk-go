package compute

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

// TestListCloudServers tests the ListCloudServers method
func TestListCloudServers(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		apiCalled := false

		// Create mock server that handles both token and API calls
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			apiCalled = true
			w.WriteHeader(http.StatusOK)
			resp := types.CloudServerList{
				ListResponse: types.ListResponse{Total: 2},
				Values: []types.CloudServerResponse{
					{Metadata: types.RegionalResourceMetadataRequest{
						ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "server-1"},
					}},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		resp, err := svc.List(context.Background(), "test-project", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !apiCalled {
			t.Error("API endpoint was not called")
		}
		if resp.Data.Total != 2 {
			t.Errorf("expected total 2, got %d", resp.Data.Total)
		}
	})

	t.Run("empty project ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		_, err := svc.List(context.Background(), "", nil)
		if err == nil {
			t.Error("expected error for empty project ID")
		}
	})
} // TestGetCloudServer tests the GetCloudServer method
func TestGetCloudServer(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			w.WriteHeader(http.StatusOK)
			resp := types.CloudServerResponse{
				Metadata: types.RegionalResourceMetadataRequest{
					ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "my-server"},
				},
				Properties: types.CloudServerPropertiesResult{Zone: "ITBG-1"},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		resp, err := svc.Get(context.Background(), "test-project", "server-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Metadata.Name != "my-server" {
			t.Errorf("expected name 'my-server', got '%s'", resp.Data.Metadata.Name)
		}
	})

	t.Run("empty project ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		_, err := svc.Get(context.Background(), "", "server-123", nil)
		if err == nil {
			t.Error("expected error for empty project ID")
		}
	})
}

// TestCreateCloudServer tests the CreateCloudServer method
func TestCreateCloudServer(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			w.WriteHeader(http.StatusCreated)
			resp := types.CloudServerResponse{
				Metadata: types.RegionalResourceMetadataRequest{
					ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "new-server"},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		req := types.CloudServerRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "new-server"},
				Location:                types.LocationRequest{Value: "ITBG-Bergamo"},
			},
			Properties: types.CloudServerPropertiesRequest{Zone: "ITBG-1"},
		}

		resp, err := svc.Create(context.Background(), "test-project", req, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}
	})
}

// TestDeleteCloudServer tests the DeleteCloudServer method
func TestDeleteCloudServer(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		_, err := svc.Delete(context.Background(), "test-project", "server-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty project ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		_, err := svc.Delete(context.Background(), "", "server-123", nil)
		if err == nil {
			t.Error("expected error for empty project ID")
		}
	})
}

// TestPowerOnCloudServer tests the PowerOn method
func TestPowerOnCloudServer(t *testing.T) {
	t.Run("successful power on", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/projects/test-project/providers/Aruba.Compute/cloudServers/server-123/poweron" {
				t.Errorf("expected poweron path, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			resp := types.CloudServerResponse{
				Metadata: types.RegionalResourceMetadataRequest{
					ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "my-server"},
				},
				Properties: types.CloudServerPropertiesResult{Zone: "ITBG-1"},
				Status: types.ResourceStatus{
					State: types.StringPtr("active"),
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		resp, err := svc.PowerOn(context.Background(), "test-project", "server-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name != "my-server" {
			t.Errorf("expected name 'my-server', got '%s'", resp.Data.Metadata.Name)
		}
		if !resp.IsSuccess() {
			t.Errorf("expected successful response, got status code %d", resp.StatusCode)
		}
	})
}

// TestPowerOffCloudServer tests the PowerOff method
func TestPowerOffCloudServer(t *testing.T) {
	t.Run("successful power off", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/projects/test-project/providers/Aruba.Compute/cloudServers/server-123/poweroff" {
				t.Errorf("expected poweroff path, got %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
			resp := types.CloudServerResponse{
				Metadata: types.RegionalResourceMetadataRequest{
					ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "my-server"},
				},
				Properties: types.CloudServerPropertiesResult{Zone: "ITBG-1"},
				Status: types.ResourceStatus{
					State: types.StringPtr("stopped"),
				},
			}
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		resp, err := svc.PowerOff(context.Background(), "test-project", "server-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil || resp.Data == nil {
			t.Fatalf("expected response data")
		}
		if resp.Data.Metadata.Name != "my-server" {
			t.Errorf("expected name 'my-server', got '%s'", resp.Data.Metadata.Name)
		}
		if !resp.IsSuccess() {
			t.Errorf("expected successful response, got status code %d", resp.StatusCode)
		}
	})
}

// TestSetPasswordCloudServer tests the SetPassword method
func TestSetPasswordCloudServer(t *testing.T) {
	t.Run("successful set password", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}

			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			if r.URL.Path != "/projects/test-project/providers/Aruba.Compute/cloudServers/server-123/password" {
				t.Errorf("expected password path, got %s", r.URL.Path)
			}

			// Verify request body
			var reqBody types.CloudServerPasswordRequest
			if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
				t.Errorf("failed to decode request body: %v", err)
			}
			if reqBody.Password != "newPassword123" {
				t.Errorf("expected password 'newPassword123', got '%s'", reqBody.Password)
			}

			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		req := types.CloudServerPasswordRequest{
			Password: "newPassword123",
		}

		resp, err := svc.SetPassword(context.Background(), "test-project", "server-123", req, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !resp.IsSuccess() {
			t.Errorf("expected successful response, got status code %d", resp.StatusCode)
		}
	})

	t.Run("empty project ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`))
				return
			}
		}))
		defer server.Close()

		var (
			baseURL    = server.URL
			httpClient = http.DefaultClient
			logger     = &noop.NoOpLogger{}
		)

		c := restclient.NewClient(baseURL, httpClient, standard.NewInterceptor(), logger)

		svc := NewCloudServersClientImpl(c)

		req := types.CloudServerPasswordRequest{
			Password: "newPassword123",
		}

		_, err := svc.SetPassword(context.Background(), "", "server-123", req, nil)
		if err == nil {
			t.Error("expected error for empty project ID")
		}
	})
}
