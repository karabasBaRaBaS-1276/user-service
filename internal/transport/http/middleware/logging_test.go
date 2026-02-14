package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLogging(t *testing.T) {
	logger := zap.NewNop()
	called := false

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true // подтверждаем, что логгер передал управление дальше
		w.WriteHeader(http.StatusNoContent)
	})

	req := httptest.NewRequest("GET", "/log", nil)
	rr := httptest.NewRecorder()

	handler := Logging(logger)(nextHandler)
	handler.ServeHTTP(rr, req)

	assert.True(t, called)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}
