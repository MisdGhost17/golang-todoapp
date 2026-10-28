package auth_service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
)

var (
	loginPattern = regexp.MustCompile("^[a-z0-9_]{3,64}$")
)

func (s *AuthService) Register(
	ctx context.Context,
	fullName string,
	phoneNumber *string,
	login string,
	password string,
) (domain.User, error) {
	login = strings.ToLower(strings.TrimSpace(login))

	if !loginPattern.MatchString(login) {
		return domain.User{}, fmt.Errorf(
			"login length must be between 3 and 64 and must consist of: a-z, 0-9 or _: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if utf8.RuneCountInString(password) < 15 || len(password) > 72 {
		return domain.User{}, fmt.Errorf(
			"password must contain at least 15 characters and at most 72 bytes: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	user := domain.NewUserUninitialized(fullName, phoneNumber)

	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user: %w", err)
	}

	passwordHash, err := s.passwordHasher.Hash(password)
	if err != nil {
		return domain.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err = s.repository.Register(ctx, user, login, passwordHash)
	if err != nil {
		return domain.User{}, fmt.Errorf("register user: %w", err)
	}

	return user, nil
}
