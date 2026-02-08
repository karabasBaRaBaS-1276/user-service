package middleware

import (
	"encoding/json"
	"net/http"
	"runtime/debug"

	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/response"
	"github.com/karabasBaRaBaS-1276/user-service/pkg/logger"
	"go.uber.org/zap"
)

func Recover(l *zap.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					// Используем нашу новую функцию для обогащения лога данными запроса
					reqLogger := logger.WithRequest(l, r)

					reqLogger.Error(
						"panic recovered",
						zap.Any("panic", rec),
						zap.ByteString("stacktrace", debug.Stack()),
					)

					errResp := response.NewError(
						http.StatusInternalServerError,
						http.StatusText(http.StatusInternalServerError),
						"Unexpected server error",
						nil,
					)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					_ = json.NewEncoder(w).Encode(errResp)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
