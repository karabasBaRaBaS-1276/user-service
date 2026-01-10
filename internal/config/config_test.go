package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("config must not be nil")
	}
	//
	assert.NotEmpty(t, cfg.Environment, "Окружение по умолчанию")
	assert.NotEmpty(t, cfg.ServiceName, "Наименование сервиса")
	// сервер
	assert.NotEmpty(t, cfg.Server.Port)
	assert.NotZero(t, cfg.Server.ReadTimeout)
	assert.NotZero(t, cfg.Server.WriteTimeout)
	assert.NotZero(t, cfg.Server.IdleTimeout)
	assert.NotZero(t, cfg.Server.ShutdownTimeout)
	// база данных
	assert.NotEmpty(t, cfg.Database.Host)
	assert.NotEmpty(t, cfg.Database.Port)
	assert.NotEmpty(t, cfg.Database.Name)
	assert.NotEmpty(t, cfg.Database.User)
	assert.NotEmpty(t, cfg.Database.Password)
	assert.NotEmpty(t, cfg.Database.SSLMode)
	// логирование
	assert.NotEmpty(t, cfg.Logging.Level)
	assert.NotEmpty(t, cfg.Logging.Environment)
	assert.NotEmpty(t, cfg.Logging.FilePath)
	assert.NotZero(t, cfg.Logging.MaxSizeMB)
	assert.NotZero(t, cfg.Logging.MaxBackups)
	assert.NotZero(t, cfg.Logging.MaxAgeDays)
	assert.NotEmpty(t, cfg.Logging.ServiceName)
	assert.NotEmpty(t, cfg.Logging.ServiceVersion)
	assert.False(t, cfg.Logging.EnableJSON)
	assert.False(t, cfg.Logging.DisableConsole)
	// авторизация
	assert.NotEmpty(t, cfg.Auth.JWTSecret)
	assert.NotZero(t, cfg.Auth.TokenDuration)
	// метрики
	assert.False(t, cfg.Metric.Enable)
	assert.NotZero(t, cfg.Metric.Port)
}

func TestLoad_BuldInfo(t *testing.T) {
	defer func() {
		SetBuildInfo("", "", "")
	}()
	strVersion := "1.0.0-RC1"
	strCommit := "7eee85bbb0af9f8117c5c05cad03c11c986df37c"
	strBuildTime := "2025-12-16"

	SetBuildInfo(strVersion, strCommit, strBuildTime)

	cfg, err := Load("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	//
	assert.Equal(t, strVersion, cfg.Logging.ServiceVersion, "Версия сервиса")
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/no/such/path/config", "")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config")

	content := []byte(`
environment: "staging"
server:
  port: "37500"
database:
  ssl_mode: "enable"
`)

	err := os.WriteFile(configPath+"test01.yaml", content, 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(configPath+"test01.yaml", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	//
	assert.Equal(t, "staging", cfg.Environment, "Новое окружение")
	assert.NotEmpty(t, cfg.ServiceName, "Наименование должно остаться из настроек по умолчанию")
	// сервер
	assert.Equal(t, "37500", cfg.Server.Port, "Новый порт у сервера")
	// база данных
	assert.NotEmpty(t, cfg.Database.Host)
	assert.Equal(t, "enable", cfg.Database.SSLMode, "Поддрежка SSL у сервера БД")
}

func TestLoad_FromFileAndEnv(t *testing.T) {
	pass := "test_password"
	jwtSecret := "jwt_secret"
	t.Setenv("DB_PASSWORD", pass)
	t.Setenv("JWT_SECRET", jwtSecret)
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config")

	content := []byte(`
database:
  host: "${DB_HOST:-localhost}"
  password: "${DB_PASSWORD}"
auth:
  jwt_secret: "${JWT_SECRET:-test}"
`)

	err := os.WriteFile(configPath+"test02.yaml", content, 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(configPath+"test02.yaml", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	//
	assert.Equal(t, "localhost", cfg.Database.Host, "url до сервера БД")
	assert.Equal(t, pass, cfg.Database.Password, "пароль к БД")
	assert.Equal(t, jwtSecret, cfg.Auth.JWTSecret, "Секрет из переменной окружения")
}

func TestLoad_PanicNotDefaultEnv(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config")

	content := []byte(`
database:
  password: "${DB_PASSWORD:-}"
`)

	err := os.WriteFile(configPath+"test03.yaml", content, 0644)
	require.NoError(t, err)

	assert.Panics(t, func() {
		_, err := Load(configPath+"test03.yaml", "")
		_ = err // unreachable, но явно фиксируем намерение
	})
}

func TestLoad_PanicNotEnv(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config")

	content := []byte(`
database:
  password: "${DB_PASSWORD}"
`)

	err := os.WriteFile(configPath+"test04.yaml", content, 0644)
	require.NoError(t, err)

	assert.Panics(t, func() {
		_, err := Load(configPath+"test04.yaml", "")
		_ = err // unreachable, но явно фиксируем намерение
	})
}

func TestLoad_ReplaceEnv(t *testing.T) {
	pass := "test_password"
	t.Setenv("DB_PASSWORD", pass)

	cfg, err := Load("", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	//
	assert.Equal(t, "test", cfg.Environment, "Имя окружения")
}
