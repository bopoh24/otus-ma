package app

import (
	"github.com/bopoh24/ma_1/app/internal/config"
	"github.com/bopoh24/ma_1/app/internal/repository/memory"
	"github.com/bopoh24/ma_1/app/internal/repository/pg"
	"github.com/bopoh24/ma_1/app/internal/service/order"
	"github.com/bopoh24/ma_1/app/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
)

type App struct {
	userSrv  *user.Service
	orderSrv *order.Service
	cfg      *config.Config
	log      *slog.Logger
}

// New returns new app
func New(cfg *config.Config, log *slog.Logger) *App {
	repo, err := pg.New(cfg.Postgres)
	if err != nil {
		panic(err)
	}

	return &App{
		userSrv:  user.New(repo),
		orderSrv: order.New(memory.New()),
		cfg:      cfg,
		log:      log,
	}
}

func (a *App) Run() error {
	a.log.Info("App started", "name", a.cfg.App.Name)
	metricsMw := NewMetricsMiddleware(newMetrics())
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "ok"}`))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Use(metricsMw.Middleware)
		r.Post("/login", a.login)
		r.Post("/logout", a.logout)
		r.Post("/register", a.register)
		r.Post("/refresh", a.refresh)
	})

	r.Route("/user", func(r chi.Router) {
		r.Use(metricsMw.Middleware)
		r.Get("/me", a.userProfile)
		r.Put("/me", a.updateUserProfile)
		r.Post("/", a.userCreate)
		r.Get("/{id}", a.userByID)
		r.Put("/{id}", a.userUpdate)
		r.Delete("/{id}", a.userDelete)
	})

	r.Route("/order", func(r chi.Router) {
		r.Use(metricsMw.Middleware)
		r.Post("/", a.orderCreate)
		r.Get("/{id}", a.orderByID)
		r.Get("/user/{id}", a.orderByUserID)
		r.Put("/{id}", a.orderUpdate)
		r.Delete("/{id}", a.orderDelete)
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

	return http.ListenAndServe(":8000", r)
}
