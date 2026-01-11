package app

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func runMigrations(db *sql.DB, migrationsPath string, logger *zap.Logger) error {

	logger.Sugar().Infof("Старт DB миграций по пути: `%s`", migrationsPath)

	m, err := newMigrator(db, migrationsPath, logger)
	if err != nil {
		return err
	}

	// Текущая версия до применения
	version, dirty, err := m.Version()
	switch {
	case err == nil:
		logger.Sugar().Infof("Информация о миграции: version = %d, dirty = %t", version, dirty)
	case err == migrate.ErrNilVersion:
		logger.Info("Информация о миграции: нет информации о миграций")
	default:
		return fmt.Errorf("get migration version: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}

	// Версия после применения
	newVersion, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("get migration version after apply: %w", err)
	}

	logger.Sugar().Infof("Миграция базы данных успешно выполнена: version = %d, dirty = %t", newVersion, dirty)

	return nil
}

func newMigrator(db *sql.DB, migrationsPath string, _ *zap.Logger) (*migrate.Migrate, error) {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return m, nil
}
