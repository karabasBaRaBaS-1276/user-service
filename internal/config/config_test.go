package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("APP_NAME", "")

	cfg, err := Load("", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg == nil {
		t.Fatal("config must not be nil")
	}
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("APP_NAME", "unit-test-service")

	cfg, err := Load("", "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServiceName != "unit-test-service" {
		t.Errorf(
			"expected App.Name from env, got %q",
			cfg.ServiceName,
		)
	}
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config")

	content := []byte(`
app:
  name: file-service
`)

	err := os.WriteFile(configPath+".test.yaml", content, 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	cfg, err := Load(configPath, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServiceName != "file-service" {
		t.Errorf("expected name from file, got %s", cfg.ServiceName)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/no/such/path/config", "test")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
