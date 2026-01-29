package http

import (
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/handler"
	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/middleware"
	pkglogger "github.com/karabasBaRaBaS-1276/user-service/pkg/logger"
	"go.uber.org/zap"
)

func NewRouter(logger *zap.Logger, deps Deps, serviceName string, serviceVersion string) http.Handler {

	mwLogger := pkglogger.NewMiddlewareLogger(
		logger,
		serviceName,
		serviceVersion,
	)

	mux := http.NewServeMux()

	// API для мониторинга и управления приложениями
	//mux.HandleFunc("actuator/health", healthHandler)

	// API запросы по работе с пользователем
	personHandler := handler.NewPersonHandler(deps.Person)
	mux.Handle("/api/v1/user", personHandler)

	return applyGlobalMiddleware(
		mux,
		middleware.Recover(mwLogger),
		middleware.Logging(mwLogger),
	)
}

func applyGlobalMiddleware(h http.Handler, mw ...middleware.Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
