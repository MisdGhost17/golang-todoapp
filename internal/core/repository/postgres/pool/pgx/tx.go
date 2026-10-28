package core_pgx_pool

import (
	"context"
	"fmt"

	core_postgres_pool "github.com/MisdGhost17/golang-todoapp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
)

type pgxTx struct {
	tx pgx.Tx
}

func (p *Pool) Begin(ctx context.Context) (core_postgres_pool.Tx, error) {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin pgx tx: %w", err)
	}

	return &pgxTx{tx: tx}, nil
}

func (t *pgxTx) QueryRow(
	ctx context.Context,
	sql string,
	args ...any,
) core_postgres_pool.Row {
	return pgxRow{t.tx.QueryRow(ctx, sql, args...)}
}

func (t *pgxTx) Exec(
	ctx context.Context,
	sql string,
	args ...any,
) (core_postgres_pool.CommandTag, error) {
	tag, err := t.tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}

	return pgxCommandTag{tag}, nil
}

func (t *pgxTx) Commit(ctx context.Context) error {
	if err := t.tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (t *pgxTx) Rollback(ctx context.Context) error {
	if err := t.tx.Rollback(ctx); err != nil {
		return fmt.Errorf("rollback transaction: %w", err)
	}

	return nil
}
