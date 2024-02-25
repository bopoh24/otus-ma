package app

import (
	"context"
	"encoding/json"
	"github.com/bopoh24/ma_1/notifier/internal/config"
	"github.com/bopoh24/ma_1/notifier/internal/service"
	"github.com/bopoh24/ma_1/notifier/pkg/model"
	"github.com/bopoh24/ma_1/pkg/http/router"
	"github.com/bopoh24/ma_1/pkg/kafka/consumer"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
)

const topic = "booking_notification"

type App struct {
	conf          *config.Config
	service       *service.Service
	log           *slog.Logger
	server        *http.Server
	kafkaConsumer *consumer.KafkaConsumer
}

func New(conf *config.Config, log *slog.Logger) *App {

	cons, err := consumer.NewKafkaConsumer(strings.Split(conf.Kafka.Hosts, ","), log)
	if err != nil {
		log.Error("Failed to create kafka consumer", "err", err)
		os.Exit(1)
	}
	return &App{
		log:           log,
		conf:          conf,
		server:        &http.Server{Addr: ":80"},
		kafkaConsumer: cons,
	}
}

func (a *App) Run(ctx context.Context) error {
	var err error
	a.service, err = service.New(a.conf.SMTP)
	if err != nil {
		return err
	}

	go func() {
		msgChan, err := a.kafkaConsumer.MessageChannel(topic)
		if err != nil {
			a.log.Error("Failed to create message channel", "err", err)
			return
		}
		for {
			if err := ctx.Err(); err != nil {
				a.log.Info("Context is done", "err", err)
				return
			}
			msg := <-msgChan

			notification := model.BookingNotification{}
			err = json.Unmarshal([]byte(msg), &notification)
			if err != nil {
				a.log.Error("Failed to unmarshal message", "err", err)
				continue
			}

			switch notification.Type {
			case model.BookingPaid:
				err = a.service.BookingPaid(notification)
			case model.BookingFailed:
				err = a.service.BookingFailed(notification)
			case model.BookingSubmitted:
				err = a.service.BookingSubmitted(notification)
			case model.BookingCompleted:
				err = a.service.BookingCompleted(notification)
			case model.BookingCancelledByCustomer:
				err = a.service.BookingCancelledByCustomer(notification)
			case model.BookingCancelledByCompany:
				err = a.service.BookingCancelledByCompany(notification)
			}
			if err != nil {
				a.log.Error("Failed to send email", "err", err)
			} else {
				a.log.Info("Email sent", "notification", notification)
			}
		}
	}()
	r := router.New("notifier")
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
	if err := a.kafkaConsumer.Close(); err != nil {
		a.log.Error("Error closing kafka consumer", "err", err)
	} else {
		a.log.Info("Kafka consumer closed")
	}
}
