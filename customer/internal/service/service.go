package service

import (
	"context"
	"encoding/json"
	"fmt"
	bookingModel "github.com/bopoh24/ma_1/booking/pkg/model"
	companyModel "github.com/bopoh24/ma_1/company/pkg/model"
	"github.com/bopoh24/ma_1/customer/internal/model"
	notifierModel "github.com/bopoh24/ma_1/notifier/pkg/model"
	"github.com/bopoh24/ma_1/pkg/http/client"
	"github.com/bopoh24/ma_1/pkg/kafka/producer"
	"github.com/bopoh24/ma_1/pkg/verifier/phone"
	"net/http"
)

const bookingNotificationTopic = "booking_notification"

//go:generate mockgen -source service.go -destination ../../mocks/repository.go -package mock Repository
type Repository interface {
	CustomerCreate(ctx context.Context, customer model.Customer) error
	CustomerUpdate(ctx context.Context, customer model.Customer) error
	CustomerByID(ctx context.Context, id string) (model.Customer, error)
	CustomerUpdatePhone(ctx context.Context, id string, phone string) error
	CustomerUpdateLocation(ctx context.Context, id string, lat float64, lng float64) error
	Close(ctx context.Context) error
}

type Service struct {
	repo          Repository
	producer      *producer.KafkaProducer
	phoneVerifier phone.Verifier
	companyClient client.HttpRequester
	bookingClient client.HttpRequester
	paymentClient client.HttpRequester
}

// New returns a new Service instance
func New(repo Repository, producer *producer.KafkaProducer, phoneVerifier phone.Verifier,
	companyClient client.HttpRequester,
	bookingClient client.HttpRequester,
	paymentClient client.HttpRequester) *Service {
	return &Service{
		repo:          repo,
		producer:      producer,
		phoneVerifier: phoneVerifier,
		companyClient: companyClient,
		bookingClient: bookingClient,
		paymentClient: paymentClient,
	}
}

// CustomerByID returns a customer profile by id
func (s *Service) CustomerByID(ctx context.Context, userID string) (model.Customer, error) {
	return s.repo.CustomerByID(ctx, userID)
}

// CreateCustomerProfile creates a new customer profile
func (s *Service) CreateCustomerProfile(ctx context.Context, customer model.Customer) error {
	return s.repo.CustomerCreate(ctx, customer)
}

// UpdateCustomerProfile updates a customer profile
func (s *Service) UpdateCustomerProfile(ctx context.Context, customer model.Customer) error {
	return s.repo.CustomerUpdate(ctx, customer)
}

// RequestPhoneVerification requests phone verification
func (s *Service) RequestPhoneVerification(ctx context.Context, phone string) error {
	return s.phoneVerifier.Send(ctx, phone)
}

// VerifyPhone verifies phone number
func (s *Service) VerifyPhone(ctx context.Context, id string, phone string, code string) error {
	err := s.phoneVerifier.Check(ctx, phone, code)
	if err != nil {
		return err
	}
	return s.repo.CustomerUpdatePhone(ctx, id, phone)
}

// UpdateCustomerLocation updates a customer location
func (s *Service) UpdateCustomerLocation(ctx context.Context, id string, lat float64, lng float64) error {
	return s.repo.CustomerUpdateLocation(ctx, id, lat, lng)
}

// CompanyByID returns a company profile by its ID
func (s *Service) CompanyByID(ctx context.Context, id int64) (companyModel.Company, error) {
	resp, err := s.companyClient.Get(ctx, fmt.Sprintf("/company/%d", id), nil)
	if err != nil {
		return companyModel.Company{}, err
	}
	var company companyModel.Company
	err = json.NewDecoder(resp.Body).Decode(&company)
	if err != nil {
		return companyModel.Company{}, err
	}
	return company, nil
}

// CompanyManagers returns company managers
func (s *Service) CompanyManagers(ctx context.Context, id int64) ([]companyModel.Manager, error) {
	resp, err := s.companyClient.Get(ctx, fmt.Sprintf("/company/%d/managers", id), nil)
	if err != nil {
		return nil, err
	}
	var managers []companyModel.Manager
	err = json.NewDecoder(resp.Body).Decode(&managers)
	if err != nil {
		return nil, err
	}
	return managers, nil
}

func (s *Service) BookOffer(ctx context.Context, id string, headers http.Header) (*bookingModel.Offer, int, error) {
	// book offer
	respBooking, err := s.bookingClient.Post(ctx,
		fmt.Sprintf("/booking/offers/%s/book", id), nil, headers)
	if err != nil {
		return nil, 0, err
	}
	defer respBooking.Body.Close()

	if respBooking.StatusCode != http.StatusOK {
		if respBooking.StatusCode == http.StatusNotFound {
			return nil, http.StatusNotFound, fmt.Errorf("offer not found")
		}

		errMessage := struct {
			Error string `json:"error"`
		}{}
		err = json.NewDecoder(respBooking.Body).Decode(&errMessage)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		return nil, respBooking.StatusCode, fmt.Errorf(errMessage.Error)
	}
	var offer bookingModel.Offer
	err = json.NewDecoder(respBooking.Body).Decode(&offer)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return &offer, http.StatusOK, nil
}

func (s *Service) ServiceNameById(ctx context.Context, serviceId int) (string, error) {
	resp, err := s.bookingClient.Get(ctx, "/booking/services", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	services := make([]bookingModel.Service, 0)
	err = json.NewDecoder(resp.Body).Decode(&services)
	if err != nil {
		return "", err
	}
	for _, service := range services {
		if service.ID == serviceId {
			return service.Name, nil
		}
	}
	return "", fmt.Errorf("service not found")
}

func (s *Service) PrepareNotification(ctx context.Context, customerId string, companyId int64) (notifierModel.BookingNotification, error) {
	var notification notifierModel.BookingNotification
	// get customer
	customer, err := s.CustomerByID(ctx, customerId)
	if err != nil {
		return notification, err
	}

	// get company
	company, err := s.CompanyByID(ctx, companyId)
	if err != nil {
		return notification, err
	}
	// get company manager
	companyManager, err := s.CompanyManagers(ctx, companyId)

	notification = notifierModel.BookingNotification{
		CustomerContacts: notifierModel.CustomerContacts{
			Email:     customer.Email,
			FirstName: customer.FirstName,
			LastName:  customer.LastName,
			Phone:     customer.Phone,
		},
		CompanyContacts: notifierModel.CompanyContacts{
			Name:    company.Name,
			Email:   company.Email,
			Phone:   company.Phone,
			Address: company.Address,
		},
		CompanyManagerContacts: make([]notifierModel.CompanyManagerContacts, len(companyManager)),
	}
	for i, manager := range companyManager {
		notification.CompanyManagerContacts[i] = notifierModel.CompanyManagerContacts{
			Email: manager.Email,
		}
	}
	return notification, nil
}

func (s *Service) MakePayment(ctx context.Context, offer bookingModel.Offer, headers http.Header) (int, error) {
	// make payment
	resp, err := s.paymentClient.Post(ctx, "/payment/make", map[string]any{
		"offer_id": offer.ID,
		"amount":   offer.Price,
	}, headers)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, s.extractError(resp)
}

func (s *Service) BookingReset(ctx context.Context, offerId int64) (int, error) {
	resp, err := s.bookingClient.Put(ctx, fmt.Sprintf("/booking/offers/%d/reset", offerId), nil, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return resp.StatusCode, s.extractError(resp)
}

func (s *Service) BookingPaid(ctx context.Context, offerId int64) (int, error) {
	resp, err := s.bookingClient.Put(ctx, fmt.Sprintf("/booking/offers/%d/paid", offerId), nil, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return resp.StatusCode, s.extractError(resp)
}

func (s *Service) SendNotification(notification *notifierModel.BookingNotification) error {
	return s.producer.SendMessage(bookingNotificationTopic, notification.ToMessage())
}

func (s *Service) extractError(resp *http.Response) error {
	errMessage := struct {
		Error string `json:"error"`
	}{}
	err := json.NewDecoder(resp.Body).Decode(&errMessage)
	if err != nil {
		return nil
	}
	return fmt.Errorf(errMessage.Error)
}

// Close closes the service
func (s *Service) Close(ctx context.Context) error {
	if err := s.producer.Close(); err != nil {
		return err
	}
	return s.repo.Close(ctx)
}
