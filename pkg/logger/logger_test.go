package logger

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNew_JSONLogger_Output(t *testing.T) {
	// перехват stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	cfg := Config{
		Environment: "production",
		EnableJSON:  true,
		ServiceName: "test-service",
	}

	logger, err := New(cfg)
	require.NoError(t, err)

	logger.Info("hello prod")

	_ = w.Close()
	out, _ := io.ReadAll(r)

	output := strings.TrimSpace(string(out))
	payload := assertValidJSONLog(t, output)

	assert.Equal(t, "hello prod", payload["message"])
	assert.Equal(t, "test-service", payload["service_id"])
	assert.Equal(t, "info", payload["level"])
}

func assertValidJSONLog(t *testing.T, output string) map[string]any {
	t.Helper()

	var payload map[string]any
	err := json.Unmarshal([]byte(output), &payload)

	require.NoError(t, err, "лог не является валидным JSON")
	require.NotEmpty(t, payload)

	return payload
}

func TestNew_ProductionHumanLogger_Output(t *testing.T) {
	// перехват stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	cfg := Config{
		Environment: "production",
		EnableJSON:  false,
		ServiceName: "test-service",
	}

	logger, err := New(cfg)
	require.NoError(t, err)

	logger.Info("hello dev")

	_ = w.Close()
	out, _ := io.ReadAll(r)

	output := strings.TrimSpace(string(out))

	// не JSON
	assert.False(t, strings.HasPrefix(output, "{"), "dev лог не должен быть JSON")

	// содержит сообщение
	assert.Contains(t, output, "hello dev")

	// содержит service_id
	assert.Contains(t, output, "test-service")
}

func TestNew_DevelopmentHumanLogger_Output(t *testing.T) {
	// перехват stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	cfg := Config{
		Environment: "development",
		EnableJSON:  false,
	}

	logger, err := New(cfg)
	require.NoError(t, err)

	logger.Info("hello prod")

	_ = w.Close()
	out, _ := io.ReadAll(r)

	output := strings.TrimSpace(string(out))

	// не JSON
	assert.False(t, strings.HasPrefix(output, "{"), "prod лог не должен быть JSON")

	// содержит сообщение
	assert.Contains(t, output, "hello prod")

	// содержит service_id по умолчанию
	assert.Contains(t, output, "unknown-service")
}

func TestNew_StagingLogger_Output(t *testing.T) {
	// перехват stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	cfg := Config{
		Environment: "staging",
	}

	logger, err := New(cfg)
	require.NoError(t, err)

	logger.Info("hello staging")

	_ = w.Close()
	out, _ := io.ReadAll(r)

	output := strings.TrimSpace(string(out))

	// не JSON
	assert.False(t, strings.HasPrefix(output, "{"), "staging лог не должен быть JSON")

	// содержит сообщение
	assert.Contains(t, output, "hello staging")

	// содержит service_id по умолчанию
	assert.Contains(t, output, "unknown-service")
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected zapcore.Level
	}{
		{"debug", "debug", zapcore.DebugLevel},
		{"warn", "warn", zapcore.WarnLevel},
		{"error", "error", zapcore.ErrorLevel},
		{"default_info", "info", zapcore.InfoLevel},
		{"unknown_fallback", "xxx", zapcore.InfoLevel},
		{"empty_fallback", "", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := getLogLevel(tt.input)
			assert.Equal(t, tt.expected, level, "Неожидаемый уровень логирования")
		})
	}
}

func TestCreateCores(t *testing.T) {
	// Тестируем создание разных ядер логгера в зависимости от конфигурации
	tests := []struct {
		name   string
		config Config
		len    int
	}{
		{"dev_cores", Config{Level: "debug", Environment: "development", DisableConsole: true, FilePath: "", EnableJSON: false}, 1},
		{"prod_cores", Config{Level: "error", Environment: "production", DisableConsole: false, FilePath: "test", EnableJSON: true}, 2},
	}

	for _, tt := range tests {
		encoderConfig := getEncoderConfig(tt.config)
		level := getLogLevel(tt.config.Level)
		t.Run(tt.name, func(t *testing.T) {
			cores, err := createCores(tt.config, encoderConfig, level)
			require.NoError(t, err, "не должно быть ошибки")
			assert.Len(t, cores, tt.len)
		})
	}
}

func TestLoggerWithRequestID(t *testing.T) {
	core, observed := observer.New(zapcore.InfoLevel)
	logger := zap.New(core)

	t.Run("empty request id", func(t *testing.T) {
		newLogger := LoggerWithRequestID(logger, "")
		assert.Equal(t, logger, newLogger)
	})

	t.Run("with request id", func(t *testing.T) {
		reqID := "test-123"
		newLogger := LoggerWithRequestID(logger, reqID)
		newLogger.Info("test message")

		require.Equal(t, 1, observed.Len())
		assert.Equal(t, reqID, observed.All()[0].ContextMap()["request_id"])
	})
}

func TestWithRequest(t *testing.T) {
	// 1. Создаем логгер с базовыми сервисными полями (имитируем работу функции New)
	serviceName := "my-service"
	serviceVersion := "1.2.3"

	core, observed := observer.New(zapcore.InfoLevel)
	// Добавляем поля так же, как это делает функция New
	logger := zap.New(core, zap.Fields(
		zap.String("service_id", serviceName),
		zap.String("service_version", serviceVersion),
	))

	// 2. Готовим запрос
	req, _ := http.NewRequest("GET", "/test-path", nil)
	req.Header.Set("X-Request-ID", "req-777")
	req.Header.Set("User-Agent", "Go-Test")
	req.RemoteAddr = "127.0.0.1:1234"

	// 3. Вызываем обновленную функцию
	l := WithRequest(logger, req)
	l.Info("request logged")

	// 4. Проверяем результат
	require.Equal(t, 1, observed.Len())
	fields := observed.All()[0].ContextMap()

	// Проверяем, что базовые поля сохранились
	assert.Equal(t, serviceName, fields["service_id"])
	assert.Equal(t, serviceVersion, fields["service_version"])

	// Проверяем, что новые поля добавились корректно
	assert.Equal(t, "GET", fields["http_method"])
	assert.Equal(t, "/test-path", fields["http_path"])
	assert.Equal(t, "req-777", fields["request_id"])
	assert.Equal(t, "Go-Test", fields["http_user_agent"])
}
