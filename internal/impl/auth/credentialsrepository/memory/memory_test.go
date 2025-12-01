package memory

import "testing"

//go:generate mockgen -package memory -destination=zz_mock_auth_test.go github.com/Arubacloud/sdk-go/internal/ports/auth CredentialsRepository

func TestCredentialsRepository_FetchCredentials(t *testing.T) {

}

func TestCredentialsProxy_FetchCredentials(t *testing.T) {

}
