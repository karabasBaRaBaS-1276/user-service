package http

import (
	"context"
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/handler"
	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/handler/actuator"
	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/middleware"
	"go.uber.org/zap"

	_ "github.com/karabasBaRaBaS-1276/user-service/docs/actuator" // Путь к сгенерированному пакету
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(ctx context.Context, logger *zap.Logger, deps Deps, serviceName string, serviceVersion string) http.Handler {

	mux := http.NewServeMux()

	// API для мониторинга и управления приложениями
	mux.Handle("GET /actuator/health/liveness", actuator.Liveness())
	mux.Handle("GET /actuator/health/readiness", actuator.Readiness(deps.DB))
	mux.Handle("GET /actuator/info", actuator.Info(deps.Config))

	// Регистрация Swagger UI
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.InstanceName("actuator"), // Указываем наше имя Swagger инстанса
	))

	// API запросы по работе с пользователем
	personHandler := handler.NewPersonHandler(deps.Person)
	mux.Handle("/api/v1/user", personHandler)

	return applyGlobalMiddleware(
		mux,
		middleware.Recover(logger),
		middleware.RequestID(),
		middleware.Logging(logger),
	)
}

func applyGlobalMiddleware(h http.Handler, mw ...middleware.Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
