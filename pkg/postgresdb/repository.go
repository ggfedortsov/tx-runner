package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	*pgxpool.Pool
}

func NewRepository(p *pgxpool.Pool) *PgRepository {
	return &PgRepository{Pool: p}
}

func (r *PgRepository) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.Query(ctx, sql, args...)
	}

	return r.Pool.Query(ctx, sql, args...)
}

func (r *PgRepository) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	if tx := txFromContext(ctx); tx != nil {
		return tx.QueryRow(ctx, sql, args...)
	}

	return r.Pool.QueryRow(ctx, sql, args...)
}

func (r *PgRepository) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	if tx := txFromContext(ctx); tx != nil {
		return tx.Exec(ctx, sql, args...)
	}

	return r.Pool.Exec(ctx, sql, args...)
}
