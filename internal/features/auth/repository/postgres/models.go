package auth_postgres_repository

import "github.com/MisdGhost17/golang-todoapp/internal/core/domain"

type UserModel struct {
	ID          int     `postgres:"id"`
	Version     int     `postgres:"version"`
	FullName    string  `postgres:"full_name"`
	PhoneNumber *string `postgres:"phone_number"`
}

func userDomainFromModel(user UserModel) domain.User {
	return domain.NewUser(
		user.ID,
		user.Version,
		user.FullName,
		user.PhoneNumber,
	)
}

type CredentialsModel struct {
	UserID       int    `postgres:"user_id"`
	PasswordHash string `postgres:"password_hash"`
}

func credentialsDomainFromModel(credentials CredentialsModel) domain.Credentials {
	return domain.NewCredentials(
		credentials.UserID,
		credentials.PasswordHash,
	)
}
