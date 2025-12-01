package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
)

type CredentialsRepository struct {
	credentials          *auth.Credentials
	locker               sync.RWMutex
	persistentRepository auth.CredentialsRepository
}

var _ auth.CredentialsRepository = (*CredentialsRepository)(nil)

func NewCredentialsRepository(clientID string, clientSecret string) *CredentialsRepository {
	return &CredentialsRepository{
		credentials: &auth.Credentials{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
	}
}

func NewCredentialsProxy(persistentRepository auth.CredentialsRepository) *CredentialsRepository {
	return &CredentialsRepository{
		persistentRepository: persistentRepository,
	}
}

func (r *CredentialsRepository) FetchCredentials(ctx context.Context) (*auth.Credentials, error) {
	r.locker.RLock()
	defer r.locker.RUnlock()

	return nil, errors.New("not implemented")
}
