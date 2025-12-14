package logger

import (
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Конфигурация логгера
type Config struct {
	Level          string // Уровень логирования: debug, info, warn, error
	Environment    string // Окружение: development, production, staging
	FilePath       string // Файл для записи логов (если пусто - только stdout)
	MaxSizeMB      int    // Максимальный размер файла в MB
	MaxBackups     int    // Количество backup файлов
	MaxAgeDays     int    // Максимальный возраст файла в днях
	ServiceName    string // наименование сервиса
	ServiceVersion string // версия сервиса
	EnableJSON     bool   // принудительный JSON формат
	DisableConsole bool   // отключить вывод в консоль
}

// MiddlewareLogger структура для работы с логгером в middleware
type MiddlewareLogger struct {
	*zap.Logger
	ServiceName    string
	serviceVersion string
}

// New создает и настраивает новый экземпляр zap.Logger на основе конфигурации.
// Функция возвращает настроенный логгер или ошибку в случае проблем.
func New(cfg Config) (*zap.Logger, error) {

	// Валидация конфигурации
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown-service"
	}
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = "0.0.0"
	}

	// Определяем конфигурацию кодировщика
	encoderConfig := getEncoderConfig(cfg)

	// Парсинг уровня логирования
	level := getLogLevel(cfg.Level)

	// Создаем ядра логгера
	cores, err := createCores(cfg, encoderConfig, level)
	if err != nil {
		return nil, err
	}

	// Объединяем ядра
	core := zapcore.NewTee(cores...)

	// Добавляем сервисные поля ко всем логам
	serviceFields := []zap.Field{
		zap.String("service_id", cfg.ServiceName),
		zap.String("service_version", cfg.ServiceVersion),
		zap.String("environment", cfg.Environment),
		zap.String("hostname", getHostname()),
	}

	// Создаем финальный логгер с дополнительными опциями:
	// - zap.AddCaller(): добавляет информацию о месте вызова (файл:строка)
	// - zap.AddStacktrace(zap.ErrorLevel): добавляет стектрейс для ошибок
	return zap.New(core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
		zap.Fields(serviceFields...),
	), nil
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// LoggerWithRequestID создает новый логгер с добавленным request_id
func LoggerWithRequestID(logger *zap.Logger, requestID string) *zap.Logger {
	if requestID == "" {
		return logger
	}
	return logger.With(zap.String("request_id", requestID))
}

// NewMiddlewareLogger создает логгер для middleware
func NewMiddlewareLogger(logger *zap.Logger, ServiceName, serviceVersion string) *MiddlewareLogger {
	return &MiddlewareLogger{
		Logger:         logger,
		ServiceName:    ServiceName,
		serviceVersion: serviceVersion,
	}
}

// WithRequest создает логгер с полями HTTP запроса
func (ml *MiddlewareLogger) WithRequest(r *http.Request) *zap.Logger {
	fields := []zap.Field{
		zap.String("service_id", ml.ServiceName),
		zap.String("service_version", ml.serviceVersion),
		zap.String("http.method", r.Method),
		zap.String("http.path", r.URL.Path),
		zap.String("http.user_agent", r.UserAgent()),
		zap.String("http.referer", r.Referer()),
		zap.String("source.ip", r.RemoteAddr),
	}

	// Добавляем request_id если есть
	if requestID := r.Header.Get("X-Request-ID"); requestID != "" {
		fields = append(fields, zap.String("request_id", requestID))
	}

	return ml.Logger.With(fields...)
}

// getLogLevel парсит строку в zapcore.Level
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// getEncoderConfig возвращает конфигурацию кодировщика
func getEncoderConfig(cfg Config) zapcore.EncoderConfig {
	var encoderConfig zapcore.EncoderConfig

	// Для development используем человекочитаемый формат
	if cfg.Environment == "development" && !cfg.EnableJSON {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		// Для production или принудительного JSON
		encoderConfig = zap.NewProductionEncoderConfig()

		// Настраиваем для ELK/Filebeat
		encoderConfig.TimeKey = "@timestamp"
		encoderConfig.MessageKey = "message"

		// Настраиваем формат времени
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	return encoderConfig
}

// createCores создает ядра логгера
func createCores(cfg Config, encoderConfig zapcore.EncoderConfig, level zapcore.Level) ([]zapcore.Core, error) {
	var cores []zapcore.Core

	// Console core (stdout)
	if !cfg.DisableConsole {
		consoleEncoder := getConsoleEncoder(cfg, encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// File core для production или если указан путь
	if cfg.FilePath != "" && (cfg.Environment == "production" || cfg.EnableJSON) {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   true,
		}

		// Для файла всегда используем JSON (удобно для ELK)
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileCore := zapcore.NewCore(
			jsonEncoder,
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// Если ядер нет, создаем хотя бы одно (в stdout)
	if len(cores) == 0 {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	return cores, nil
}

// getConsoleEncoder возвращает кодировщик для консоли
func getConsoleEncoder(cfg Config, baseConfig zapcore.EncoderConfig) zapcore.Encoder {
	if cfg.EnableJSON {
		// Если включен принудительный JSON, используем его даже в консоли
		return zapcore.NewJSONEncoder(baseConfig)
	}

	if cfg.Environment == "development" {
		// Для development - читаемый формат
		return zapcore.NewConsoleEncoder(baseConfig)
	}

	// Для production в консоли тоже можно использовать JSON
	// или вернуть к консольному формату без цветов
	if cfg.Environment == "production" {
		prodConfig := baseConfig
		prodConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(prodConfig)
	}

	return zapcore.NewConsoleEncoder(baseConfig)
}
