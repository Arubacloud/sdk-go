package memory

import (
	"testing"

	"github.com/stretchr/testify/require"
)

//go:generate mockgen -package memory -destination=zz_mock_auth_test.go github.com/Arubacloud/sdk-go/internal/ports/auth CredentialsRepository

const (
	clientID     = "client id"
	clientSecret = "client secret"
)

func TestCredentialsRepository_FetchCredentials(t *testing.T) {
	t.Run("should return the credentials", func(t *testing.T) {
		// Given a fresh new credentials repository which contains valid credentials
		credentialsRepository := NewCredentialsRepository(clientID, clientSecret)

		// When we try to fetch the credentials
		credentials, err := credentialsRepository.FetchCredentials(t.Context())

		// Then no error should be reported
		require.NoError(t, err)

		// And some credentials should be returned
		require.NotNil(t, credentials)

		// And the credentials should match with the ones stored on the repository
		require.Equal(t, clientID, credentials.ClientID)
		require.Equal(t, clientSecret, credentials.ClientSecret)

		// And the credential holders should not be the same
		require.NotSame(t, credentialsRepository.credentials, credentials)
	})
}

func TestCredentialsProxy_FetchCredentials(t *testing.T) {

}
