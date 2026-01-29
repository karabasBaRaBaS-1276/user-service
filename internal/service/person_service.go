package service

import (
	"context"

	"github.com/karabasBaRaBaS-1276/user-service/internal/domain/user"
)

type personService struct {
	personRepo         user.PersonRepository
	userCredentialRepo user.UserCredentialRepository
}

var _ PersonService = (*personService)(nil)

func NewPersonService(personRepo user.PersonRepository, userCredentialRepo user.UserCredentialRepository) PersonService {
	return &personService{
		personRepo:         personRepo,
		userCredentialRepo: userCredentialRepo,
	}
}

func (s *personService) CreatePerson(ctx context.Context) error {
	// TODO: orchestration:
	// 1. validation
	// 2. create person
	// 3. create credentials
	// 4. transaction handling
	return nil
}
