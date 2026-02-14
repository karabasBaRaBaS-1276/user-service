package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/karabasBaRaBaS-1276/user-service/internal/domain/user"
)

type UserCredentialRepository struct {
	db *sql.DB
}

var _ user.UserCredentialRepository = (*UserCredentialRepository)(nil)

func NewUserCredentialRepository(db *sql.DB) *UserCredentialRepository {
	return &UserCredentialRepository{db: db}
}

// Create implements [user.UserCredentialRepository].
func (u *UserCredentialRepository) Create(ctx context.Context) error {
	return fmt.Errorf("метод Create не реализован")
}
