//go:build integration
// +build integration

package app

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var migrateOnce sync.Once

const migrationsPath = "file://../../migrations"

func runMigrationsOnce(t *testing.T) {
	t.Helper()

	migrateOnce.Do(func() {
		require.NotNil(t, testDB, "testDB must be initialized in TestMain")

		err := runMigrations(testDB, migrationsPath, zap.NewNop())
		require.NoError(t, err)
	})
}

func TestMigrations_ApplyAll(t *testing.T) {
	runMigrationsOnce(t)
}

func TestMigrations_TablesExist(t *testing.T) {
	runMigrationsOnce(t)

	tables := []string{
		"person",
		"user_credential",
		"gloss_auth_client",
		"temp_auth_init",
		"public_key",
	}

	for _, table := range tables {
		var exists bool
		err := testDB.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = current_schema()
				  AND table_name = $1
			)
		`, table).Scan(&exists)

		require.NoError(t, err)
		require.True(t, exists, "table %s should exist", table)
	}
}

func TestMigrations_SeedGlossAuthClient(t *testing.T) {
	runMigrationsOnce(t)

	var count int
	err := testDB.QueryRow(`
		SELECT count(*)
		FROM gloss_auth_client
		WHERE id = 'myBank'
	`).Scan(&count)

	require.NoError(t, err)
	require.Equal(t, 1, count)
}

func TestTrigger_PreventPersonIDUpdate(t *testing.T) {
	runMigrationsOnce(t)

	_, err := testDB.Exec(`
		INSERT INTO temp_auth_init (
			id, person_id, response_type, client_id,
			redirect_uri, scope, state,
			code_challenge_method, code_challenge,
			create_time_ms, last_modify_time_ms
		) VALUES (
			'test1', 'person-1', 'code', 'myBank',
			'/cb', 'dbo', 'state',
			'S256', 'challenge',
			1, 1
		)
	`)
	require.NoError(t, err)

	_, err = testDB.Exec(`
		UPDATE temp_auth_init
		SET person_id = 'person-2'
		WHERE id = 'test1'
	`)
	require.Error(t, err, "updating person_id must be forbidden by trigger")
}

func TestMigrations_DownAll(t *testing.T) {
	runMigrationsOnce(t)

	m, err := newMigrator(testDB, migrationsPath, zap.NewNop())
	require.NoError(t, err)

	err = m.Down()
	require.NoError(t, err)

	// Проверяем, что бизнес-таблицы удалены
	tables := []string{
		"person",
		"user_credential",
		"gloss_auth_client",
		"temp_auth_init",
		"public_key",
	}

	for _, table := range tables {
		var exists bool
		err := testDB.QueryRow(`
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = current_schema()
				  AND table_name = $1
			)
		`, table).Scan(&exists)

		require.NoError(t, err)
		require.False(t, exists, "Таблица %s должна быть удалена при down миграции", table)
	}
}

func TestMigrations_RollbackToVersion7(t *testing.T) {
	runMigrationsOnce(t)

	m, err := newMigrator(testDB, migrationsPath, zap.NewNop())
	require.NoError(t, err)

	// 1. Откатываемся конкретно на версию 7. Так как 8 миграция создает функцию и триггер
	err = m.Migrate(7)
	require.NoError(t, err)

	// 2. Проверяем состояние в таблице schema_migrations
	var version int
	var dirty bool
	err = testDB.QueryRow(`SELECT version, dirty FROM schema_migrations`).Scan(&version, &dirty)
	require.NoError(t, err)
	assert.Equal(t, 7, version, "Текущая версия миграции должна быть 7")
	assert.False(t, dirty, "База данных не должна быть в состоянии dirty")

	// 3. Проверяем удаление триггера через information_schema.triggers
	var triggerExists bool
	err = testDB.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.triggers 
            WHERE event_object_schema = 'user_service' 
              AND event_object_table = 'temp_auth_init'
              AND trigger_name = 'prevent_person_id_update_trigger'
        )
    `).Scan(&triggerExists)
	require.NoError(t, err)
	assert.False(t, triggerExists, "Триггер должен быть удален из таблицы temp_auth_init")

	// 4. Проверяем удаление функции через information_schema.routines
	var funcExists bool
	err = testDB.QueryRow(`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.routines 
            WHERE routine_schema = 'user_service' 
              AND routine_name = 'check_person_id'
        )
    `).Scan(&funcExists)
	require.NoError(t, err)
	assert.False(t, funcExists, "Функция check_person_id должна быть удалена")

}
