package auth_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	core_errors "github.com/MisdGhost17/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/MisdGhost17/golang-todoapp/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

func (r *AuthRepository) RotateRefreshToken(
	ctx context.Context,
	oldTokenHash []byte,
	newTokenID uuid.UUID,
	newTokenHash []byte,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin refresh transaction: %w", err)
	}

	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)

		defer cleanupCancel()

		_ = tx.Rollback(cleanupCtx)
	}()

	var sessionID uuid.UUID

	err = tx.QueryRow(
		ctx, `
		SELECT session_id
		FROM todoapp.refresh_tokens
		WHERE token_hash = $1;`,
		oldTokenHash,
	).Scan(&sessionID)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return 0, core_errors.ErrUnauthorized
		}
		return 0, fmt.Errorf("find refresh session: %w", err)
	}

	var (
		userID           int
		sessionExpiresAt time.Time
		revokedAt        *time.Time
	)

	err = tx.QueryRow(
		ctx,
		`
		SELECT user_id, expires_at, revoked_at
		FROM todoapp.auth_sessions
		WHERE id = $1
		FOR UPDATE;
		`,
		sessionID,
	).Scan(
		&userID,
		&sessionExpiresAt,
		&revokedAt,
	)
	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return 0, core_errors.ErrUnauthorized
		}
		return 0, fmt.Errorf("lock refresh session: %w", err)
	}

	var (
		usedAt         *time.Time
		tokenExpiresAt time.Time
		dbNow          time.Time
	)

	err = tx.QueryRow(
		ctx,
		`
		SELECT used_at, expires_at, clock_timestamp()
		FROM todoapp.refresh_tokens
		WHERE token_hash = $1 AND session_id = $2;
		`,
		oldTokenHash,
		sessionID,
	).Scan(
		&usedAt,
		&tokenExpiresAt,
		&dbNow,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return 0, core_errors.ErrUnauthorized
		}
		return 0, fmt.Errorf("read refresh token: %w", err)
	}

	if revokedAt != nil || !sessionExpiresAt.After(dbNow) {
		return 0, core_errors.ErrUnauthorized
	}

	if usedAt != nil {
		_, err = tx.Exec(
			ctx,
			`
			UPDATE todoapp.auth_sessions
			SET revoked_at = clock_timestamp()
			WHERE id = $1;
			`,
			sessionID,
		)
		if err != nil {
			return 0, fmt.Errorf("revoke reused session: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return 0, fmt.Errorf("commit session revocation: %w", err)
		}

		return 0, core_errors.ErrUnauthorized
	}

	if !tokenExpiresAt.After(dbNow) {
		return 0, core_errors.ErrUnauthorized
	}

	_, err = tx.Exec(
		ctx,
		`
		UPDATE todoapp.refresh_tokens
		SET used_at = clock_timestamp()
		WHERE token_hash = $1;
		`,
		oldTokenHash,
	)
	if err != nil {
		return 0, fmt.Errorf("mark refresh token used: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO todoapp.refresh_tokens (
			id, session_id, token_hash, expires_at
		)
		VALUES ($1, $2, $3, $4);
		`,
		newTokenID,
		sessionID,
		newTokenHash,
		sessionExpiresAt,
	)

	if err != nil {
		return 0, fmt.Errorf("insert replace refresh: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit refresh rotation: %w", err)
	}

	return userID, nil
}
