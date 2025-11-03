package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMockTokenServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := TokenResponse{
			AccessToken: "mock-access-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestNewClient(t *testing.T) {
	tokenServer := setupMockTokenServer(t)
	defer tokenServer.Close()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				BaseURL:        "https://api.example.com",
				TokenIssuerURL: tokenServer.URL,
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				HTTPClient:     http.DefaultClient,
			},
			wantErr: false,
		},
		{
			name:    "nil config uses defaults but fails validation",
			config:  nil,
			wantErr: true,
		},
		{
			name: "missing base URL",
			config: &Config{
				TokenIssuerURL: tokenServer.URL,
				ClientID:       "test-client-id",
				ClientSecret:   "test-client-secret",
				HTTPClient:     http.DefaultClient,
			},
			wantErr: true,
		},
		{
			name: "missing token issuer URL",
			config: &Config{
				BaseURL:      "https://api.example.com",
				ClientID:     "test-client-id",
				ClientSecret: "test-client-secret",
				HTTPClient:   http.DefaultClient,
			},
			wantErr: true,
		},
		{
			name: "missing client ID",
			config: &Config{
				BaseURL:        "https://api.example.com",
				TokenIssuerURL: tokenServer.URL,
				ClientSecret:   "test-client-secret",
				HTTPClient:     http.DefaultClient,
			},
			wantErr: true,
		},
		{
			name: "missing client secret",
			config: &Config{
				BaseURL:        "https://api.example.com",
				TokenIssuerURL: tokenServer.URL,
				ClientID:       "test-client-id",
				HTTPClient:     http.DefaultClient,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClient_WithContext(t *testing.T) {
	tokenServer := setupMockTokenServer(t)
	defer tokenServer.Close()

	config := &Config{
		BaseURL:        "https://api.example.com",
		TokenIssuerURL: tokenServer.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		HTTPClient:     http.DefaultClient,
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	ctx := context.WithValue(context.Background(), "key", "value")
	newClient := client.WithContext(ctx)

	if newClient.Context() != ctx {
		t.Error("WithContext() did not set context correctly")
	}

	// Original client should not be modified
	if client.Context() == ctx {
		t.Error("WithContext() modified original client")
	}
}

func TestClient_RequestEditorFn(t *testing.T) {
	tokenServer := setupMockTokenServer(t)
	defer tokenServer.Close()

	config := &Config{
		BaseURL:        "https://api.example.com",
		TokenIssuerURL: tokenServer.URL,
		ClientID:       "test-client-id",
		ClientSecret:   "test-client-secret",
		HTTPClient:     http.DefaultClient,
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
	}

	client, err := NewClient(config)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	// Create a test request
	req, err := http.NewRequest(http.MethodGet, "https://api.example.com/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Apply the request editor
	editorFn := client.RequestEditorFn()
	if err := editorFn(context.Background(), req); err != nil {
		t.Fatalf("RequestEditorFn() error = %v", err)
	}

	// Check Authorization header
	authHeader := req.Header.Get("Authorization")
	if authHeader != "Bearer mock-access-token" {
		t.Errorf("Expected Authorization header 'Bearer mock-access-token', got '%s'", authHeader)
	}

	// Check custom header
	customHeader := req.Header.Get("X-Custom-Header")
	if customHeader != "custom-value" {
		t.Errorf("Expected X-Custom-Header 'custom-value', got '%s'", customHeader)
	}
}
