package service

import (
	"context"
)

//go:generate mockgen -source service.go -destination ../../mocks/repository.go -package mock Repository
type Repository interface {
	CreateAccount(ctx context.Context, customerID string) error
	Balance(ctx context.Context, customerID string) (float32, error)
	TopUp(ctx context.Context, customerID string, amount float32) error
	PaymentMake(ctx context.Context, orderId int64, customerID string, amount float32) error
	PaymentCancel(ctx context.Context, orderId int64) error
	Close(ctx context.Context) error
}

type Service struct {
	repo Repository
}

// New returns a new Service instance
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateAccount creates a new account
func (s *Service) CreateAccount(ctx context.Context, customerID string) error {
	return s.repo.CreateAccount(ctx, customerID)
}

// TopUp tops up the account
func (s *Service) TopUp(ctx context.Context, customerID string, amount float32) error {
	return s.repo.TopUp(ctx, customerID, amount)
}

// Balance returns the balance of the account
func (s *Service) Balance(ctx context.Context, customerID string) (float32, error) {
	return s.repo.Balance(ctx, customerID)
}

// PaymentMake creates a new payment
func (s *Service) PaymentMake(ctx context.Context, orderId int64, customerID string, amount float32) error {
	return s.repo.PaymentMake(ctx, orderId, customerID, amount)
}

// PaymentCancel cancels the payment
func (s *Service) PaymentCancel(ctx context.Context, orderId int64) error {
	return s.repo.PaymentCancel(ctx, orderId)
}

// Close closes the Service
func (s *Service) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
