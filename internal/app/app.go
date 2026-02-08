package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/karabasBaRaBaS-1276/user-service/internal/config"
	"go.uber.org/zap"
)

// App отвечает за инициализацию и жизненный цикл приложения.
// Здесь происходит сборка инфраструктуры и связывание компонентов.
type App struct {
	cfg        *config.Config
	logger     *zap.Logger
	httpServer *http.Server
	db         *sql.DB
}

// Конструктор приложения с необходимыми для его работы ресурсами
func New(ctx context.Context, cfg *config.Config, log *zap.Logger) (*App, error) {
	app := &App{
		cfg:    cfg,
		logger: log,
	}

	if err := app.init(ctx); err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	return app, nil
}

// Инициализация приложения
func (a *App) init(ctx context.Context) error {
	// 1. База данных
	db, err := initDB(ctx, a.cfg, a.logger)
	if err != nil {
		return err
	}
	a.db = db

	// 2. Миграции
	if err := runMigrations(db, "file://migrations", a.logger); err != nil {
		return err
	}
	// 3. Репозитории
	repos := initRepositories(db)

	// 4. Сервисы
	services := initServices(repos)

	// 5. HTTP сервер
	a.httpServer = initHTTPServer(ctx, a.cfg, a.logger, services, db)

	return nil
}

// Запуск приложения
func (a *App) Run(ctx context.Context) error {

	if a.httpServer == nil {
		return errors.New("http server is not initialized")
	}

	logger := a.logger

	serverErr := make(chan error, 1)

	go func() {
		logger.Sugar().Infof("Запуск приложения по адресу: %s", a.httpServer.Addr)

		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info("Контекст завершён, начинаем graceful shutdown")
		return a.shutdown(ctx)

	case err := <-serverErr:
		return fmt.Errorf("http server startup: %w", err)
	}
}

// Остановка приложения
func (a *App) shutdown(parentCtx context.Context) error {
	timeout := a.cfg.Server.ShutdownTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()

	a.logger.Info("Выключение http сервера",
		zap.Duration("timeout", timeout),
	)

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	if a.db != nil {
		a.logger.Info("Закрытие соединения с БД")
		if err := a.db.Close(); err != nil {
			a.logger.Warn("Ошибка при закрытии БД", zap.Error(err))
		}
	}

	a.logger.Info("Приложение безопасно остановлено")
	return nil
}
