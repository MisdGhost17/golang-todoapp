package auth_service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
)

type AuthRepository interface {
	Register(
		ctx context.Context,
		user domain.User,
		login string,
		passwordHash string,
	) (domain.User, error)

	GetCredentialsByLogin(
		ctx context.Context,
		login string,
	) (domain.Credentials, error)

	CreateSession(
		ctx context.Context,
		sessionID uuid.UUID,
		userID int,
		tokenID uuid.UUID,
		tokenHash []byte,
		expiresAt time.Time,
	) error

	RotateRefreshToken(
		ctx context.Context,
		oldTokenHash []byte,
		newTokenID uuid.UUID,
		newTokenHash []byte,
	) (int, error)

	RevokeSessionByRefreshToken(
		ctx context.Context,
		tokenHash []byte,
	) error
}

type RefreshTokenManager interface {
	Generate() (string, []byte, error)
	Hash(token string) []byte
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(passwordHash string, password string) (bool, error)
}

type TokenIssuer interface {
	Issue(userID int) (string, error)
}

type AuthService struct {
	repository          AuthRepository
	passwordHasher      PasswordHasher
	tokenIssuer         TokenIssuer
	refreshTokenManager RefreshTokenManager
}

func NewAuthService(
	repository AuthRepository,
	passwordHasher PasswordHasher,
	tokenIssuer TokenIssuer,
	refreshTokenManager RefreshTokenManager,
) *AuthService {
	return &AuthService{
		repository:          repository,
		passwordHasher:      passwordHasher,
		tokenIssuer:         tokenIssuer,
		refreshTokenManager: refreshTokenManager,
	}
}
