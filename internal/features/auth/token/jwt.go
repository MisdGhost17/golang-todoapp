package auth_token

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTManager(
	secret []byte,
	issuer string,
	ttl time.Duration,
) (*JWTManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("JWT secret must contain at least 32 bytes")
	}

	if strings.TrimSpace(issuer) == "" {
		return nil, fmt.Errorf("JWT issuer must not be empty")
	}

	if ttl < time.Second {
		return nil, fmt.Errorf("JWT TTL must be at least one second")
	}

	return &JWTManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}, nil
}

func (m *JWTManager) Issue(userID int) (string, error) {
	if userID <= 0 {
		return "", fmt.Errorf("invalid user ID")
	}

	now := time.Now()

	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		Issuer:    m.issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return signed, nil
}

func (m *JWTManager) Verify(tokenString string) (int, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (any, error) {
			return m.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer(m.issuer),
		jwt.WithIssuedAt(),
	)
	if err != nil || !token.Valid {
		return 0, fmt.Errorf(
			"invalid or expired access token: %w",
			core_errors.ErrUnauthorized,
		)
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil || userID <= 0 {
		return 0, fmt.Errorf(
			"invalid token subject: %w",
			core_errors.ErrUnauthorized,
		)
	}

	return userID, nil
}
