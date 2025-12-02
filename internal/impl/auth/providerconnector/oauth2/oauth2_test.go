package oauth2

import "testing"

//go:generate mockgen -package oauth2 -destination=zz_mock_auth_test.go github.com/Arubacloud/sdk-go/internal/ports/auth CredentialsRepository

func TestProviderConnector_RequestToken(t *testing.T) {

}
