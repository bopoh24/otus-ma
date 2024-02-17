package app

import (
	"context"
	"github.com/bopoh24/ma_1/notifier/internal/config"
	"github.com/bopoh24/ma_1/notifier/internal/service"
	"github.com/bopoh24/ma_1/pkg/http/router"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net"
	"net/http"
)

type App struct {
	conf    *config.Config
	service *service.Service
	log     *slog.Logger
	server  *http.Server
}

func New(conf *config.Config, log *slog.Logger) *App {
	return &App{
		log:    log,
		conf:   conf,
		server: &http.Server{Addr: ":80"},
	}
}

func (a *App) Run(ctx context.Context) error {
	var err error
	a.service, err = service.New(a.conf.SMTP)
	if err != nil {
		return err
	}
	r := router.New("notifier")
	r.Route("/notifier", func(r chi.Router) {
		r.Post("/send", a.handlerSend)
	})
	// set base context
	a.server.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}
	a.server.Handler = r
	return a.server.ListenAndServe()
}

func (a *App) Close(ctx context.Context) {
	if err := a.service.Close(ctx); err != nil {
		a.log.Error("Error closing service", "err", err)
	} else {
		a.log.Info("Service closed")
	}
	a.log.Info("App is closing")
	if err := a.server.Close(); err != nil {
		a.log.Error("Error closing server", "err", err)
	} else {
		a.log.Info("Server closed")
	}
}
