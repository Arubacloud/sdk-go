package memory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"

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
	t.Run("should forward errors from the persistent repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Given a persistent repository which still not has a token
		persistentRepository := NewMockTokenRepository(ctrl)

		errConnection := errors.New("connection error")

		persistentRepository.EXPECT().FetchToken(
			gomock.AssignableToTypeOf(context.TODO()),
		).Return(nil, errConnection).Times(1)

		//
		// And a fresh new proxy using that last
		proxy := NewTokenProxy(persistentRepository)

		// When we try to fetch the token from the proxy
		token, err := proxy.FetchToken(context.TODO())

		// Then the same error obtained from the persistent repository should be reported
		require.ErrorIs(t, err, errConnection)

		// And no token should be returned
		require.Nil(t, token)
	})

	t.Run("should fetch the token from the persistent repository when it has no token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Given a persistent repository which already has a token
		persistentRepository := NewMockTokenRepository(ctrl)

		persistentRepository.EXPECT().FetchToken(gomock.AssignableToTypeOf(context.TODO())).Return(
			&auth.Token{
				AccessToken: accessToken,
				Expiry:      expiry,
			}, nil).Times(1)

		//
		// And a fresh new proxy using that last
		proxy := NewTokenProxy(persistentRepository)

		// When we try to fetch the token from the proxy
		token, err := proxy.FetchToken(context.TODO())

		// Then no error should be reported
		require.NoError(t, err)

		// And the token should be stored on memory by the proxy
		require.NotNil(t, proxy.token)

		// And the token should match the one returned from the persistent repository
		require.Equal(t, accessToken, token.AccessToken)
		require.Equal(t, expiry, token.Expiry)

		// And the tokens should not be the same
		require.NotSame(t, token, proxy.token)
	})

	t.Run("should fetch the token from the persistent repository when its own is not valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Given a persistent repository which has a valid token
		persistentRepository := NewMockTokenRepository(ctrl)

		persistentRepository.EXPECT().FetchToken(gomock.AssignableToTypeOf(context.TODO())).Return(
			&auth.Token{
				AccessToken: accessToken,
				Expiry:      expiry,
			}, nil).Times(1)

		//
		// And a proxy using that last which contains an expired token
		proxy := NewTokenProxy(persistentRepository)

		proxy.token = &auth.Token{
			AccessToken: "this is an expired access token",
			Expiry:      time.Now().Add(-1 * time.Hour),
		}

		//
		// When we try to fetch the token from the proxy
		token, err := proxy.FetchToken(context.TODO())

		// Then no error should be reported
		require.NoError(t, err)

		// And the token should be stored on memory by the proxy
		require.NotNil(t, proxy.token)

		// And the token should match the one returned from the persistent repository
		require.Equal(t, accessToken, token.AccessToken)
		require.Equal(t, expiry, token.Expiry)

		// And the tokens should not be the same
		require.NotSame(t, token, proxy.token)
	})

	t.Run("should not overlap persistent repository calls", func(t *testing.T) {

	})

	t.Run("should return its own token when it is valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Given a persistent repository which has a valid token
		persistentRepository := NewMockTokenRepository(ctrl)

		persistentRepository.EXPECT().FetchToken(gomock.AssignableToTypeOf(context.TODO())).Return(
			&auth.Token{
				AccessToken: accessToken,
				Expiry:      expiry,
			}, nil).Times(0)

		//
		// And a proxy using that last which contains a different valid token
		proxy := NewTokenProxy(persistentRepository)

		proxy.token = &auth.Token{
			AccessToken: "this is a different but valid access token",
			Expiry:      time.Now().Add(1 * time.Hour),
		}

		//
		// When we try to fetch the token from the proxy
		token, err := proxy.FetchToken(context.TODO())

		// Then no error should be reported
		require.NoError(t, err)

		// And the token should match the one already on memory
		require.Equal(t, proxy.token.AccessToken, token.AccessToken)
		require.Equal(t, proxy.token.Expiry, token.Expiry)

		// And the tokens should not be the same
		require.NotSame(t, token, proxy.token)
	})
}

func TestTokenProxy_SaveToken(t *testing.T) {
	t.Run("should not be overwriten by fetch calls", func(t *testing.T) {

	})
}

func TestTokenProxyWithRandonExpirationDriftSeconds_FetchToken(t *testing.T) {

}
