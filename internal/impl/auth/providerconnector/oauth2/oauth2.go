package oauth2

import (
	"context"
	"errors"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
)

type ProviderConnector struct {
	credentialsRepository auth.CredentialsRepository
	tokenURL              string
	scopes                []string
}

var _ auth.ProviderConnector = (*ProviderConnector)(nil)

func NewProviderConnector(credentialsRepository auth.CredentialsRepository, tokenURL string, scopes []string) *ProviderConnector {
	return &ProviderConnector{
		credentialsRepository: credentialsRepository,
		tokenURL:              tokenURL,
		scopes:                scopes,
	}
}

func (c *ProviderConnector) RequestToken(ctx context.Context) (*auth.Token, error) {
	return nil, errors.New("not implemented")
}
