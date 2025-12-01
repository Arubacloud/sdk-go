package redis

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
	"github.com/Arubacloud/sdk-go/internal/restclient"
	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	token       *auth.Token
	locker      sync.RWMutex
	redisClient IRedis
	clientId    string
}

var _ auth.TokenRepository = (*TokenRepository)(nil)

func NeWRedisTokenRepository(cfg restclient.Config) *TokenRepository {
	opt, err := redis.ParseURL(cfg.Redis.RedisURI)

	if err != nil {
		log.Fatal("Cannot parse Redis URI", err)
		panic(err)
	}

	rdb := redis.NewClient(opt)

	return &TokenRepository{redisClient: rdb,
		clientId: cfg.ClientID}
}

func NewTokenRepositoryWithRedis(cfg restclient.Config, r IRedis) *TokenRepository {
	return &TokenRepository{
		redisClient: r,
		clientId:    cfg.ClientID,
	}
}

func (r *TokenRepository) FetchToken(ctx context.Context) (*auth.Token, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	if r.token != nil {
		return r.token.Copy(), nil
	}

	x, err := r.redisClient.Get(ctx, r.clientId).Result()

	if err != nil {
		return nil, auth.ErrTokenNotFound
	}

	if x == "" {
		return nil, auth.ErrTokenNotFound
	}

	var token auth.Token
	if err := json.Unmarshal([]byte(x), &token); err != nil {
		return nil, err
	}

	r.token = &token
	return &token, nil
}

func (r *TokenRepository) SaveToken(ctx context.Context, token *auth.Token) error {
	r.locker.Lock()
	defer r.locker.Unlock()

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
