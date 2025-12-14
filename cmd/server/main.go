package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"github.com/karabasBaRaBaS-1276/user-service/pkg/logger"
	"go.uber.org/zap"
)

// Эти переменные будут заполнены при сборке (build metadata)
var (
	version   = "dev" // будет перезаписана через ldflags
	commit    = "none"
	buildTime = "unknown"

	configPath  = flag.String("config", "", "Path to config file (default: use embedded config)")
	showVersion = flag.Bool("version", false, "Show version and exit")
	env         = flag.String("env", "", "Environment (development, production, staging)")
)

func main() {
	// Флаги запуска приложения
	flag.Parse()

	// Показать версию и выйти
	if *showVersion {
		fmt.Printf("User Service\n")
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Commit:     %s\n", commit)
		fmt.Printf("Build Time: %s\n", buildTime)
		os.Exit(0)
	}

	// Передаем переменные сборки в конфигурацию
	config.SetBuildInfo(version, commit, buildTime)

	// Определяем путь к конфигурационному файлу
	finalConfigPath := *configPath
	if finalConfigPath == "" && *env != "" {
		// Если указано окружение, ищем соответствующий файл
		finalConfigPath = findConfigByEnv(*env)
	}

	// Загрузка конфигурации
	cfg, err := config.Load(finalConfigPath, *env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	logger, err := logger.New(logger.Config(cfg.Logging))
	if err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// todo. Убрать, чтобы не показывать секреты
	logger.Debug("Конфигурация запуска",
		zap.Any("config_summary", cfg),
	)

	// Создание и запуск приложения
	/*
		application, err := app.New(cfg, logger)
		if err != nil {
			logger.Fatal("Failed to create app", err)
		}

		if err := application.Run(); err != nil {
			logger.Error("Application error", err)
		}
	*/

}

// findConfigByEnv ищет конфигурационный файл по имени окружения
func findConfigByEnv(env string) string {
	// Проверяем стандартные места
	possiblePaths := []string{
		filepath.Join("config", env+".yaml"),
		filepath.Join("config", env+".yml"),
		filepath.Join("config", env+".json"),
		filepath.Join("configs", env+".yaml"),
		env + ".yaml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return "" // Используем встроенную (embedded) конфигурацию
}
