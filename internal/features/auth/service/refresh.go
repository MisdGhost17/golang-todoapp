package auth_service

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
	"github.com/google/uuid"
)

func (s *AuthService) Refresh(
	ctx context.Context,
	refreshToken string,
) (domain.TokenPair, error) {
	if len(refreshToken) != 43 {
		return domain.TokenPair{}, core_errors.ErrUnauthorized
	}

	decoded, err := base64.RawURLEncoding.Strict().DecodeString(refreshToken)
	if err != nil || len(decoded) != 32 {
		return domain.TokenPair{}, core_errors.ErrUnauthorized
	}

	oldTokenHash := s.refreshTokenManager.Hash(refreshToken)

	newRefreshToken, newTokenHash, err := s.refreshTokenManager.Generate()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	newTokenID, err := uuid.NewRandom()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate refresh token random ID: %w", err)
	}

	userID, err := s.repository.RotateRefreshToken(
		ctx,
		oldTokenHash,
		newTokenID,
		newTokenHash,
	)

	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("rotate refresh token: %w", err)
	}

	accessToken, err := s.tokenIssuer.Issue(userID)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("issue access token: %w", err)
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
