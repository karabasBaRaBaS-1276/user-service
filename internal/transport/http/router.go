package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/handler/actuator"
	userapi "github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/handler/user"
	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/middleware"
	"github.com/karabasBaRaBaS-1276/user-service/internal/transport/http/response"
	"go.uber.org/zap"

	_ "github.com/karabasBaRaBaS-1276/user-service/docs/actuator" // Путь к сгенерированному пакету
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(ctx context.Context, logger *zap.Logger, deps Deps, serviceName string, serviceVersion string) http.Handler {

	mux := http.NewServeMux()

	// 1. API для мониторинга и управления приложениями
	mux.Handle("GET /actuator/health/liveness", actuator.Liveness())
	mux.Handle("GET /actuator/health/readiness", actuator.Readiness(deps.DB))
	mux.Handle("GET /actuator/info", actuator.Info(deps.Config))

	apiPrefix := "/" + serviceName
	// Регистрация Swagger UI для actuator
	mux.Handle("GET "+apiPrefix+"/swagger/actuator/", httpSwagger.Handler(
		httpSwagger.InstanceName("actuator"), // Указываем наше имя Swagger инстанса
	))

	// 2. API с бизнес функционалом подключаем используя кодогенерированные файлы
	personHandler := userapi.NewPersonHandler(deps.Person)
	strictHandler := userapi.NewStrictHandler(personHandler, nil)

	// Создаем свою функцию-обработчик ошибок для ошибок валидации запросов
	customErrorHandler := func(w http.ResponseWriter, r *http.Request, err error) {

		// Достаем логгер из контекста
		reqLogger := middleware.FromContext(r.Context())

		// Логируем проблему как Warn, так как это ошибка клиента, а не сервера
		reqLogger.Sugar().Warnf("Ошибка на уровне валидации по OpenAPI схеме: %s", err.Error())

		errResp := response.NewError(
			http.StatusBadRequest,
			http.StatusText(http.StatusBadRequest),
			err.Error(),
			nil,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResp)
	}

	// Хендлер для отдачи вшитой спецификации openapi
	specURL := apiPrefix + "/swagger/openapi.json"
	mux.HandleFunc("GET "+specURL, func(w http.ResponseWriter, r *http.Request) {
		swagger, err := userapi.GetSwagger()
		if err != nil {
			logger.Error("Failed to get swagger spec", zap.Error(err))
			http.Error(w, "Failed to get spec", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(swagger)
	})

	// Регистрация Swagger UI для бизнес-API
	mux.Handle("GET "+apiPrefix+"/swagger/api/", httpSwagger.Handler(
		httpSwagger.InstanceName("api"), // Указываем наше имя Swagger инстанса
		httpSwagger.URL(specURL),        // Указываем путь к нашему JSON
	))

	// 2. Настраиваем опции для генератора
	options := userapi.StdHTTPServerOptions{
		BaseRouter:       mux,
		BaseURL:          apiPrefix + "/api/v1",
		ErrorHandlerFunc: customErrorHandler,
	}

	// 3. Вызываем HandlerWithOptions
	userapi.HandlerWithOptions(strictHandler, options)

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
