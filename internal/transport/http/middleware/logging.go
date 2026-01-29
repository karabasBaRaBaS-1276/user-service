package middleware

import (
	"context"
	"net/http"
	"time"

	logger "github.com/karabasBaRaBaS-1276/user-service/pkg/logger"
	"go.uber.org/zap"
)

type contextKey string

const loggerKey contextKey = "requestLogger"

// FromContext возвращает request-scoped логгер из context
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}

// Logging middleware логирует HTTP-запросы
func Logging(ml *logger.MiddlewareLogger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// request-scoped логгер с HTTP-полями
			reqLogger := ml.WithRequest(r)

			ctx := context.WithValue(r.Context(), loggerKey, reqLogger)

			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r.WithContext(ctx))

			reqLogger.Info(
				"http request completed",
				zap.Int("http_status", rw.status),
				zap.Duration("duration", time.Since(start)),
			)
		})
	}
}

// responseWriter перехватывает HTTP status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
