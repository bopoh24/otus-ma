package main

import (
	"github.com/bopoh24/ma_1/internal/config"
	"github.com/bopoh24/ma_1/internal/repository/pg"
	"github.com/bopoh24/ma_1/internal/service"
	"log/slog"
	"os"
)

func initLogLevel(level string) slog.Level {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	default:
		logLevel = slog.LevelError
	}
	return logLevel
}

func main() {
	// init config
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	// init logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: initLogLevel(cfg.App.LogLevel),
	}))

	logger.Info("App started", slog.With("name", cfg.App.Name))
	// init repository
	repo, err := pg.New(cfg.Postgres)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	srv := service.NewUserService(cfg, repo)
	if err := srv.Run(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
