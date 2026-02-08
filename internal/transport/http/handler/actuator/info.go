package actuator

import (
	"encoding/json"
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/middleware" // Путь к твоему FromContext
	"go.uber.org/zap"
)

func Info(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		logger := middleware.FromContext(r.Context())

		resp := map[string]any{
			"service":     cfg.Logging.ServiceName,
			"version":     cfg.Logging.ServiceVersion,
			"environment": cfg.Environment,
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to encode info response", zap.Error(err))
		}
	}
}
