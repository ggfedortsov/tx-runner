package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct {
	*pgxpool.Pool

	opts pgx.TxOptions
}

func NewRepository(p *pgxpool.Pool, opts pgx.TxOptions) *PgRepository {
	return &PgRepository{
		Pool: p,
		opts: opts,
	}
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

func (r *PgRepository) RunInTx(ctx context.Context, txFunc func(ctx context.Context) error) error {
	if tx := txFromContext(ctx); tx != nil {
		return pgx.BeginTxFunc(ctx, tx.Conn(), pgx.TxOptions{}, func(tx pgx.Tx) error {
			return txFunc(withTx(ctx, tx))
		})
	}

	return pgx.BeginTxFunc(ctx, r.Pool, r.opts, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
