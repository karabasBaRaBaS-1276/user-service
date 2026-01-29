package http

import "github.com/karabasBaRaBaS-1276/user-service/internal/service"

// Deps описывает зависимости HTTP-транспортного слоя.
// Содержит только те интерфейсы, которые требуются обработчикам запросов.
type Deps struct {
	Person service.PersonService
}
