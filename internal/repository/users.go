package repository

import (
	"context"

	"github.com/ggfedortsov/tx-runner/internal/model"
	"github.com/ggfedortsov/tx-runner/pkg/postgresdb"
	"github.com/jackc/pgx/v5"
)

type Users struct {
	*postgresdb.PgRepository
}

func (a *Users) CreateUser(ctx context.Context, u model.User) error {
	_, err := a.PgRepository.Exec(ctx, "INSERT INTO users(name, age) VALUES ($1, $2)", u.Username, u.Age)

	return err
}

func (a *Users) GatAll(ctx context.Context) ([]model.User, error) {
	rows, err := a.PgRepository.Query(ctx, "SELECT name, age FROM users")
	if err != nil {
		return nil, err
	}

	return pgx.CollectRows(rows, pgx.RowToStructByPos[model.User])
}
