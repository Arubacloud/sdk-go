package vault

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

//go:generate mockgen -package vault -destination=zz_mock_vault_test.go github.com/Arubacloud/sdk-go/internal/impl/auth/credentialsrepository/vault VaultClient,LogicalAPI,KvAPI,AuthAPI,AuthTokenAPI

func TestCredentialsRepository_LoginWithAppRole(t *testing.T) {
	t.Run("do nothing if token is already set", func(t *testing.T) {

		repo := &CredentialsRepository{
			tokenExist: true,
		}

		err := repo.loginWithAppRole(t.Context())

		require.NoError(t, err)
		require.True(t, repo.tokenExist)
	})
	t.Run("Error on write secret", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockVaultClient(ctrl)
		mockLogicalAPI := NewMockLogicalAPI(ctrl)

		mockClient.EXPECT().Logical().Return(mockLogicalAPI)

		data := map[string]interface{}{
			"role_id":   "test-role-id",
			"secret_id": "test-secret-id",
		}

		mockLogicalAPI.
			EXPECT().Write("test-role-path", data).
			Return(nil, fmt.Errorf("mock error"))

		repo := &CredentialsRepository{
			client:   mockClient,
			rolePath: "test-role-path",
			roleID:   "test-role-id",
			secretID: "test-secret-id",
		}

		err := repo.loginWithAppRole(t.Context())

		require.Error(t, err)
		require.False(t, repo.tokenExist)
	})
	t.Run("successful login sets token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := NewMockVaultClient(ctrl)
		mockLogicalAPI := NewMockLogicalAPI(ctrl)

		// Setup the mock responses

		calls := []string{}

		mockClient.
			EXPECT().
			SetToken("mock-token").
			Do(func(token string) {
				calls = append(calls, token)
			}).
			AnyTimes().MinTimes(1)

		mockClient.EXPECT().Logical().Return(mockLogicalAPI)

		data := map[string]interface{}{
			"role_id":   "test-role-id",
			"secret_id": "test-secret-id",
		}

		mockLogicalAPI.
			EXPECT().Write("test-role-path", data).
			Return(&vaultapi.Secret{
				Auth: &vaultapi.SecretAuth{
					ClientToken:   "mock-token",
					Renewable:     false,
					LeaseDuration: 3600,
				},
			}, nil)

		repo := &CredentialsRepository{
			client:   mockClient,
			rolePath: "test-role-path",
			roleID:   "test-role-id",
			secretID: "test-secret-id",
		}

		err := repo.loginWithAppRole(t.Context())

		require.NoError(t, err)
		require.True(t, repo.tokenExist)
		require.Equal(t, "mock-token", calls[0])
		require.Equal(t, false, repo.renewable)
		require.Equal(t, 3600*time.Second, repo.ttl)
	})
	t.Run("successful login sets token and namespace", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		mockClient := NewMockVaultClient(ctrl)
		mockLogicalAPI := NewMockLogicalAPI(ctrl)

		// Setup the mock responses

		mockClient.
			EXPECT().
			SetToken("mock-token").
			Return()

		mockClient.EXPECT().
			SetNamespace("test-namespace").
			Return()

		mockClient.EXPECT().Logical().Return(mockLogicalAPI)

		data := map[string]interface{}{
			"role_id":   "test-role-id",
			"secret_id": "test-secret-id",
		}

		mockLogicalAPI.
			EXPECT().Write("test-role-path", data).
			Return(&vaultapi.Secret{
				Auth: &vaultapi.SecretAuth{
					ClientToken:   "mock-token",
					Renewable:     true,
					LeaseDuration: 3600,
				},
			}, nil)

		repo := &CredentialsRepository{
			client:    mockClient,
			rolePath:  "test-role-path",
			roleID:    "test-role-id",
			secretID:  "test-secret-id",
			namespace: "test-namespace",
		}

		var wg sync.WaitGroup
		wg.Add(1)

		var err error
		go func() {
			defer wg.Done()
			err = repo.loginWithAppRole(ctx)
		}()
		wg.Wait()
		cancel()
		require.NoError(t, err)
		require.True(t, repo.tokenExist)
		require.Equal(t, true, repo.renewable)
		require.Equal(t, 3600*time.Second, repo.ttl)
	})
}

func TestCredentialsRepository_GetCredentialsFromVaultSecret(t *testing.T) {
	t.Run("should return credentials when secret contains client_id and client_secret", func(t *testing.T) {
		secret := &vaultapi.KVSecret{
			Data: map[string]interface{}{
				"client_id":     "test-client-id",
				"client_secret": "test-client-secret",
			},
		}

		creds, err := getCredentialsFromVaultSecret(secret)

		require.NoError(t, err)
		require.Equal(t, "test-client-id", creds.ClientID)
		require.Equal(t, "test-client-secret", creds.ClientSecret)
	})

	t.Run("should return error when client_id is missing", func(t *testing.T) {
		secret := &vaultapi.KVSecret{
			Data: map[string]interface{}{
				"client_secret": "test-client-secret",
			},
		}

		creds, err := getCredentialsFromVaultSecret(secret)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrCredentialsNotFound, err)
		require.Nil(t, creds)
	})

	t.Run("should return error when client_secret is missing", func(t *testing.T) {
		secret := &vaultapi.KVSecret{
			Data: map[string]interface{}{
				"client_id": "test-client-id",
			},
		}

		creds, err := getCredentialsFromVaultSecret(secret)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrCredentialsNotFound, err)
		require.Nil(t, creds)
	})
	t.Run("should return error when client_id is not a string", func(t *testing.T) {
		secret := &vaultapi.KVSecret{
			Data: map[string]interface{}{
				"client_id":     12345,
				"client_secret": "test-client-secret",
			},
		}

		creds, err := getCredentialsFromVaultSecret(secret)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrCredentialsNotFound, err)
		require.Nil(t, creds)
	})
	t.Run("should return error when client_secret is not a string", func(t *testing.T) {
		secret := &vaultapi.KVSecret{
			Data: map[string]interface{}{
				"client_id":     "test-client-id",
				"client_secret": 67890,
			},
		}

		creds, err := getCredentialsFromVaultSecret(secret)

		require.Error(t, err)
		require.ErrorIs(t, auth.ErrCredentialsNotFound, err)
		require.Nil(t, creds)
	})
}

func TestCredentialsRepository_FetchCredentials(t *testing.T) {
	t.Run("should report credentials not found error when vault returns nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVaultClient := NewMockVaultClient(ctrl)
		mockKvAPI := NewMockKvAPI(ctrl)

		tokenRepository := &CredentialsRepository{
			client:     mockVaultClient,
			tokenExist: true,
			ttl:        time.Until(time.Now()),
			kvMount:    "test-kv",
			kvPath:     "test-path",
		}

		mockVaultClient.EXPECT().
			KVv2(tokenRepository.kvMount).Return(mockKvAPI)

		mockKvAPI.EXPECT().
			Get(gomock.Any(), tokenRepository.kvPath).
			Return(nil, fmt.Errorf("vault: secret not found"))

		// When we try to fetch the token
		data, err := tokenRepository.FetchCredentials(t.Context())

		// And no token should be returned
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrCredentialsNotFound)
		require.Nil(t, data)
	})

	t.Run("should report credentials not found error when vault secret is missing data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVaultClient := NewMockVaultClient(ctrl)
		mockKvAPI := NewMockKvAPI(ctrl)

		tokenRepository := &CredentialsRepository{
			client:     mockVaultClient,
			tokenExist: true,
			ttl:        time.Until(time.Now()),
			kvMount:    "test-kv",
			kvPath:     "test-path",
		}

		mockVaultClient.EXPECT().
			KVv2(tokenRepository.kvMount).Return(mockKvAPI)

		mockKvAPI.EXPECT().
			Get(gomock.Any(), tokenRepository.kvPath).
			Return(&vaultapi.KVSecret{
				Data: map[string]interface{}{},
			}, nil)

		// When we try to fetch the token
		data, err := tokenRepository.FetchCredentials(t.Context())

		// And no token should be returned
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrCredentialsNotFound)
		require.Nil(t, data)
	})
	t.Run("Run ok when vault secret contains credentials", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockVaultClient := NewMockVaultClient(ctrl)
		mockKvAPI := NewMockKvAPI(ctrl)

		tokenRepository := &CredentialsRepository{
			client:     mockVaultClient,
			tokenExist: true,
			ttl:        time.Until(time.Now()),
			kvMount:    "test-kv",
			kvPath:     "test-path",
		}

		mockVaultClient.EXPECT().
			KVv2(tokenRepository.kvMount).Return(mockKvAPI)

		mockKvAPI.EXPECT().
			Get(gomock.Any(), tokenRepository.kvPath).
			Return(&vaultapi.KVSecret{
				Data: map[string]interface{}{
					"client_id":     "test-client-id",
					"client_secret": "test-client-secret",
				},
			}, nil)

		// When we try to fetch the token
		data, err := tokenRepository.FetchCredentials(t.Context())

		// And no token should be returned
		require.NoError(t, err)
		require.NotNil(t, data)
		require.Equal(t, "test-client-id", data.ClientID)
		require.Equal(t, "test-client-secret", data.ClientSecret)
	})

}

func TestCredentialsRepository_RenewTokenBeforeExpiration(t *testing.T) {
	t.Run("token renewal success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		mockClient := NewMockVaultClient(ctrl)
		mockAuthAPI := NewMockAuthAPI(ctrl)
		mockAuthTokenAPI := NewMockAuthTokenAPI(ctrl)
		repo := &CredentialsRepository{

			client:     mockClient,
			tokenExist: true,
			ttl:        2 * time.Millisecond,
			renewable:  true,
		}

		mockClient.EXPECT().
			Auth().
			Return(mockAuthAPI).
			AnyTimes()

		calls := []string{}
		mockClient.EXPECT().SetToken(gomock.Any()).
			Do(func(token string) {
				calls = append(calls, token)
			}).Return().MinTimes(1)

		mockAuthAPI.EXPECT().Token().Return(mockAuthTokenAPI).AnyTimes()

		mockAuthTokenAPI.EXPECT().
			RenewSelfWithContext(ctx, gomock.Any()).
			Return(&vaultapi.Secret{
				Auth: &vaultapi.SecretAuth{
					ClientToken:   "renewed-token",
					Renewable:     true,
					LeaseDuration: 5,
				},
			}, nil).
			AnyTimes()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			repo.renewTokenBeforeExpiration(ctx)
		}()

		// Wait enough time for the token to be renewed at least once
		time.Sleep(100 * time.Millisecond)
		cancel()
		wg.Wait()

		require.True(t, repo.tokenExist)
		require.Equal(t, 5*time.Second, repo.ttl)
	})
	t.Run("token renewal failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		mockClient := NewMockVaultClient(ctrl)
		mockAuthAPI := NewMockAuthAPI(ctrl)
		mockAuthTokenAPI := NewMockAuthTokenAPI(ctrl)
		repo := &CredentialsRepository{

			client:     mockClient,
			tokenExist: true,
			ttl:        2 * time.Millisecond,
			renewable:  true,
		}

		mockClient.EXPECT().
			Auth().
			Return(mockAuthAPI).
			AnyTimes()

		mockClient.EXPECT().SetToken(gomock.Any()).Return().MaxTimes(0)

		mockAuthAPI.EXPECT().Token().Return(mockAuthTokenAPI).AnyTimes()

		mockAuthTokenAPI.EXPECT().
			RenewSelfWithContext(ctx, gomock.Any()).
			Return(nil, fmt.Errorf("renewal error")).
			Times(1)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			repo.renewTokenBeforeExpiration(ctx)
		}()

		// Wait enough time for the token to be attempted to renew at least once
		time.Sleep(100 * time.Millisecond)

		cancel()
		wg.Wait()

		require.False(t, repo.tokenExist)
	})
}
