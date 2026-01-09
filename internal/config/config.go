package config

import (
	"embed"
	"fmt"
	"os"
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

// Config содержит итоговую конфигурацию приложения,
// загружаемую из файла конфигурации и/или переменных окружения.
type Config struct {
	// Environment — имя окружения выполнения приложения.
	// Обычно используется для переключения поведения и конфигурации.
	//
	// Типичные значения:
	//   - development
	//   - staging
	//   - production
	Environment string `json:"environment" yaml:"environment"`

	// ServiceName — логическое имя сервиса.
	// Используется в логировании, метриках и трассировке.
	ServiceName string `json:"service_name" yaml:"service_name"`

	// Server — конфигурация HTTP-сервера приложения.
	Server ServerConfig `json:"server" yaml:"server"`

	// Database — конфигурация подключения к базе данных.
	Database DatabaseConfig `json:"database" yaml:"database"`

	// Logging — конфигурация системы логирования.
	Logging LoggingConfig `json:"logging" yaml:"logging"`

	// Auth — конфигурация модуля авторизации и аутентификации.
	Auth AuthConfig `json:"auth" yaml:"auth"`

	// Metric — конфигурация подсистемы сбора и экспорта метрик.
	Metric MetricConfig `json:"metric" yaml:"metric"`
}

// BuildInfo содержит информацию о сборке приложения.
//
// Значения данной структуры не загружаются из конфигурационных файлов
// и, как правило, задаются на этапе сборки (например, через ldflags).
type BuildInfo struct {
	// Version — версия приложения.
	Version string

	// Commit — идентификатор коммита системы контроля версий.
	Commit string

	// BuildTime — дата и время сборки приложения.
	BuildTime string
}

// ServerConfig описывает параметры HTTP-сервера приложения.
type ServerConfig struct {
	// Port — порт, на котором сервер принимает входящие соединения.
	Port string `json:"port" yaml:"port"`

	// ReadTimeout — максимальное время чтения запроса целиком,
	// включая тело запроса.
	ReadTimeout time.Duration `json:"read_timeout" yaml:"read_timeout"`

	// WriteTimeout — максимальное время записи ответа клиенту.
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`

	// IdleTimeout — максимальное время ожидания следующего запроса
	// при использовании keep-alive соединений.
	IdleTimeout time.Duration `json:"idle_timeout" yaml:"idle_timeout"`

	// ShutdownTimeout — максимальное время, отводимое
	// на корректное завершение работы сервера (graceful shutdown).
	ShutdownTimeout time.Duration `json:"shutdown_timeout" yaml:"shutdown_timeout"`
}

// MetricConfig описывает параметры подсистемы сбора метрик.
type MetricConfig struct {
	// Enable — флаг включения экспорта метрик.
	Enable bool `json:"enable" yaml:"enable"`

	// Port — порт, на котором публикуется endpoint метрик
	// (например, для Prometheus).
	Port string `json:"port" yaml:"port"`
}

// DatabaseConfig описывает параметры подключения к PostgreSQL.
type DatabaseConfig struct {
	// Host — хост для подключения к серверу БД.
	Host string `json:"host" yaml:"host"`

	// Port — порт для подключения к серверу БД.
	Port string `json:"port" yaml:"port"`

	// User — имя технической учетной записи (ТУЗ) для подключения к БД.
	User string `json:"user" yaml:"user"`

	// Password — пароль технической учетной записи (ТУЗ).
	Password string `json:"password" yaml:"password"`

	// Name — имя базы данных.
	Name string `json:"name" yaml:"name"`

	// Schema — схема внутри указанной базы данных,
	// используемая для работы с таблицами.
	Schema string `json:"schema" yaml:"schema"`

	// SSLMode определяет режим использования SSL при подключении.
	//
	// Допустимые значения (при корректной настройке сервера БД):
	//   - disable     — SSL не используется.
	//   - allow       — SSL используется, если на этом настаивает сервер.
	//   - prefer      — SSL используется предпочтительно, если поддерживается сервером.
	//   - require     — SSL обязателен, без проверки сертификата сервера.
	//   - verify-ca   — SSL обязателен, с проверкой центра сертификации.
	//   - verify-full — SSL обязателен, с полной проверкой сертификата и имени сервера.
	//
	// Подробнее см. документацию PostgreSQL:
	// https://postgrespro.ru/docs/postgresql/current/libpq-ssl
	SSLMode string `json:"ssl_mode" yaml:"ssl_mode"`
}

// LoggingConfig описывает конфигурацию системы логирования приложения.
type LoggingConfig struct {
	// Level — минимальный уровень логирования.
	//
	// Типичные значения:
	//   - debug
	//   - info
	//   - warn
	//   - error
	Level string `json:"level" yaml:"level"`

	// Environment — имя окружения выполнения приложения.
	// Заполняется из основного Config и не задаётся явно в конфигурационном файле.
	Environment string

	// FilePath — путь к файлу логов.
	// Если не задан, логирование в файл может быть отключено.
	FilePath string `json:"file_path" yaml:"file_path"`

	// MaxSizeMB — максимальный размер одного файла логов в мегабайтах
	// перед его ротацией.
	MaxSizeMB int `json:"max_size_mb" yaml:"max_size_mb"`

	// MaxBackups — максимальное количество архивных файлов логов,
	// сохраняемых после ротации.
	MaxBackups int `json:"max_backups" yaml:"max_backups"`

	// MaxAgeDays — максимальный срок хранения файлов логов в днях.
	MaxAgeDays int `json:"max_age_days" yaml:"max_age_days"`

	// ServiceName — логическое имя сервиса.
	// Заполняется из основного Config и используется в структурированных логах.
	ServiceName string

	// ServiceVersion — версия сервиса.
	// Заполняется из данных сборки и используется в структурированных логах.
	ServiceVersion string

	// EnableJSON — включает вывод логов в структурированном JSON-формате.
	// Как правило, используется в production-окружении.
	EnableJSON bool `json:"enable_json" yaml:"enable_json"`

	// DisableConsole — отключает вывод логов в stdout/stderr.
	// Полезно при использовании исключительно файлового логирования.
	DisableConsole bool `json:"disable_console" yaml:"disable_console"`
}

// AuthConfig описывает конфигурацию модуля аутентификации и авторизации.
type AuthConfig struct {
	// JWTSecret — секретный ключ для подписи и валидации JWT-токенов.
	// Должен храниться в защищённом виде и не коммититься в репозиторий.
	JWTSecret string `json:"jwt_secret" yaml:"jwt_secret"`

	// TokenDuration — срок действия JWT-токена.
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

// Load выполняет загрузку конфигурации
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
// NOTE:
// default config is embedded at compile time.
// ReadFile / Unmarshal errors are unreachable in a valid build.
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
				panic(fmt.Sprintf("Environment variable %s is not set and no default", envKey))
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
