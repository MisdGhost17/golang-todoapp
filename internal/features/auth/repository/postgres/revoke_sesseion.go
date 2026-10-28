package auth_postgres_repository

import (
	"context"
	"fmt"
)

func (r *AuthRepository) RevokeSessionByRefreshToken(
	ctx context.Context,
	tokenHash []byte,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.auth_sessions
		SET revoked_at = clock_timestamp()
		WHERE id = (
			SELECT session_id
			FROM todoapp.refresh_tokens
			WHERE token_hash = $1
		)
		AND revoked_at IS NULL;
	`

	if _, err := r.pool.Exec(ctx, query, tokenHash); err != nil {
		return fmt.Errorf("revoke auth session: %w", err)
	}

	return nil
}
