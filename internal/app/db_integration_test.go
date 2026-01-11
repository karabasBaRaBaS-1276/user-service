//go:build integration
// +build integration

package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"
)

var (
	testDB     *sql.DB
	dbHost     string
	dbPort     string
	dbName     = "testdb"
	dbUser     = "go_user_service"
	dbPassword = "go_user_service"
	dbSchema   = "user_service"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := postgres.Run(
		ctx,
		"postgres:15.4",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
	)
	if err != nil {
		panic(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		panic(err)
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		panic(err)
	}

	dbHost = host
	dbPort = port.Port()

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			Schema:   dbSchema,
			SSLMode:  "disable",
		},
	}

	logger := zap.NewNop()

	testDB, err = waitForDB(ctx, cfg, logger)
	if err != nil {
		panic(err)
	}

	code := m.Run()

	_ = testDB.Close()
	_ = container.Terminate(ctx)

	os.Exit(code)
}

func waitForDB(
	ctx context.Context,
	cfg *config.Config,
	logger *zap.Logger,
) (*sql.DB, error) {

	var lastErr error

	for i := 0; i < 10; i++ {
		db, err := initDB(cfg, logger)
		if err == nil {
			return db, nil
		}

		lastErr = err

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("database not ready: %w", ctx.Err())
		case <-time.After(1 * time.Second):
		}
	}

	return nil, fmt.Errorf("database not ready: %w", lastErr)
}

// cleanDatabase очищает все таблицы текущей схемы.
// Используется для изоляции integration-тестов.
func cleanDatabase(t *testing.T) {
	t.Helper()

	_, err := testDB.Exec(`
		DO $$
		DECLARE
			r RECORD;
		BEGIN
			FOR r IN (
				SELECT tablename
				FROM pg_tables
				WHERE schemaname = current_schema()
			) LOOP
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE';
			END LOOP;
		END $$;
	`)
	require.NoError(t, err)
}

func TestSomething(t *testing.T) {
	cleanDatabase(t)
	t.Cleanup(func() { cleanDatabase(t) })

	// тест
}

// Проверяем, что initDB корректно инициализировал соединение.
// Основная логика уже проверена в TestMain, здесь — smoke-test.
func TestInitDB_Integration(t *testing.T) {
	require.NotNil(t, testDB)
	require.NoError(t, testDB.Ping())
}
