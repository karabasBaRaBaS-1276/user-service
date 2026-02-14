package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	// Пустой хендлер, который просто подтверждает, что до него дошел запрос
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("should generate new ID if not present", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		rr := httptest.NewRecorder()

		handler := RequestID()(nextHandler)
		handler.ServeHTTP(rr, req)

		// Проверяем заголовок в ответе
		respID := rr.Header().Get(HeaderXRequestID)
		assert.NotEmpty(t, respID, "RequestID должен быть сгенерирован")

		// Проверяем, что ID проброшен в сам запрос (для логгера)
		assert.Equal(t, respID, req.Header.Get(HeaderXRequestID))
	})

	t.Run("should preserve existing ID", func(t *testing.T) {
		existingID := "existing-uuid-123"
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set(HeaderXRequestID, existingID)
		rr := httptest.NewRecorder()

		handler := RequestID()(nextHandler)
		handler.ServeHTTP(rr, req)

		assert.Equal(t, existingID, rr.Header().Get(HeaderXRequestID), "Должен сохраниться старый ID")
	})
}
