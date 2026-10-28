package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/MisdGhost17/golang-todoapp/internal/core/domain"
	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/MisdGhost17/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *AuthRepository) GetCredentialsByLogin(
	ctx context.Context,
	login string,
) (domain.Credentials, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT user_id, password_hash
	FROM todoapp.user_credentials
	WHERE login = $1;
	`

	var credentialsModel CredentialsModel

	row := r.pool.QueryRow(ctx, query, login)

	err := row.Scan(
		&credentialsModel.UserID,
		&credentialsModel.PasswordHash,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Credentials{}, fmt.Errorf(
				"credentials not found: %w",
				core_errors.ErrNotFound,
			)
		}

		return domain.Credentials{}, fmt.Errorf(
			"get credentials by login: %w",
			err,
		)
	}

	credentials := credentialsDomainFromModel(credentialsModel)

	return credentials, nil
}
