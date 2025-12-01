package memory

import (
	"context"
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

	var credentialsCopy *auth.Credentials

	if r.credentials != nil {
		credentialsCopy = r.credentials.Copy()
	}

	r.locker.RUnlock()

	if credentialsCopy == nil && r.persistentRepository != nil {
		r.locker.Lock()
		defer r.locker.Unlock()

		credentials, err := r.persistentRepository.FetchCredentials(ctx)
		if err != nil {
			return nil, err
		}

		r.credentials = credentials.Copy()

		credentialsCopy = credentials.Copy()
	}

	return credentialsCopy, nil
}
