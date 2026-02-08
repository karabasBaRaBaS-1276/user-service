package actuator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type contextKey string

const loggerKey contextKey = "requestLogger"

// MockDB для имитации работы sql.DB
type MockDB struct {
	err error
}

func (m *MockDB) Ping() error {
	return m.err
}

func TestLiveness(t *testing.T) {
	handler := Liveness()
	req := httptest.NewRequest("GET", "/liveness", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "UP")
}

func TestInfo(t *testing.T) {
	cfg := &config.Config{
		Environment: "test",
		Logging: config.LoggingConfig{
			ServiceName:    "test-app",
			ServiceVersion: "1.2.3",
		},
	}

	handler := Info(cfg)
	req := httptest.NewRequest("GET", "/info", nil)

	ctx := context.WithValue(req.Context(), loggerKey, zap.NewNop())
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, "test-app", resp["service"]) // Ключ из хендлера: "service"
	assert.Equal(t, "1.2.3", resp["version"])    // Ключ из хендлера: "version"
	assert.Equal(t, "test", resp["environment"])
}

func TestReadiness_Integration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))

		// Ожидаем пинг
		mock.ExpectPing()
		// Ожидаем запрос версии миграций
		rows := sqlmock.NewRows([]string{"version", "dirty"}).AddRow(10, false)
		mock.ExpectQuery("SELECT version, dirty FROM schema_migrations").WillReturnRows(rows)

		handler := Readiness(db)
		req := httptest.NewRequest("GET", "/readiness", nil)

		req = req.WithContext(context.WithValue(req.Context(), loggerKey, zap.NewNop()))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "UP")
		assert.Contains(t, rr.Body.String(), "version: 10")
	})

	t.Run("database_down", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
		mock.ExpectPing().WillReturnError(fmt.Errorf("db error"))

		handler := Readiness(db)
		req := httptest.NewRequest("GET", "/readiness", nil)
		req = req.WithContext(context.WithValue(req.Context(), loggerKey, zap.NewNop()))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
		assert.Contains(t, rr.Body.String(), "DOWN")
	})

	t.Run("dirty_migration", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
		mock.ExpectPing()
		rows := sqlmock.NewRows([]string{"version", "dirty"}).AddRow(10, true) // DIRTY = true
		mock.ExpectQuery("SELECT version, dirty FROM schema_migrations").WillReturnRows(rows)

		handler := Readiness(db)
		req := httptest.NewRequest("GET", "/readiness", nil)
		req = req.WithContext(context.WithValue(req.Context(), loggerKey, zap.NewNop()))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
		assert.Contains(t, rr.Body.String(), "dirty")
	})
}
