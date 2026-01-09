package app

import (
	"testing"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/stretchr/testify/require"
)

func TestBuildPostgresDSN(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.DatabaseConfig
		want    string
		wantErr bool
	}{
		{
			name:    "error on empty config",
			cfg:     config.DatabaseConfig{},
			wantErr: true,
		},
		{
			name: "minimal valid config",
			cfg: config.DatabaseConfig{
				Host: "localhost",
				Port: "5432",
				User: "user",
				Name: "db",
			},
			want: "postgres://user@localhost:5432/db?client_encoding=UTF8",
		},
		{
			name: "with password",
			cfg: config.DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "user",
				Password: "secret",
				Name:     "db",
			},
			want: "postgres://user:secret@localhost:5432/db?client_encoding=UTF8",
		},
		{
			name: "with sslmode",
			cfg: config.DatabaseConfig{
				Host:    "localhost",
				Port:    "5432",
				User:    "user",
				Name:    "db",
				SSLMode: "require",
			},
			want: "postgres://user@localhost:5432/db?client_encoding=UTF8&sslmode=require",
		},
		{
			name: "with schema",
			cfg: config.DatabaseConfig{
				Host:   "localhost",
				Port:   "5432",
				User:   "user",
				Name:   "db",
				Schema: "user_service",
			},
			want: "postgres://user@localhost:5432/db?client_encoding=UTF8&search_path=user_service",
		},
		{
			name: "with sslmode and schema",
			cfg: config.DatabaseConfig{
				Host:    "localhost",
				Port:    "5432",
				User:    "user",
				Name:    "db",
				Schema:  "user_service",
				SSLMode: "verify-full",
			},
			want: "postgres://user@localhost:5432/db?client_encoding=UTF8&search_path=user_service&sslmode=verify-full",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn, err := buildPostgresDSN(tt.cfg)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, dsn)
		})
	}
}
