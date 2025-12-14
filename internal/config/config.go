package config

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	buildInfo     *BuildInfo
	buildInfoOnce sync.Once
	//go:embed default.yaml
	defaultConfigFS embed.FS
)

const defaultConfigPath = "default.yaml"

// Итоговая конфигурация приложения
type Config struct {
	Environment string         `json:"environment" yaml:"environment"`   // Окружение: development, production, staging
	ServiceName string         `json:"service_name" yaml:"service_name"` // Имя сервиса
	Server      ServerConfig   `json:"server" yaml:"server"`             // Конфигурация сервера
	Database    DatabaseConfig `json:"database" yaml:"database"`         // Конфигурация базы данных
	Logging     LoggingConfig  `json:"logging" yaml:"logging"`           // Конфигурация логгера
	Auth        AuthConfig     `json:"auth" yaml:"auth"`                 // Конфигурация модуля авторизации
	Metric      MetricConfig   `json:"metric" yaml:"metric"`             // Конфигурация для сбора метрик
}

// Информация о сборке приложения
type BuildInfo struct {
	Version   string // Не может быть определено из файла конфигурации
	Commit    string // Не может быть определено из файла конфигурации
	BuildTime string // Не может быть определено из файла конфигурации
}

// Конфигурация сервера
type ServerConfig struct {
	Port            string        `json:"port" yaml:"port"`
	ReadTimeout     time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout     time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" yaml:"shutdown_timeout"`
}

// Конфигурация для сбора метрик
type MetricConfig struct {
	Enable bool   `json:"enable" yaml:"enable"`
	Port   string `json:"port" yaml:"port"`
}

// Конфигурация базы данных
type DatabaseConfig struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Name     string `json:"name" yaml:"name"`
	SSLMode  string `json:"ssl_mode" yaml:"ssl_mode"`
}

// Конфигурация логгера
type LoggingConfig struct {
	Level          string `json:"level" yaml:"level"`
	Environment    string // Заполняется из основного Config
	FilePath       string `json:"file_path" yaml:"file_path"`
	MaxSizeMB      int    `json:"max_size_mb" yaml:"max_size_mb"`
	MaxBackups     int    `json:"max_backups" yaml:"max_backups"`
	MaxAgeDays     int    `json:"max_age_days" yaml:"max_age_days"`
	ServiceName    string // Заполняется из основного Config
	ServiceVersion string // Заполняется из данных сборки
	EnableJSON     bool   `json:"enable_json" yaml:"enable_json"`
	DisableConsole bool   `json:"disable_console" yaml:"disable_console"`
}

// Конфигурация модуля авторизации
type AuthConfig struct {
	JWTSecret     string        `json:"jwt_secret" yaml:"jwt_secret"`
	TokenDuration time.Duration `json:"token_duration" yaml:"token_duration"`
}

// SetBuildInfo устанавливает информацию о сборке (вызывается из main)
func SetBuildInfo(version, commit, buildTime string) {
	buildInfoOnce.Do(func() {
		buildInfo = &BuildInfo{
			Version:   version,
			Commit:    commit,
			BuildTime: buildTime,
		}
	})
}

// GetBuildInfo возвращает информацию о сборке
func GetBuildInfo() *BuildInfo {
	if buildInfo == nil {
		return &BuildInfo{
			Version:   "dev",
			Commit:    "unknown",
			BuildTime: time.Now().Format(time.RFC3339),
		}
	}
	return buildInfo
}

// Загрузка конфигурации
func Load(configPath string, env string) (*Config, error) {

	cfg := &Config{}

	// 1. Загружаем конфигурацию по умолчанию из embedded файла
	if err := loadDefaultConfig(cfg); err != nil {
		return nil, fmt.Errorf("load default config: %w", err)
	}

	// 2. Загружаем из файла если указан путь (переопределяет default)
	if configPath != "" {
		if err := loadYAMLWithEnv(cfg, configPath); err != nil {
			return nil, fmt.Errorf("load from file %s: %w", configPath, err)
		}
	}

	// Переопределяем окружение если указано флагом
	if env != "" {
		cfg.Environment = env
	}

	cfg.Logging.Environment = cfg.Environment
	cfg.Logging.ServiceName = cfg.ServiceName
	cfg.Logging.ServiceVersion = GetBuildInfo().Version

	return cfg, nil
}

// loadDefaultConfig загружает встроенную конфигурацию по умолчанию
func loadDefaultConfig(cfg *Config) error {
	data, err := defaultConfigFS.ReadFile(defaultConfigPath)
	if err != nil {
		return fmt.Errorf("read embedded config: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parse embedded config: %w", err)
	}

	return nil
}

// loadYAMLWithEnv загружает YAML с подстановкой переменных окружения
func loadYAMLWithEnv(cfg *Config, configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Подставляем переменные окружения
	expanded := os.Expand(string(data), func(key string) string {
		// Поддержка ${VAR:-default}
		if idx := strings.Index(key, ":-"); idx != -1 {
			envKey := key[:idx]         // часть до ":-"
			defaultValue := key[idx+2:] // часть после ":-"
			if val := os.Getenv(envKey); val != "" {
				return val // возвращаем значение env
			}
			if defaultValue == "" {
				log.Printf("WARN: Environment variable %s is not set and no default", envKey)
			}
			return defaultValue
		}
		val := os.Getenv(key)
		if val == "" {
			panic(fmt.Sprintf("Required environment variable %s is not set", key))
		}
		return val
	})

	return yaml.Unmarshal([]byte(expanded), cfg)
}

// Загрузка строковых переменных окружения
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Загрузка целочисленных переменных окружения
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Загрузка булевых переменных окружения
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

// Загрузка переменных окружения с отрезками времени
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
