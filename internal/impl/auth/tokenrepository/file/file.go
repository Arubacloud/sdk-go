package file

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Arubacloud/sdk-go/internal/ports/auth"
)

type TokenRepository struct {
	baseDir string
	path    string
}

var _ auth.TokenRepository = (*TokenRepository)(nil)

func NeWFileTokenRepository(baseDir, clientID string) *TokenRepository {
	name := fmt.Sprintf("%s.token.json", clientID)
	return &TokenRepository{
		baseDir: baseDir,

		path: filepath.Join(baseDir, name),
	}
}

func (r *TokenRepository) FetchToken(ctx context.Context) (*auth.Token, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}

	var token auth.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}

	// Optional: check if token expired
	if time.Now().After(token.Expiry) {
		return nil, fmt.Errorf("token expired")
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

	// Ensure directory exists
	if err := os.MkdirAll(r.baseDir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(r.path, tokenJSON, 0o600)
}
