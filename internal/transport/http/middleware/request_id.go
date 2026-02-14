package middleware

import (
	"net/http"

	"github.com/google/uuid"
)

const HeaderXRequestID = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(HeaderXRequestID)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Прокидываем ID в заголовки ответа, чтобы клиент тоже его видел
			// (удобно для техподдержки: "у меня ошибка, вот мой request-id")
			w.Header().Set(HeaderXRequestID, requestID)

			// Обновляем заголовок в запросе, чтобы метод WithRequest его прочитал
			r.Header.Set(HeaderXRequestID, requestID)

			next.ServeHTTP(w, r)
		})
	}
}
