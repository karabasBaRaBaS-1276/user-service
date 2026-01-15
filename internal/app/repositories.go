package app

import (
	"database/sql"

	"github.com/karabasBaRaBaS-1276/user-service/internal/domain/user"
	"github.com/karabasBaRaBaS-1276/user-service/internal/repository/postgres"
)

// Список хранилищ приложения
type Repositories struct {
	// Хранилище для персональных данных пользователя
	Person user.PersonRepository
	// Хранилище для авторизации пользователей
	UserCredential user.UserCredentialRepository
}

func initRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Person:         postgres.NewPersonRepository(db),
		UserCredential: postgres.NewUserCredentialRepository(db),
		/*
			AuthClient:     postgres.NewAuthClientRepository(db),
			TempAuth:       postgres.NewTempAuthRepository(db),
			PublicKey:      postgres.NewPublicKeyRepository(db),
		*/
	}
}
