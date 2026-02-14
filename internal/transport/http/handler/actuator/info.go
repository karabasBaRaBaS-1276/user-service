package actuator

import (
	"encoding/json"
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/middleware" // Путь к твоему FromContext
	"go.uber.org/zap"
)

type InfoResponse struct {
	Service     string `json:"service" example:"user-service"`
	Version     string `json:"version" example:"1.0.5"`
	Environment string `json:"environment" example:"production"`
}

// Info godoc
// @Summary      Информация о приложении
// @Description  Информация о работающем приложении
// @Tags         Actuator
// @Produce      json
// @Success      200  {object}  InfoResponse
// @Router       /info [get]
func Info(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger := middleware.FromContext(r.Context())

		infoResponse := InfoResponse{
			Service:     cfg.Logging.ServiceName,
			Version:     cfg.Logging.ServiceVersion,
			Environment: cfg.Environment,
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(infoResponse); err != nil {
			logger.Error("failed to encode info response", zap.Error(err))
		}
	}
}
