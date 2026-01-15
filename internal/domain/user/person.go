package user

import "context"

// Интерфейс по работе с хранилищем для персональных данных пользователя
type PersonRepository interface {
	//Create(ctx context.Context, p *Person) error
	Create(ctx context.Context) error
}
