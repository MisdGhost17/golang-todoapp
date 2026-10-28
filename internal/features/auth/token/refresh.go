package auth_token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type RefreshTokenManager struct{}

func NewRefreshTokenManager() *RefreshTokenManager {
	return &RefreshTokenManager{}
}

func (m *RefreshTokenManager) Generate() (string, []byte, error) {
	randBytes := make([]byte, 32)

	if _, err := rand.Read(randBytes); err != nil {
		return "", nil, fmt.Errorf("failed generate random bytes: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(randBytes)
	hash := m.Hash(token)

	return token, hash, nil
}

func (m *RefreshTokenManager) Hash(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}
