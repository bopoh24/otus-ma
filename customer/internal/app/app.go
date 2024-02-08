package app

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/repository/pg"
	"github.com/bopoh24/ma_1/customer/internal/service"
	"github.com/bopoh24/ma_1/pkg/http/router"
	"github.com/bopoh24/ma_1/pkg/verifier/phone"
	"github.com/go-chi/chi/v5"
	"log/slog"
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
	repo, err := pg.New(a.conf.Postgres)
	if err != nil {
		return err
	}

	// init phone verifier
	phoneVerifier, err := phone.NewStubPhoneVerify()
	if err != nil {
		return err
	}
	a.service = service.New(a.conf, repo, phoneVerifier)

	r := router.New("customer")
	r.Route("/customer", func(r chi.Router) {
		r.Post("/login", a.handlerLogin)
		r.Post("/logout", a.hanlderLogout)
		r.Post("/register", a.handlerRegister)
		r.Post("/refresh", a.handlerRefresh)

		r.Get("/profile", a.handlerProfile)
		r.Put("/profile", a.handlerProfileUpdate)
		r.Post("/phone/verify", a.handlerRequestPhoneVerification)
		r.Post("/phone/verify/check", a.handlerVerifyPhone)
		r.Get("/{id}", a.handlerCustomerByID)
	})
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
