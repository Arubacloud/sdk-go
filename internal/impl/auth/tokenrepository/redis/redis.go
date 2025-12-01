package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	redisClient IRedis
	clientId    string
}

var _ auth.TokenRepository = (*TokenRepository)(nil)

func NeWRedisTokenRepository(clientID string, r IRedis) *TokenRepository {
	return &TokenRepository{
		redisClient: r,
		clientId:    clientID,
	}
}

func (r *TokenRepository) FetchToken(ctx context.Context) (*auth.Token, error) {
	x := r.redisClient.Get(ctx, r.clientId)

	if x == nil {
		return nil, auth.ErrTokenNotFound
	}

	val, err := x.Result()

	if err != nil || val == "" {
		return nil, auth.ErrTokenNotFound
	}

	var token auth.Token
	if err := json.Unmarshal([]byte(val), &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *TokenRepository) SaveToken(ctx context.Context, token *auth.Token) error {
	if token == nil {
		return fmt.Errorf("token cannot be nil")
	}

	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}

	cmd := r.redisClient.Set(ctx, r.clientId, tokenJSON, time.Until(token.Expiry))

	return cmd.Err()
}

type IRedis interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, expiration time.Duration) *redis.StatusCmd
}
