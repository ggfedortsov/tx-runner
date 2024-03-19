package service

import (
	"context"
	"github.com/ggfedortsov/tx-runner/internal/model"
	"github.com/ggfedortsov/tx-runner/internal/repository"
	"github.com/ggfedortsov/tx-runner/pkg/postgresdb"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestService_MethodOk(t *testing.T) {
	ctx := context.Background()

	dbName := "users"
	dbUser := "user"
	dbPassword := "password"

	// 1. Start the postgres container and run any migrations on it
	container, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Run any migrations on the database
	_, _, err = container.Exec(ctx, []string{"psql", "-U", dbUser, "-d", dbName, "-c", "CREATE TABLE users (id SERIAL, name TEXT NOT NULL, age INT NOT NULL)"})
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a snapshot of the database to restore later
	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	dbURL, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	pool, err := postgresdb.NewPgPool(ctx, postgresdb.PgConfig{Conn: dbURL})
	if err != nil {
		t.Fatal(err)
	}

	r := postgresdb.NewRepository(pool)
	users := &repository.Users{Db: r}
	runner := postgresdb.NewTxRunner(pool, pgx.TxOptions{})

	service := Service{
		UserStorage: users,
		Runner:      runner,
	}

	us, err := service.MethodOk(ctx, model.User{
		Username: "alex",
		Age:      10,
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.MethodError(ctx, model.User{
		Username: "alex",
		Age:      10,
	})
	if err == nil {
		t.Fatal()
	}

	us, err = service.MethodOk(ctx, model.User{
		Username: "alex",
		Age:      10,
	})
	if err != nil {
		t.Fatal(err)
	}

	println(len(us))
}
