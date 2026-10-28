package auth_service

import (
	"context"
	"fmt"
	"time"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	"github.com/google/uuid"
)

func (s *AuthService) Login(
	ctx context.Context,
	login string,
	password string,
) (domain.TokenPair, error) {
	userID, err := s.Authenticate(ctx, login, password)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("authenticate: %w", err)
	}

	accessToken, err := s.tokenIssuer.Issue(userID)
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("issue access token: %w", err)
	}

	refereshToken, tokenHash, err := s.refreshTokenManager.Generate()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate refresh token: %w", err)
	}

	sessionID, err := uuid.NewRandom()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate session random ID: %w", err)
	}

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return domain.TokenPair{}, fmt.Errorf("generate token ID: %w", err)
	}

	expiresAt := time.Now().UTC().Add(30 * 24 * time.Hour)

	if err := s.repository.CreateSession(
		ctx,
		sessionID,
		userID,
		tokenID,
		tokenHash,
		expiresAt,
	); err != nil {
		return domain.TokenPair{}, fmt.Errorf("create session: %w", err)
	}

	return domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refereshToken,
	}, nil
}
