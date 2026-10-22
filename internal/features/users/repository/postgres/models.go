package users_postgres_repository

import (
	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
)

type UserModel struct {
	ID          int     `postgres:"id"`
	Version     int     `postgres:"version"`
	FullName    string  `postgres:"full_name"`
	PhoneNumber *string `postgres:"phone_number"`
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.ID,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}

	return userDomains
}
