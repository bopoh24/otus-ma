package app

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/bopoh24/ma_1/customer/internal/config"
	"github.com/bopoh24/ma_1/customer/internal/repository"
	"github.com/bopoh24/ma_1/customer/internal/service"
	"github.com/bopoh24/ma_1/pkg/http/client"
	"github.com/bopoh24/ma_1/pkg/http/router"
	"github.com/bopoh24/ma_1/pkg/kafka/producer"
	"github.com/bopoh24/ma_1/pkg/verifier/phone"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net"
	"net/http"
	"strings"
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

	// init phone verifier
	phoneVerifier := phone.NewStubPhoneVerify()
	if err != nil {
		return err
	}

	prod, err := producer.NewKafkaProducer(strings.Split(a.conf.Kafka.Hosts, ","), a.log)
	if err != nil {
		return err
	}

	a.service = service.New(repo, prod,
		phoneVerifier,
		client.NewHttpClient(a.conf.CompanyUrl),
		client.NewHttpClient(a.conf.BookingUrl),
		client.NewHttpClient(a.conf.PaymentUrl))

	r := router.New("customer")
	r.Route("/customer", func(r chi.Router) {
		// auth
		r.Post("/login", a.handlerLogin)
		r.Post("/logout", a.hanlderLogout)
		r.Post("/register", a.handlerRegister)
		r.Post("/refresh", a.handlerRefresh)

		// profile
		r.Get("/profile", a.handlerProfile)
		r.Put("/profile", a.handlerProfileUpdate)
		r.Post("/phone/verify", a.handlerRequestPhoneVerification)
		r.Post("/phone/verify/check", a.handlerVerifyPhone)
		r.Get("/{id}", a.handlerCustomerByID)

		// booking
		r.Post("/offer/{id}/book", a.handlerBookOffer)

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
