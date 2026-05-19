// Package app build all dependencies and provide methods for run and shutdown application
package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"todo/internal/api"
	"todo/internal/api/router"
	"todo/internal/config"
	"todo/internal/repository/db"
	"todo/internal/usecase/task"

	"github.com/76Parker/golib/httplib"
	"github.com/76Parker/golib/loglib"
	"github.com/76Parker/golib/pglib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

// App build app and run app
type App struct {
	s         *http.Server
	l         *loglib.Slog
	closeOnce sync.Once
}

// New is a constructor for App
func New(ctx context.Context, cfg config.Config) (*App, error) {
	log, err := initLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}
	dbPool, err := initPostgres(ctx, cfg.Postgres)
	if err != nil {
		return nil, err
	}

	if err := pglib.ApplyMigrations(dbPool, "./migrations", cfg.Postgres.Database); err != nil {
		return nil, err
	}

	pgAddress := net.JoinHostPort(cfg.Postgres.Host, strconv.Itoa(cfg.Postgres.Port))
	log.Info("Connected to PostgreSQL", "address", pgAddress, "database", cfg.Postgres.Database)

	taskRepo := db.NewTaskRepository(dbPool)
	taskSvc := task.NewService(taskRepo)
	taskHandler := api.NewTaskHandler(taskSvc)

	handlers := router.Handlers{
		Tasks: taskHandler,
	}

	serverHandler := initRouter(handlers, log)
	serverHandlerWithDefaultCORS := cors.Default().Handler(serverHandler)

	server := httplib.NewHTTPServer(ctx, cfg.HTTP, serverHandlerWithDefaultCORS)

	return &App{
		s:         server,
		l:         log,
		closeOnce: sync.Once{},
	}, nil
}

// Start run HTTP server
func (a *App) Start() error {
	a.l.Info("Starting HTTP server", "address", a.s.Addr)
	swaggerMsg := fmt.Sprintf("Swagger available at http://%s/swagger", a.s.Addr)
	a.l.Info(swaggerMsg)
	if err := a.s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown implement graceful shutdown for app
func (a *App) Shutdown(ctx context.Context) error {
	a.l.Info("Shutting down HTTP server...")
	a.closeOnce.Do(func() {
		_ = a.l.Close()
	})
	return a.s.Shutdown(ctx)
}

func initLogger(cfg loglib.SlogConfig) (*loglib.Slog, error) {
	return loglib.NewSlog(cfg)
}

func initPostgres(ctx context.Context, cfg pglib.Config) (*pgxpool.Pool, error) {
	return pglib.New(ctx, &cfg)
}

func initRouter(handlers router.Handlers, log loglib.Logger) http.Handler {
	r := router.New(handlers, log)
	return r
}
