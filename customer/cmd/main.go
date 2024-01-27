package main

import (
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/repository/pg"
	"github.com/bopoh24/ma_1/customer/internal/service"
	"github.com/bopoh24/ma_1/pkg/logger"
	"log/slog"
	"os"
)

func main() {
	// init config
	cfg, err := config.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	// init logger
	log := logger.New(cfg.App.LogLevel)

	log.Info("App started", cfg.App.Name)
	// init repository
	repo, err := pg.New(cfg.Postgres)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	srv := service.NewUserService(cfg, repo)
	if err := srv.Run(); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}
