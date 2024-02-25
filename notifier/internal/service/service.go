package service

import (
	"bytes"
	"context"
	"crypto/tls"
	"embed"
	"fmt"
	"github.com/bopoh24/ma_1/notifier/internal/config"
	"github.com/bopoh24/ma_1/notifier/pkg/model"
	mail "github.com/xhit/go-simple-mail/v2"
	"strings"
	"text/template"
	"time"
)

type Service struct {
	conf       config.SMTP
	smtpServer *mail.SMTPServer
}

var (
	//go:embed templates
	templates embed.FS
)

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

	return &Service{conf: conf, smtpServer: server}, nil
}

func (s *Service) sendEmail(to, subject, body string, cc ...string) error {

	// SMTP client
	smtpClient, err := s.smtpServer.Connect()
	if err != nil {
		return err
	}
	defer func() {
		smtpClient.Quit()
		smtpClient.Close()
	}()

	email := mail.NewMSG()
	email.SetFrom(s.conf.From).AddTo(to).SetSubject(subject).SetBody(mail.TextHTML, body)
	if len(cc) > 0 {
		email.AddCc(cc...)
	}
	err = email.Send(smtpClient)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return nil
}

func (s *Service) BookingFailed(notification model.BookingNotification) error {
	return s.sendNotifications(notification, "failed to pay")
}

func (s *Service) BookingPaid(notification model.BookingNotification) error {
	return s.sendNotifications(notification, "paid")
}

func (s *Service) BookingSubmitted(notification model.BookingNotification) error {
	return s.sendNotifications(notification, "submitted")
}

func (s *Service) BookingCompleted(notification model.BookingNotification) error {
	return s.sendNotifications(notification, "completed")
}

func (s *Service) BookingCancelledByCustomer(notification model.BookingNotification) error {
	return s.sendNotifications(notification, "cancelled by customer")
}

func (s *Service) BookingCancelledByCompany(notification model.BookingNotification) error {

	return s.sendNotifications(notification, "cancelled by company")
}

func (s *Service) tplToString(templateName string, booking model.BookingNotification) (string, error) {

	tpl, err := template.New(templateName+".gohtml").ParseFS(templates, fmt.Sprintf("templates/%s.gohtml", templateName))
	if err != nil {
		return "", err
	}
	var output bytes.Buffer
	if err := tpl.Execute(&output, booking); err != nil {
		return "", err
	}
	return output.String(), nil
}

func (s *Service) sendNotifications(notification model.BookingNotification, subject string) error {
	notification.Status = subject

	cc := make([]string, 0, len(notification.CompanyManagerContacts))
	for _, c := range notification.CompanyManagerContacts {
		cc = append(cc, c.Email)
	}
	body, err := s.tplToString("company", notification)
	if err != nil {
		return err
	}

	subj := fmt.Sprintf("#%d OFFER - %s", notification.Offer.ID, strings.ToUpper(subject))

	// send email to company managers
	err = s.sendEmail(notification.CompanyContacts.Email, subj, body, cc...)
	if err != nil {
		return err
	}

	body, err = s.tplToString("customer", notification)
	if err != nil {
		return err
	}
	// send email to customer
	return s.sendEmail(notification.CustomerContacts.Email, subj, body)
}

// Close closes the Service
func (s *Service) Close(_ context.Context) error {
	return nil
}
