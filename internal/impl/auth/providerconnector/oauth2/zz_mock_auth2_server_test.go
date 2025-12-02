package oauth2

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockServerConfig struct {
	StatusCode  int
	AccessToken string
	ExpiresIn   int
	ErrorBody   string
}

func SetupConfigurableTokenServer(t *testing.T, config MockServerConfig) *httptest.Server {
	t.Helper()

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Helper()

		w.WriteHeader(config.StatusCode)

		if config.StatusCode == http.StatusOK {
			w.Header().Set("Content-Type", "application/json")

			response := map[string]interface{}{
				"access_token": config.AccessToken,
				"token_type":   "Bearer",
				"expires_in":   config.ExpiresIn,
			}
			json.NewEncoder(w).Encode(response)

		} else {
			w.Header().Set("Content-Type", "application/json")

			errorResponse := map[string]string{
				"error":             http.StatusText(config.StatusCode),
				"error_description": config.ErrorBody,
			}

			json.NewEncoder(w).Encode(errorResponse)
		}
	}))
}
