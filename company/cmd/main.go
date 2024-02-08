package main

import (
	"context"
	"github.com/bopoh24/ma_1/company/internal/app"
	"github.com/bopoh24/ma_1/company/internal/config"
	"github.com/bopoh24/ma_1/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"time"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// init app
	a := app.New(cfg, log)
	if err := a.Run(ctx); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	// graceful shutdown
	<-ctx.Done()
	// closing context
	closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer closeCancel()
	a.Close(closeCtx)
}
