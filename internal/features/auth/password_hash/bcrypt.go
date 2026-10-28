package auth_passwordhash

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher {
	return &BcryptHasher{}
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf(
			"bcrypt generate hash: %w",
			err,
		)
	}

	return string(hash), nil
}

func (h *BcryptHasher) Verify(
	passwordHash string,
	password string,
) (bool, error) {
	if len(password) > 72 {
		return false, nil
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	)

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("bcrypt compare and hash: %w", err)
	}

	return true, nil
}
