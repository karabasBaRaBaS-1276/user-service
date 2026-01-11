package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"go.uber.org/zap"
)

const (
	dbDriver      = "pgx"
	dbPingTimeout = 5 * time.Second
)

func initDB(cfg *config.Config, logger *zap.Logger) (*sql.DB, error) {
	logger.Info("Инициализация подключения в базе данных")

	dsn, err := buildPostgresDSN(cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("build postgres dsn: %w", err)
	}

	logger.Sugar().Infof("Базовые параметры подключения в базе данных: host=%s, port=%s, database=%s, username=%s, ssl_mode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.User,
		cfg.Database.SSLMode)

	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %s", err)
	}
	logger.Info("Соединение с базой данных успешно установлено")

	// --- Работа со схемой
	schema := cfg.Database.Schema
	if schema == "" {
		return nil, fmt.Errorf("database schema is empty")
	}

	// CREATE SCHEMA
	createSchemaQuery := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)

	logger.Sugar().Infof("Проверка / создание схемы: '%s'", createSchemaQuery)
	if _, err := db.Exec(createSchemaQuery); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create schema %s: %w", schema, err)
	}

	// SET search_path
	setSearchPathQuery := fmt.Sprintf("SET search_path TO %s", schema)
	logger.Sugar().Infof("Установка search_path по умолчанию: %s", setSearchPathQuery)

	if _, err := db.Exec(setSearchPathQuery); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("set search_path to %s: %w", schema, err)
	}

	logger.Info("Схема базы данных успешно подготовлена для работы")

	return db, nil
}

func buildPostgresDSN(cfg config.DatabaseConfig) (string, error) {
	if cfg.Host == "" || cfg.Port == "" || cfg.User == "" || cfg.Name == "" {
		return "", fmt.Errorf("incomplete database configuration")
	}

	var user *url.Userinfo
	if cfg.Password != "" {
		user = url.UserPassword(cfg.User, cfg.Password)
	} else {
		user = url.User(cfg.User)
	}

	u := &url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Path:   cfg.Name,
	}

	q := u.Query()
	if cfg.SSLMode != "" {
		q.Set("sslmode", cfg.SSLMode)
	}
	if cfg.Schema != "" {
		q.Set("search_path", cfg.Schema)
	}
	q.Set("client_encoding", "UTF8")

	u.RawQuery = q.Encode()

	return u.String(), nil
}
