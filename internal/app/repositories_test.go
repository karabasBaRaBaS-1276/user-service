package app

import (
	"database/sql"
	"testing"

	"github.com/karabasBaRaBaS-1276/user-service/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRepositories(t *testing.T) {
	var db *sql.DB
	repos := initRepositories(db)

	require.NotNil(t, repos)

	t.Run("PersonRepository", func(t *testing.T) {
		assert.NotNil(t, repos.Person)
		assert.Implements(t, (*user.PersonRepository)(nil), repos.Person)
	})

	t.Run("UserCredentialRepository", func(t *testing.T) {
		assert.NotNil(t, repos.UserCredential)
		assert.Implements(t, (*user.UserCredentialRepository)(nil), repos.UserCredential)
	})
}
