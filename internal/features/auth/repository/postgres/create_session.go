package auth_postgres_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *AuthRepository) CreateSession(
	ctx context.Context,
	sessionID uuid.UUID,
	userID int,
	tokenID uuid.UUID,
	tokenHash []byte,
	expiresAt time.Time,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin session transaction: %w", err)
	}

	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		_ = tx.Rollback(cleanupCtx)
	}()

	sessionQuery := `
	INSERT INTO todoapp.auth_sessions (
		id, user_id, expires_at
	)
	VALUES ($1, $2, $3);
	`

	_, err = tx.Exec(
		ctx,
		sessionQuery,
		sessionID,
		userID,
		expiresAt,
	)

	if err != nil {
		return fmt.Errorf("insert auth session: %w", err)
	}

	tokenQuery := `
	INSERT INTO todoapp.refresh_tokens (
		id, session_id, token_hash, expires_at
	)
	VALUES ($1, $2, $3, $4);
	`

	_, err = tx.Exec(
		ctx,
		tokenQuery,
		tokenID,
		sessionID,
		tokenHash,
		expiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert refresh token: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit auth session: %w", err)
	}

	return nil
}
