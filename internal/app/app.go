package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"go.uber.org/zap"
)

type App struct {
	cfg        *config.Config
	logger     *zap.Logger
	httpServer *http.Server
	db         *sql.DB
}

// Конструктор приложения с необходимыми для его работы ресурсами
func New(cfg *config.Config, log *zap.Logger) (*App, error) {
	app := &App{
		cfg:    cfg,
		logger: log,
	}

	if err := app.init(); err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	return app, nil
}

// Инициализация приложения
func (a *App) init() error {
	// 1. База данных
	db, err := initDB(a.cfg, a.logger)
	if err != nil {
		return err
	}
	a.db = db

	/*
		// 2. Миграции
		if a.cfg.DB.Migrate {
			if err := runMigrations(a.cfg); err != nil {
				return err
			}
		}

		// 3. Репозитории
		repos := initRepositories(db)

		// 4. Сервисы
		services := initServices(repos)

		// 5. HTTP сервер
		a.httpServer = initHTTPServer(a.cfg, a.logger, services)
	*/
	return nil
}

// Запуск приложения
func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Запуск приложения",
		zap.String("service", a.cfg.ServiceName),
		zap.String("environment", a.cfg.Environment),
	)

	errCh := make(chan error, 1)

	go func() {
		a.logger.Info("Запущен http server",
			zap.String("addr", a.httpServer.Addr),
		)

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		a.logger.Info("Получен сигнал на выключение приложения")
		return a.shutdown()

	case err := <-errCh:
		return fmt.Errorf("Ошибка http сервера: %w", err)
	}
}

// Остановка приложения
func (a *App) shutdown() error {
	timeout := a.cfg.Server.ShutdownTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	a.logger.Info("Выключение http сервера",
		zap.Duration("timeout", timeout),
	)

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("Ошибка выключения http сервера: %w", err)
	}

	a.logger.Info("Приложение безопасно остановлено")
	return nil
}
