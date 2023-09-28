package main

import (
	"github.com/bopoh24/ma_1/internal/config"
	"github.com/bopoh24/ma_1/internal/repository/pg"
	"github.com/bopoh24/ma_1/internal/service"
	"log/slog"
	"os"
)

func main() {
	// init logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// init config
	cfg, err := config.New()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("App started")
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
