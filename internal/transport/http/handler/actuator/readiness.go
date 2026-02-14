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

// ReadinessCheck описывает состояние конкретного компонента (БД, очереди и т.д.)
type ReadinessCheck struct {
	Status  string `json:"status" example:"UP"`
	Message string `json:"message,omitempty" example:"connected"`
}

// ReadinessResponse — основной объект ответа для Readiness probe
type ReadinessResponse struct {
	Status string                    `json:"status" example:"UP"`
	Checks map[string]ReadinessCheck `json:"checks"`
}

// Приложение готово принимать трафик
// Readiness godoc
// @Summary      Проверка приложения к приему трафика
// @Description  Простая проверка работоспособности, чтобы узнать, готово ли приложение принимать запросы на обработку
// @Tags         Actuator
// @Produce      json
// @Success      200  {object}  ReadinessResponse
// @Failure      503  {object}  ReadinessResponse
// @Router       /health/readiness [get]
func Readiness(db DatabasePinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := middleware.FromContext(r.Context())
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		checks := make(map[string]ReadinessCheck)
		overallStatus := "UP"

		// 1. Проверка БД
		if err := db.PingContext(ctx); err != nil {
			overallStatus = "DOWN"
			checks["database"] = ReadinessCheck{Status: "DOWN", Message: err.Error()}
		} else {
			checks["database"] = ReadinessCheck{Status: "UP", Message: "connected"}
			// 2. Получаем версию миграции (прямым запросом, чтобы не плодить зависимости от golang-migrate в хендлере)
			var version int64
			var dirty bool
			err := db.QueryRowContext(ctx, "SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty)
			if err != nil {
				checks["database"] = ReadinessCheck{Status: "UP", Message: "migration info unavailable"}
			} else {
				migrationStatus := "завершена"
				if dirty {
					migrationStatus = "прервалась (dirty)"
					overallStatus = "DOWN" // Если миграция в состоянии dirty, сервис не готов принимать трафик
				}
				checks["database"] = ReadinessCheck{Status: "UP", Message: fmt.Sprintf("Миграция %s. Версия: %d", migrationStatus, version)}
			}
		}

		// 3. Пример будущей проверки JWT ключей (внешних)
		// if !keysLoaded() {
		//    status = "DOWN"
		//    checks["jwt_keys"] = "not_loaded"
		// }

		resp := ReadinessResponse{
			Status: overallStatus,
			Checks: checks,
		}

		w.Header().Set("Content-Type", "application/json")

		if overallStatus == "DOWN" {
			w.WriteHeader(http.StatusServiceUnavailable) // 503 если не готовы
		} else {
			w.WriteHeader(http.StatusOK)
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write readiness response", zap.Error(err))
		}
	}
}
