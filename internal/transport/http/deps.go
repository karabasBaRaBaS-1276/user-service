package http

import (
	"database/sql"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/karabasBaRaBaS-1276/user-service/internal/service"
)

// Deps описывает зависимости HTTP-транспортного слоя.
// Содержит только те интерфейсы и структуры, которые требуются обработчикам запросов.
type Deps struct {
	Person service.PersonService

	// actuator
	DB     *sql.DB
	Config *config.Config
}
