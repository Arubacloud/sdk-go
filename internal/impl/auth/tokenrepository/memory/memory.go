package memory

import (
	"context"
	"math/rand/v2"
	"sync"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
)

type TokenRepository struct {
	token                  *auth.Token
	locker                 sync.RWMutex
	persistentRepository   auth.TokenRepository
	expirationDriftSeconds uint32
}

var _ auth.TokenRepository = (*TokenRepository)(nil)

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func NewTokenProxy(persistentRepository auth.TokenRepository) *TokenRepository {
	return &TokenRepository{
		persistentRepository: persistentRepository,
	}
}

func NewTokenProxyWithRandonExpirationDriftSeconds(persistentRepository auth.TokenRepository, maxExpirationDriftSeconds uint32) *TokenRepository {
	return &TokenRepository{
		persistentRepository:   persistentRepository,
		expirationDriftSeconds: rand.Uint32N(maxExpirationDriftSeconds),
	}
}

func (r *TokenRepository) FetchToken(ctx context.Context) (*auth.Token, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	if r.token == nil {
		return nil, auth.ErrTokenNotFound
	}

	return r.token.Copy(), nil
}

func (r *TokenRepository) SaveToken(ctx context.Context, token *auth.Token) error {
	r.locker.Lock()
	defer r.locker.Unlock()

	r.token = token.Copy()

	return nil
}
