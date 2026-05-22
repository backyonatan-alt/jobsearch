package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/backyonatan-alt/jobsearch/internal/auth"
	"github.com/backyonatan-alt/jobsearch/internal/config"
	"github.com/backyonatan-alt/jobsearch/internal/db"
	"github.com/backyonatan-alt/jobsearch/internal/httpsrv"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.FromEnv()
	if err != nil {
		logger.Error("config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("db connect", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	migrationsDir := getenv("MIGRATIONS_DIR", "migrations")
	if err := db.Migrate(ctx, pool, migrationsDir); err != nil {
		logger.Error("migrate", "err", err)
		os.Exit(1)
	}
	logger.Info("migrations applied", "dir", migrationsDir)

	staticDir := getenv("STATIC_DIR", "web/build")
	abs, _ := filepath.Abs(staticDir)
	logger.Info("serving static", "dir", abs)

	srv := &httpsrv.Server{
		Cfg:    cfg,
		Pool:   pool,
		Auth:   auth.NewService(pool, cfg.SessionTTL),
		Google: auth.NewGoogle(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL()),
		Logger: logger,
		Static: http.Dir(staticDir),
	}

	httpSrv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           srv.Routes(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		logger.Info("listening", "port", cfg.Port, "brand", cfg.BrandName, "base_url", cfg.BaseURL)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
}

func getenv(k, def string) string {
	if v, ok := os.LookupEnv(k); ok && v != "" {
		return v
	}
	return def
}
