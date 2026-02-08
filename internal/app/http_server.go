package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	httptransport "github.com/karabasBaRaBaS-1276/user-service/internal/transport/http"
	"go.uber.org/zap"
)

func initHTTPServer(ctx context.Context, cfg *config.Config, logger *zap.Logger, services *Services, db *sql.DB) *http.Server {

	// Формируем ссылки на указатели с реализацией бизнес-слоя
	impl := httptransport.Deps{
		Person: services.Person,
		Config: cfg,
		DB:     db,
	}
	return httptransport.NewServer(ctx, cfg, logger, impl)

}
