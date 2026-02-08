package actuator

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/middleware"
	"go.uber.org/zap"
)

// DatabasePinger описывает методы БД, необходимые для проверки готовности
type DatabasePinger interface {
	PingContext(ctx context.Context) error
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// Приложение готово принимать трафик
func Readiness(db DatabasePinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.FromContext(r.Context())
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// Структура ответа
		status := "UP"
		checks := make(map[string]any)

		// 1. Проверка БД
		if err := db.PingContext(ctx); err != nil {
			status = "DOWN"
			checks["database"] = "unreachable"
			logger.Error("readiness: database unreachable", zap.Error(err))
		} else {
			// 2. Получаем версию миграции (прямым запросом, чтобы не плодить зависимости от golang-migrate в хендлере)
			var version int64
			var dirty bool
			err := db.QueryRowContext(ctx, "SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty)
			if err != nil {
				checks["database"] = "OK (migration info unavailable)"
			} else {
				migrationStatus := "clean"
				if dirty {
					migrationStatus = "dirty"
					status = "DOWN" // Если миграция в состоянии dirty, сервис не готов
				}
				checks["database"] = fmt.Sprintf("OK (version: %d, status: %s)", version, migrationStatus)
			}
		}

		// 3. Пример будущей проверки JWT ключей (внешних)
		// if !keysLoaded() {
		//    status = "DOWN"
		//    checks["jwt_keys"] = "not_loaded"
		// }

		w.Header().Set("Content-Type", "application/json")

		if status == "DOWN" {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		resp := map[string]any{
			"status": status,
			"checks": checks,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write readiness response", zap.Error(err))
		}
	}
}
