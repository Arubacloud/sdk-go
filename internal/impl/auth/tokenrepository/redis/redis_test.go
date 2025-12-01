package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

//go:generate mockgen -package redis -destination=zz_mock_redis_test.go github.com/Arubacloud/sdk-go/internal/impl/auth/tokenrepository/redis IRedis

var (
	accessToken = "this is a valid token"
	expiry      = time.Now().Add(24 * time.Hour)
)

func TestTokenRepository_FetchToken(t *testing.T) {
	t.Run("should report a token not found error when it has not a token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRedis := NewMockIRedis(ctrl)
		tokenRepository := NeWRedisTokenRepository("user-123", mockRedis)
		mockRedis.
			EXPECT().
			Get(gomock.Any(), "user-123").
			Return(nil)

		// When we try to fetch the token
		token, err := tokenRepository.FetchToken(context.Background())

		// And no token should be returned
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrTokenNotFound)

		require.Nil(t, token)
	})

	t.Run("should report a token not found error when it has empty token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRedis := NewMockIRedis(ctrl)
		tokenRepository := NeWRedisTokenRepository("user-123", mockRedis)

		cmd := redis.NewStringCmd(context.Background())
		cmd.SetVal("") // Empty value to simulate no token

		mockRedis.
			EXPECT().
			Get(gomock.Any(), "user-123").
			Return(cmd)

		// When we try to fetch the token
		token, err := tokenRepository.FetchToken(context.Background())

		// And no token should be returned
		require.Error(t, err)
		require.ErrorIs(t, err, auth.ErrTokenNotFound)

		require.Nil(t, token)
	})

	t.Run("should report a token not found error when it has no jwt token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRedis := NewMockIRedis(ctrl)
		tokenRepository := NeWRedisTokenRepository("user-123", mockRedis)

		cmd := redis.NewStringCmd(context.Background())
		cmd.SetVal("token") // Empty value to simulate no token

		mockRedis.
			EXPECT().
			Get(gomock.Any(), "user-123").
			Return(cmd)

		// When we try to fetch the token
		token, err := tokenRepository.FetchToken(context.Background())

		// And no token should be returned
		require.Error(t, err)
		require.Nil(t, token)
	})

	t.Run("should report a token when it has a jwt token on redis", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRedis := NewMockIRedis(ctrl)
		tokenRepository := NeWRedisTokenRepository("user-123", mockRedis)

		cmd := redis.NewStringCmd(context.Background())

		authToken := &auth.Token{
			AccessToken: accessToken,
			Expiry:      expiry,
		}

		tokenJSON, _ := json.Marshal(authToken)
		cmd.SetVal(string(tokenJSON))

		mockRedis.
			EXPECT().
			Get(gomock.Any(), "user-123").
			Return(cmd)

		// When we try to fetch the token
		token, err := tokenRepository.FetchToken(context.Background())

		// And no token should be returned
		require.NoError(t, err)
		require.NotNil(t, token)
		require.Equal(t, accessToken, token.AccessToken)
		require.Equal(t, expiry.Unix(), token.Expiry.Unix())
	})
}

func TestTokenRepository_SaveToken(t *testing.T) {

	t.Run("should save a token if not nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRedis := NewMockIRedis(ctrl)

		// And a valid token
		token := &auth.Token{AccessToken: accessToken, Expiry: expiry}
		tokenRepository := NeWRedisTokenRepository("user-123", mockRedis)
		tokenJSON, _ := json.Marshal(token)

		cmd := redis.NewStatusCmd(context.Background())
		cmd.SetErr(nil)

		mockRedis.
			EXPECT().
			Set(gomock.Any(), "user-123", tokenJSON, gomock.Any()).
			Return(cmd)
		// When we try to save the token
		err := tokenRepository.SaveToken(context.TODO(), token)

		// Then no error should be reported
		require.NoError(t, err)
	})

	t.Run("should not save token if malformed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRedis := NewMockIRedis(ctrl)

		// And a valid token
		token := &auth.Token{
			AccessToken: string([]byte{0xff, 0xfe, 0xfd}), // Invalid UTF-8
			Expiry:      expiry,
		}
		tokenRepository := NeWRedisTokenRepository("user-123", mockRedis)

		tokenJSON, _ := json.Marshal(token)

		cmd := redis.NewStatusCmd(context.Background())
		cmd.SetErr(errors.New("malformed token"))

		mockRedis.
			EXPECT().
			Set(gomock.Any(), "user-123", tokenJSON, gomock.Any()).
			Return(cmd)
		// When we try to

		// When we try to save the token
		err := tokenRepository.SaveToken(context.TODO(), token)

		// Then no error should be reported
		require.Error(t, err)
		require.Equal(t, "malformed token", err.Error())
	})

}
