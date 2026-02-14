package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/karabasBaRaBaS-1276/user-service/internal/domain/user"
)

type PersonRepository struct {
	db *sql.DB
}

var _ user.PersonRepository = (*PersonRepository)(nil)

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

// Create implements [user.PersonRepository].
func (p *PersonRepository) Create(ctx context.Context) error {
	return fmt.Errorf("метод Create не реализован")
}
