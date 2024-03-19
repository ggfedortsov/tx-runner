package postgresdb

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type txStarter interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Runner struct {
	db   txStarter
	opts pgx.TxOptions
}

func NewTxRunner(db txStarter, o pgx.TxOptions) *Runner {
	return &Runner{
		db:   db,
		opts: o,
	}
}

func (r *Runner) Run(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return pgx.BeginTxFunc(ctx, r.db, r.opts, func(tx pgx.Tx) error {
		return txFunc(withTx(ctx, tx))
	})
}
