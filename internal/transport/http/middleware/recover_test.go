package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRecover(t *testing.T) {
	// Создаем логгер-заглушку, чтобы не мусорить в консоль при тестах
	logger := zap.NewNop()

	// Хендлер-камикадзе
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom!")
	})

	t.Run("should recover from panic and return 500", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/panic", nil)
		rr := httptest.NewRecorder()

		// Оборачиваем паникующий хендлер в наш Recover
		handler := Recover(logger)(panicHandler)

		// Выполняем. Если Recover не сработает, тест упадет с грохотом
		assert.NotPanics(t, func() {
			handler.ServeHTTP(rr, req)
		})

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Unexpected server error")
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	})
}
