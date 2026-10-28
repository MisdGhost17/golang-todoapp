package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/MisdGhost17/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) Register(
	ctx context.Context,
	user domain.User,
	login string,
	passwordHash string,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	WITH new_user AS (
		INSERT INTO todoapp.users (full_name, phone_number)
		VALUES ($1, $2)
		RETURNING id, version, full_name, phone_number
	),
	new_credentials AS (
		INSERT INTO todoapp.user_credentials (
			user_id, login, password_hash
		)
		SELECT id, $3, $4
		FROM new_user
		RETURNING user_id
	)
	SELECT id, version, full_name, phone_number
	FROM new_user
	JOIN new_credentials ON new_credentials.user_id = new_user.id;
	`

	var userCreated UserModel

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		login,
		passwordHash,
	)

	err := row.Scan(
		&userCreated.ID,
		&userCreated.Version,
		&userCreated.FullName,
		&userCreated.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrViolatesUnique) {
			return domain.User{}, fmt.Errorf(
				"conflict with existing data: %w",
				core_errors.ErrConflict,
			)
		}

		return domain.User{}, fmt.Errorf("register user: %w", err)
	}

	userDomain := userDomainFromModel(userCreated)

	return userDomain, nil
}
