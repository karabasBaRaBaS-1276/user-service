package http

import (
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"go.uber.org/zap"
)

func NewServer(cfg *config.Config, logger *zap.Logger, deps Deps) *http.Server {

	handler := NewRouter(logger, deps, cfg.ServiceName, cfg.Logging.ServiceVersion)

	return &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}
