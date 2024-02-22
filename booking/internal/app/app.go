package app

import (
	"context"
	"github.com/bopoh24/ma_1/booking/internal/config"
	"github.com/bopoh24/ma_1/booking/internal/repository"
	"github.com/bopoh24/ma_1/booking/internal/service"
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
	// init repository
	repo, err := repository.New(a.conf.Postgres)
	if err != nil {
		return err
	}
	a.service = service.New(repo)
	r := router.New("booking")
	r.Route("/booking", func(r chi.Router) {

		r.Get("/services", a.handlerGetServices)
		r.Post("/services", a.handlerAddService)
		// add offer
		r.Post("/offers", a.handlerAddOffer)
		// delete offer
		r.Delete("/offers/{id}", a.handlerDeleteOffer)
		// change offer status
		r.Put("/offers/{id}/status", a.handlerChangeOfferStatus)

		r.Post("/offers/{id}/cancel", a.handlerCancelOfferByCompany)
		r.Post("/offers/{id}/cancel/customer", a.handlerCancelOfferByCustomer)

		// get company offers
		r.Get("/company/offers/{id}", a.handlerGetCompanyOffers)
		// get customer offers
		r.Get("/customer/offers", a.handlerGetCustomerOffers)

		// search offers
		r.Get("/offers", a.handlerSearchOffers)

		// book offer
		r.Post("/offers/{id}/book", a.handlerBookOffer)
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
