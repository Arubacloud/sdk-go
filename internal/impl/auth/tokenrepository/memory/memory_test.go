package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
)

//go:generate mockgen -package memory -destination=zz_mock_auth_test.go github.com/Arubacloud/sdk-go/internal/ports/auth TokenRepository

// Common parameters
var (
	accessToken = "this is a valid token"
	expiry      = time.Now().Add(24 * time.Hour)
)

func TestTokenRepository_FetchToken(t *testing.T) {
	t.Run("should report a token not found error when it has not a token", func(t *testing.T) {
		// Given a fresh new TokenRepository which contains no token
		tokenRepository := NewTokenRepository()

		// When we try to fetch the token
		token, err := tokenRepository.FetchToken(context.TODO())

		// Then a token not found error should be reported
		require.ErrorIs(t, err, auth.ErrTokenNotFound)

		// And no token should be returned
		require.Nil(t, token)
	})

	t.Run("should return an expired token with no error", func(t *testing.T) {
		// Given a fresh new TokenRepository which contains an expired token
		tokenRepository := NewTokenRepository()

		passedExpiry := time.Now().Add(-24 * time.Hour)

		tokenRepository.token = &auth.Token{
			AccessToken: accessToken,
			Expiry:      passedExpiry,
		}

		//
		// When we try to fetch the token
		token, err := tokenRepository.FetchToken(context.TODO())

		// Then no error shoudt be reported
		require.NoError(t, err)

		// And a token containing the same data should be returned
		require.NotNil(t, token)
		require.Equal(t, accessToken, token.AccessToken)
		require.Equal(t, passedExpiry, token.Expiry)

		// And the token should not be valid
		require.False(t, token.IsValid())

		// And the tokens should not be the same
		require.NotSame(t, tokenRepository.token, token)
	})

	t.Run("should return a valid token", func(t *testing.T) {
		// Given a fresh new TokenRepository which contains an non expired token
		tokenRepository := NewTokenRepository()

		tokenRepository.token = &auth.Token{
			AccessToken: accessToken,
			Expiry:      expiry,
		}

		//
		// When we try to fetch the token
		token, err := tokenRepository.FetchToken(context.TODO())

		// Then no error shoudt be reported
		require.NoError(t, err)

		// And a token containing the same data should be returned
		require.NotNil(t, token)
		require.Equal(t, accessToken, token.AccessToken)
		require.Equal(t, expiry, token.Expiry)

		// And the token should be valid
		require.True(t, token.IsValid())

		// And the tokens should not be the same
		require.NotSame(t, tokenRepository.token, token)
	})
}

func TestTokenRepository_SaveToken(t *testing.T) {
	t.Run("should save a token when the repository does not have a token", func(t *testing.T) {
		// Given a fresh new TokenRepository which contains no token
		tokenRepository := NewTokenRepository()

		// And a valid token
		token := &auth.Token{AccessToken: accessToken, Expiry: expiry}

		// When we try to save the token
		err := tokenRepository.SaveToken(context.TODO(), token)

		// Then no error should be reported
		require.NoError(t, err)

		// And the repository should contain a token with the same data of the given one
		require.NotNil(t, tokenRepository.token)
		require.Equal(t, accessToken, tokenRepository.token.AccessToken)
		require.Equal(t, expiry, tokenRepository.token.Expiry)

		// And the tokens should not be the same
		require.NotSame(t, token, tokenRepository.token)
	})

	t.Run("should replace a token when the repository has an expired token", func(t *testing.T) {
		// Given a fresh new TokenRepository which contains an expired token
		tokenRepository := NewTokenRepository()

		passedExpiry := time.Now().Add(-24 * time.Hour)

		tokenRepository.token = &auth.Token{
			AccessToken: accessToken,
			Expiry:      passedExpiry,
		}

		// And a valid token
		token := &auth.Token{AccessToken: accessToken, Expiry: expiry}

		// When we try to save the token
		err := tokenRepository.SaveToken(context.TODO(), token)

		// Then no error should be reported
		require.NoError(t, err)

		// And the repository should contain a token with the same data of the given one
		require.NotNil(t, tokenRepository.token)
		require.Equal(t, accessToken, tokenRepository.token.AccessToken)
		require.Equal(t, expiry, tokenRepository.token.Expiry)

		// And the tokens should not be the same
		require.NotSame(t, token, tokenRepository.token)
	})

	t.Run("should replace a token when the repository has a valid token", func(t *testing.T) {
		// Given a fresh new TokenRepository which contains a valid token
		tokenRepository := NewTokenRepository()

		differentAcessToken := "different access token"
		differentExpiry := time.Now().Add(1 * time.Hour)

		tokenRepository.token = &auth.Token{
			AccessToken: differentAcessToken,
			Expiry:      differentExpiry,
		}

		// And a valid token
		token := &auth.Token{AccessToken: accessToken, Expiry: expiry}

		// When we try to save the token
		err := tokenRepository.SaveToken(context.TODO(), token)

		// Then no error should be reported
		require.NoError(t, err)

		// And the repository should contain a token with the same data of the given one
		require.NotNil(t, tokenRepository.token)
		require.Equal(t, accessToken, tokenRepository.token.AccessToken)
		require.Equal(t, expiry, tokenRepository.token.Expiry)

		// And the tokens should not be the same
		require.NotSame(t, token, tokenRepository.token)
	})
}

func TestTokenProxy_FetchToken(t *testing.T) {

}

func TestTokenProxy_SaveToken(t *testing.T) {

}

func TestTokenProxyWithRandonExpirationDriftSeconds_FetchToken(t *testing.T) {

}

func TestTokenProxyWithRandonExpirationDriftSeconds_SaveToken(t *testing.T) {

}
