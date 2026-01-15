package user

import "context"

// Интерфейс по работе с хранилищем данных для авторизации пользователей
type UserCredentialRepository interface {
	//Create(ctx context.Context, p *Person) error
	Create(ctx context.Context) error
}
