package app

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/company/internal/config"
	"github.com/bopoh24/ma_1/company/internal/repository"
	"github.com/bopoh24/ma_1/company/internal/service"
	"github.com/bopoh24/ma_1/pkg/http/router"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net"
	"net/http"
)

type App struct {
	conf           *config.Config
	keycloakClient *gocloak.GoCloak
	service        *service.Service
	log            *slog.Logger
	server         *http.Server
}

func New(conf *config.Config, log *slog.Logger) *App {
	return &App{
		log:            log,
		conf:           conf,
		keycloakClient: gocloak.NewClient(conf.Keycloak.URL),
		server:         &http.Server{Addr: ":80"},
	}
}

func (a *App) Run(ctx context.Context) error {
	// init repository
	repo, err := repository.New(a.conf.Postgres)
	if err != nil {
		return err
	}
	a.service = service.New(repo)
	r := router.New("company")
	r.Route("/company", func(r chi.Router) {
		r.Post("/login", a.handlerLogin)
		r.Post("/logout", a.handlerLogout)
		r.Post("/register", a.handlerRegister)
		r.Post("/refresh", a.handlerRefresh)

		r.Get("/{id}", a.handlerCompanyDetails)
		r.Put("/{id}", a.handlerUpdateCompany)
		r.Post("/{id}/logo", a.handlerUpdateLogo)
		r.Post("/{id}/location", a.handlerUpdateLocation)
		r.Post("/{id}/activate", a.handlerActivateDeactivate(true))
		r.Post("/{id}/deactivate", a.handlerActivateDeactivate(false))
		r.Get("/{id}/managers", a.handlerGetManagers)
		r.Post("/", a.handlerCreateCompany)

	})
	// set base context
	a.server.BaseContext = func(listener net.Listener) context.Context {
		return ctx
	}
	a.server.Handler = r
	return a.server.ListenAndServe()
}

func (a *App) Close(ctx context.Context) {
	a.log.Info("App is closing")
	if err := a.server.Close(); err != nil {
		a.log.Error("Error closing server", "err", err)
	} else {
		a.log.Info("Server closed")
	}
}
