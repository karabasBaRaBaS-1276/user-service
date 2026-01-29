package app

import (
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	httptransport "github.com/karabasBaRaBaS-1276/user-service/internal/transport/http"
	"go.uber.org/zap"
)

func initHTTPServer(cfg *config.Config, logger *zap.Logger, services *Services) *http.Server {

	// Формируем ссылки на указатели с реализацией бизнес-слоя
	impl := httptransport.Deps{
		Person: services.Person,
	}
	return httptransport.NewServer(cfg, logger, impl)

}
