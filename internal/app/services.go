package app

import "github.com/karabasBaRaBaS-1276/user-service/internal/service"

// Набор бизнес-сервисов приложения
type Services struct {
	// Сервисы по работе с пользователем
	Person service.PersonService
}

func initServices(repos *Repositories) *Services {
	return &Services{
		Person: service.NewPersonService(
			repos.Person,
			repos.UserCredential,
		),
	}
}
