package app

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/repository/pg"
	"github.com/bopoh24/ma_1/customer/internal/service"
	appMiddleware "github.com/bopoh24/ma_1/pkg/http/middleware"
	"github.com/bopoh24/ma_1/pkg/verifier/phone"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	mw := appMiddleware.NewMetricsMiddleware("customer")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Use(mw.Middleware)
		r.Post("/login", a.login)
		r.Post("/logout", a.logout)
		r.Post("/register", a.register)
		r.Post("/refresh", a.refresh)
	})

	r.Route("/customer", func(r chi.Router) {
		r.Use(mw.Middleware)
		r.Get("/profile", a.customerProfile)
		r.Put("/profile", a.updateCustomerProfile)
		r.Post("/phone/verify", a.requestPhoneVerification)
		r.Post("/phone/verify/check", a.verifyPhone)
		r.Get("/{id}", a.customerByID)
	})

	// metrics handler
	r.Handle("/metrics", promhttp.Handler())

	// Readiness and liveness probes
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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
