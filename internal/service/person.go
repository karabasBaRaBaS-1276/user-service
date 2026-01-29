package service

import "context"

// PersonService описывает бизнес процессы, связанные с пользователями
type PersonService interface {
	// CreatePerson — создание пользователя (пока заглушка)
	CreatePerson(ctx context.Context) error
}
