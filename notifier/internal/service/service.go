package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bopoh24/ma_1/notifier/internal/config"
	"github.com/bopoh24/ma_1/notifier/internal/model"
	mail "github.com/xhit/go-simple-mail/v2"
	"text/template"
	"time"
)

type Service struct {
	conf   config.SMTP
	client *mail.SMTPClient
}

// New returns a new Service instance
func New(conf config.SMTP) (*Service, error) {
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = conf.Host
	server.Port = conf.Port
	if conf.Username != "" && conf.Password != "" {
		server.Authentication = mail.AuthAuto
		server.Username = conf.Username
		server.Password = conf.Password
	}
	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// SMTP client
	smtpClient, err := server.Connect()

	if err != nil {
		return nil, err
	}
	return &Service{conf: conf, client: smtpClient}, nil
}

func (s *Service) sendEmail(to, subject, body string, cc ...string) error {
	email := mail.NewMSG()
	email.SetFrom(s.conf.From).AddTo(to).SetSubject(subject).SetBody(mail.TextHTML, body)
	if len(cc) > 0 {
		email.AddCc(cc...)
	}
	err := email.Send(s.client)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}

func (s *Service) BookingFailed(notification model.BookingNotification) error {
	notification.Status = "failed to pay"
	return s.sendNotifications(notification)
}

func (s *Service) BookingPaid(order model.BookingNotification) error {
	order.Status = "paid"
	return s.sendNotifications(order)
}

func (s *Service) BookingSubmitted(notification model.BookingNotification) error {
	notification.Status = "submitted"
	return s.sendNotifications(notification)
}

func (s *Service) BookingCompleted(notification model.BookingNotification) error {
	notification.Status = "completed"
	return s.sendNotifications(notification)
}

func (s *Service) BookingCancelledByCustomer(notification model.BookingNotification) error {
	notification.Status = "cancelled by customer"
	return s.sendNotifications(notification)
}

func (s *Service) BookingCancelledByCompany(notification model.BookingNotification) error {
	notification.Status = "cancelled by company"
	return s.sendNotifications(notification)
}

func (s *Service) tplToString(templateName string, booking model.BookingNotification) (string, error) {
	tpl, err := template.New(templateName).ParseFiles(templateName + ".gohtml")
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	if err := tpl.Execute(&output, booking); err != nil {
		return "", err
	}
	return output.String(), nil
}

func (s *Service) sendNotifications(order model.BookingNotification) error {
	cc := make([]string, 0, len(order.CompanyManagerContacts))
	for _, c := range order.CompanyManagerContacts {
		cc = append(cc, c.Email)
	}
	body, err := s.tplToString("company", order)
	if err != nil {
		return err
	}
	// send email to company managers
	err = s.sendEmail(order.CompanyContacts.Email, "Order paid", body, cc...)
	if err != nil {
		return err
	}

	body, err = s.tplToString("customer", order)
	if err != nil {
		return err
	}
	// send email to customer
	return s.sendEmail(order.CustomerContacts.Email, "Order paid", body)
}

// Close closes the Service
func (s *Service) Close(_ context.Context) error {
	return s.client.Quit()
}
