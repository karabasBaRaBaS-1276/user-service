package actuator

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/middleware"
	"go.uber.org/zap"
)

type LivenessResponse struct {
	Status    string `json:"status" example:"UP"`
	Timestamp string `json:"timestamp" example:"2026-02-09T15:04:05Z"`
}

// Liveness godoc
// @Summary      Проверка приложения на запуск
// @Description  Простая проверка работоспособности, чтобы узнать, запущено ли приложение
// @Tags         Actuator
// @Produce      json
// @Success      200  {object}  LivenessResponse
// @Router       /health/liveness [get]
func Liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := middleware.FromContext(r.Context())

		// Используем структуру вместо мапы
		resp := LivenessResponse{
			Status:    "UP",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Кодируем структуру напрямую
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			l.Error("failed to write liveness response", zap.Error(err))
		}
	}
}
