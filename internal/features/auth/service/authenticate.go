package auth_service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
)

func (s *AuthService) Authenticate(
	ctx context.Context,
	login string,
	password string,
) (int, error) {
	login = strings.ToLower(strings.TrimSpace(login))

	var (
		invalidCredentials = fmt.Errorf(
			"invalid login or password: %w",
			core_errors.ErrUnauthorized,
		)
	)

	if !loginPattern.MatchString(login) ||
		len(password) == 0 ||
		len(password) > 72 {
		return 0, invalidCredentials
	}

	credentials, err := s.repository.GetCredentialsByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, core_errors.ErrNotFound) {
			return 0, invalidCredentials
		}

		return 0, fmt.Errorf("get credentials: %w", err)
	}

	matches, err := s.passwordHasher.Verify(
		credentials.PasswordHash,
		password,
	)

	if err != nil {
		return 0, fmt.Errorf("verify password: %w", err)
	}

	if !matches {
		return 0, invalidCredentials
	}

	return credentials.UserID, nil
}
