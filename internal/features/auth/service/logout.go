package auth_service

import (
	"context"
	"encoding/base64"
	"fmt"
)

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	if len(refreshToken) != 43 {
		return nil
	}

	decoded, err := base64.RawURLEncoding.Strict().DecodeString(refreshToken)
	if err != nil || len(decoded) != 32 {
		return nil
	}

	tokenHash := s.refreshTokenManager.Hash(refreshToken)
	if err := s.repository.RevokeSessionByRefreshToken(ctx, tokenHash); err != nil {
		return fmt.Errorf("logout: %w", err)
	}

	return nil
}
