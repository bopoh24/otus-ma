package service

import (
	"context"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/pkg/verifier/phone"
)

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
	phoneVerifier phone.Verifier
}

// New returns a new Service instance
func New(repo Repository, phoneVerifier phone.Verifier) *Service {
	return &Service{
		repo:          repo,
		phoneVerifier: phoneVerifier,
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

// Close closes the service
func (s *Service) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
