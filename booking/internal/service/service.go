package service

import (
	"context"
	"github.com/bopoh24/ma_1/booking/internal/config"
	"github.com/bopoh24/ma_1/booking/internal/model"
	"time"
)

//go:generate mockgen -source service.go -destination ../../mocks/repository.go -package mock Repository
type Repository interface {
	Services(ctx context.Context) ([]model.Service, error)
	ServiceAdd(ctx context.Context, service model.Service) error
	OfferAdd(ctx context.Context, offer model.Offer) error
	OfferDelete(ctx context.Context, id int64, companyId int64) error
	OfferChangeStatus(ctx context.Context, id int64, status model.OrderStatus) error
	OfferCancelByCompany(ctx context.Context, id int64, reason string, companyId int64, managerId string) error
	OfferCancelByCustomer(ctx context.Context, id int64, reason string, customerId string) error

	// TODO: add search by location
	OfferSearch(ctx context.Context, serviceId int64, from, to time.Time, page, limit int) ([]model.Offer, error)
	Book(ctx context.Context, offerId int64, customerId string) error

	CompanyOffers(ctx context.Context, companyId int64, page, limit int) ([]model.Offer, error)
	CustomerOffers(ctx context.Context, customerId string, page, limit int) ([]model.Offer, error)

	Close(ctx context.Context) error
}

type Service struct {
	repo     Repository
	conf     *config.Config
	services []model.Service
}

// New returns a new Service instance
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// Services returns all services
func (s *Service) Services(ctx context.Context) ([]model.Service, error) {
	var err error
	if s.services == nil {
		s.services, err = s.repo.Services(ctx)
		if err != nil {
			return nil, err
		}
	}
	return s.services, nil
}

// ServiceAdd adds a new service
func (s *Service) ServiceAdd(ctx context.Context, service model.Service) error {
	return s.repo.ServiceAdd(ctx, service)
}

// OfferAdd adds a new offer
func (s *Service) OfferAdd(ctx context.Context, offer model.Offer) error {
	return s.repo.OfferAdd(ctx, offer)
}

// OfferDelete deletes an offer
func (s *Service) OfferDelete(ctx context.Context, id int64, companyId int64) error {
	return s.repo.OfferDelete(ctx, id, companyId)
}

// OfferChangeStatus changes the status of an offer
func (s *Service) OfferChangeStatus(ctx context.Context, id int64, status model.OrderStatus) error {
	return s.repo.OfferChangeStatus(ctx, id, status)
}

// OfferCancelByCompany cancels an offer by company
func (s *Service) OfferCancelByCompany(ctx context.Context, id int64, reason string,
	companyId int64, managerId string) error {
	return s.repo.OfferCancelByCompany(ctx, id, reason, companyId, managerId)
}

// OfferCancelByCustomer cancels an offer by user
func (s *Service) OfferCancelByCustomer(ctx context.Context, id int64, reason string, customerId string) error {
	return s.repo.OfferCancelByCustomer(ctx, id, reason, customerId)
}

// OfferSearch searches for offers
func (s *Service) OfferSearch(ctx context.Context, serviceId int64, from, to time.Time, page, limit int) ([]model.Offer, error) {
	return s.repo.OfferSearch(ctx, serviceId, from, to, page, limit)
}

// Book books an offer
func (s *Service) Book(ctx context.Context, offerId int64, customerId string) error {
	return s.repo.Book(ctx, offerId, customerId)
}

// CompanyOffers returns offers of a company
func (s *Service) CompanyOffers(ctx context.Context, companyId int64, page, limit int) ([]model.Offer, error) {
	services, err := s.Services(ctx)
	if err != nil {
		return nil, err
	}
	offers, err := s.repo.CompanyOffers(ctx, companyId, page, limit)
	if err != nil {
		return nil, err
	}
	// add service name to offers
	for i := range offers {
		for _, service := range services {
			if offers[i].ServiceID == service.ID {
				offers[i].ServiceName = service.Name
				break
			}
		}
	}
	return offers, nil
}

// CustomerOffers returns offers of a customer
func (s *Service) CustomerOffers(ctx context.Context, customerId string, page, limit int) ([]model.Offer, error) {
	services, err := s.Services(ctx)
	if err != nil {
		return nil, err
	}
	offers, err := s.repo.CustomerOffers(ctx, customerId, page, limit)
	if err != nil {
		return nil, err
	}
	// add service name to offers
	for i := range offers {
		for _, service := range services {
			if offers[i].ServiceID == service.ID {
				offers[i].ServiceName = service.Name
				break
			}
		}
	}
	return offers, nil
}

// Close closes the Service
func (s *Service) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
