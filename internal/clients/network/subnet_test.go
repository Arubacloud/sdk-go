package network

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Arubacloud/sdk-go/internal/testutil"
	"github.com/Arubacloud/sdk-go/pkg/types"
)

func TestListSubnets(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"total":1,"values":[{"metadata":{"name":"subnet-1"}}]}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.List(context.Background(), "test-project", "vpc-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Total != 1 {
			t.Errorf("expected total 1, got %d", resp.Data.Total)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.List(context.Background(), "", "vpc-123", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpcID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.List(context.Background(), "test-project", "", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.List(context.Background(), "test-project", "vpc-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.List(context.Background(), "test-project", "vpc-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.List(context.Background(), "test-project", "vpc-123", nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"total":0,"values":[]}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.List(context.Background(), "test-project", "vpc-123", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestGetSubnet(t *testing.T) {
	t.Run("successful get", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"metadata":{"name":"my-subnet"}}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Get(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Data.Metadata.Name == nil || *resp.Data.Metadata.Name != "my-subnet" {
			t.Errorf("expected name 'my-subnet', got %v", resp.Data.Metadata.Name)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Get(context.Background(), "", "vpc-123", "subnet-456", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpcID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Get(context.Background(), "test-project", "", "subnet-456", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty subnetID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Get(context.Background(), "test-project", "vpc-123", "", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Get(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Get(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Get(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Get(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestCreateSubnet(t *testing.T) {
	// TODO(TD-020): unskip once VPC-active polling is mockable.
	t.Skip("Skipping CreateSubnet test - requires complex VPC polling mock setup")
	// NOTE: CreateSubnet calls waitForVPCActive() which polls the VPC status
	// To properly test this, you need to mock the VPC GET endpoint to return "active" status
	// Example path: /projects/test-project/providers/Aruba.Network/vpcs/vpc-123
	t.Run("successful create", func(t *testing.T) {
		requestCount := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCount++
			t.Logf("Request #%d: %s %s", requestCount, r.Method, r.URL.Path)

			// Limit requests to prevent infinite loops during testing
			if requestCount > 50 {
				t.Error("Too many requests - infinite loop detected")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Handle token endpoint
			if r.URL.Path == "/token" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				tokenResp := `{"access_token":"test-token","token_type":"Bearer","expires_in":3600}`
				t.Logf("Returning token response: %s", tokenResp)
				w.Write([]byte(tokenResp))
				return
			}

			// Handle VPC status polling - GET request to VPC endpoint
			// Path: /projects/test-project/providers/Aruba.Network/vpcs/vpc-123
			if r.Method == http.MethodGet && r.URL.Path == "/projects/test-project/providers/Aruba.Network/vpcs/vpc-123" {
				t.Logf("Returning active VPC status")
				w.WriteHeader(http.StatusOK)
				vpcResp := types.VPCResponse{
					Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("test-vpc")},
					Status:   types.ResourceStatus{State: types.StringPtr("active")},
				}
				json.NewEncoder(w).Encode(vpcResp)
				return
			}

			// Handle subnet creation - POST request
			if r.Method == http.MethodPost {
				t.Logf("Creating subnet")
				w.WriteHeader(http.StatusCreated)
				resp := types.SubnetResponse{
					Metadata: types.ResourceMetadataResponse{Name: types.StringPtr("new-subnet")},
					Properties: types.SubnetPropertiesResponse{
						Type: types.SubnetTypeAdvanced,
						Network: &types.SubnetNetwork{
							Address: "192.168.1.0/25",
						},
						DHCP: &types.SubnetDHCP{
							Enabled: true,
							Range: &types.SubnetDHCPRange{
								Start: "192.168.1.10",
								Count: 50,
							},
							Routes: []types.SubnetDHCPRoute{
								{
									Address: "0.0.0.0/0",
									Gateway: "192.168.1.1",
								},
							},
							DNS: []string{"8.8.8.8", "8.8.4.4"},
						},
					},
				}
				json.NewEncoder(w).Encode(resp)
				return
			}

			// If we get here, something unexpected happened
			t.Logf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		c := testutil.NewClient(t, server.URL)

		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))

		req := types.SubnetRequest{
			Metadata: types.RegionalResourceMetadataRequest{
				ResourceMetadataRequest: types.ResourceMetadataRequest{Name: "new-subnet"},
			},
		}

		resp, err := svc.Create(context.Background(), "test-project", "vpc-123", req, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Errorf("expected status 201, got %d", resp.StatusCode)
		}
	})
}

func TestUpdateSubnet(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", r.Method)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{"metadata":{"name":"updated-subnet"}}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Update(context.Background(), "test-project", "vpc-123", "subnet-456", types.SubnetRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Update(context.Background(), "", "vpc-123", "subnet-456", types.SubnetRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpcID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Update(context.Background(), "test-project", "", "subnet-456", types.SubnetRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty subnetID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Update(context.Background(), "test-project", "vpc-123", "", types.SubnetRequest{}, nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Update(context.Background(), "test-project", "vpc-123", "subnet-456", types.SubnetRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		// TODO(TD-010): Create/Update's manual response build silently swallows non-JSON
		// unmarshal errors (diverges from ParseResponseBody which logs at DEBUG).
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Update(context.Background(), "test-project", "vpc-123", "subnet-456", types.SubnetRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Update(context.Background(), "test-project", "vpc-123", "subnet-456", types.SubnetRequest{}, nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
		}
	})

	t.Run("nil params injects default api-version", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.URL.Query().Get("api-version"); got != "1.0" {
				t.Errorf("expected api-version=1.0, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, `{}`)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Update(context.Background(), "test-project", "vpc-123", "subnet-456", types.SubnetRequest{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}
	})
}

func TestDeleteSubnet(t *testing.T) {
	t.Run("successful delete", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodDelete {
				t.Errorf("expected DELETE, got %s", r.Method)
			}
			w.WriteHeader(http.StatusNoContent)
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Delete(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("empty project", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Delete(context.Background(), "", "vpc-123", "subnet-456", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty vpcID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Delete(context.Background(), "test-project", "", "subnet-456", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("empty subnetID", func(t *testing.T) {
		c := testutil.NewClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Delete(context.Background(), "test-project", "vpc-123", "", nil)
		if err == nil {
			t.Fatal("expected validation error, got nil")
		}
	})

	t.Run("not found", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, testutil.ErrorBodyJSON("Not Found", "resource not found", 404))
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Delete(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", resp.StatusCode)
		}
		if resp.Error == nil || resp.Error.Title == nil || *resp.Error.Title != "Not Found" {
			t.Errorf("expected error title 'Not Found', got %v", resp.Error)
		}
	})

	t.Run("bad gateway non-json", func(t *testing.T) {
		server := testutil.NewMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			fmt.Fprint(w, "Bad Gateway")
		})
		c := testutil.NewClient(t, server.URL)
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		resp, err := svc.Delete(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.StatusCode != http.StatusBadGateway {
			t.Errorf("expected status 502, got %d", resp.StatusCode)
		}
		if resp.Error != nil {
			t.Errorf("expected nil Error, got %v", resp.Error)
		}
		if string(resp.RawBody) != "Bad Gateway" {
			t.Errorf("expected raw body 'Bad Gateway', got %q", string(resp.RawBody))
		}
	})

	t.Run("network error", func(t *testing.T) {
		c := testutil.NewBrokenClient(t, "http://unused.invalid")
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Delete(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err == nil {
			t.Fatal("expected transport error, got nil")
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
		svc := NewSubnetsClientImpl(c, NewVPCsClientImpl(c))
		_, err := svc.Delete(context.Background(), "test-project", "vpc-123", "subnet-456", nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestNewSubnetsClientImpl_panicsOnNilVPCClient(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic on nil vpcClient but got none")
		}
		if !strings.Contains(fmt.Sprint(r), "vpcClient") {
			t.Fatalf("expected panic message to mention vpcClient, got: %v", r)
		}
	}()
	NewSubnetsClientImpl(nil, nil)
}
