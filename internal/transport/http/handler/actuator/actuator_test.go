package actuator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	req = req.WithContext(context.WithValue(req.Context(), loggerKey, zap.NewNop()))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp LivenessResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)

	require.NoError(t, err, "Ответ должен быть валидным JSON")
	assert.Equal(t, "UP", resp.Status)

	_, err = time.Parse(time.RFC3339, resp.Timestamp)
	assert.NoError(t, err, "Timestamp должен соответствовать формату RFC3339")
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

	var resp InfoResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)

	require.NoError(t, err, "Ответ должен соответствовать структуре InfoResponse")

	assert.Equal(t, cfg.Logging.ServiceName, resp.Service)
	assert.Equal(t, cfg.Logging.ServiceVersion, resp.Version)
	assert.Equal(t, cfg.Environment, resp.Environment)
}

func TestReadiness_Integration(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))

		mock.ExpectPing()
		rows := sqlmock.NewRows([]string{"version", "dirty"}).AddRow(10, false)
		mock.ExpectQuery("SELECT version, dirty FROM schema_migrations").WillReturnRows(rows)

		handler := Readiness(db)
		req := httptest.NewRequest("GET", "/readiness", nil)
		req = req.WithContext(context.WithValue(req.Context(), loggerKey, zap.NewNop()))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var resp ReadinessResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "UP", resp.Status)
		assert.Equal(t, "UP", resp.Checks["database"].Status)
		assert.Contains(t, resp.Checks["database"].Message, "Версия: 10")
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

		var resp ReadinessResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "DOWN", resp.Status)
		assert.Equal(t, "DOWN", resp.Checks["database"].Status)
		assert.Contains(t, resp.Checks["database"].Message, "db error")
	})

	t.Run("dirty_migration", func(t *testing.T) {
		db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
		mock.ExpectPing()
		rows := sqlmock.NewRows([]string{"version", "dirty"}).AddRow(10, true)
		mock.ExpectQuery("SELECT version, dirty FROM schema_migrations").WillReturnRows(rows)

		handler := Readiness(db)
		req := httptest.NewRequest("GET", "/readiness", nil)
		req = req.WithContext(context.WithValue(req.Context(), loggerKey, zap.NewNop()))
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)

		var resp ReadinessResponse
		err := json.Unmarshal(rr.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.Equal(t, "DOWN", resp.Status)
		assert.Contains(t, resp.Checks["database"].Message, "dirty")
	})
}
