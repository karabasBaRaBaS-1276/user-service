package api

// Генерируем модели (Domain)
//go:generate go tool oapi-codegen -config ../internal/transport/http/handler/user/config_models.yaml openapi.yaml

// Генерируем сервер (Transport) с привязкой к моделям
//go:generate go tool oapi-codegen -config ../internal/transport/http/handler/user/config_server.yaml openapi.yaml
