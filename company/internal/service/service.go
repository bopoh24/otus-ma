package service

import (
	"context"
	"github.com/bopoh24/ma_1/company/internal/config"
	"github.com/bopoh24/ma_1/company/internal/model"
)

type Repository interface {
	CompanyCreate(ctx context.Context, company model.Company) error
	CompanyUpdate(ctx context.Context, company model.Company) error
	CompanyByID(ctx context.Context, id int64) (model.Company, error)
	CompanyActivateDeactivate(ctx context.Context, id int64, active bool) error
	CompanyUpdateLocation(ctx context.Context, id int64, lat float64, lng float64) error
	CompanyUpdateLogo(ctx context.Context, id int64, logo string) error
}

type Service struct {
	repo Repository
	conf *config.Config
}

// New returns a new Service instance
func New(cfg *config.Config, repo Repository) *Service {
	return &Service{
		conf: cfg,
		repo: repo,
	}
}

// CompanyByID returns a company profile by its ID
func (s *Service) CompanyByID(ctx context.Context, id int64) (model.Company, error) {
	return s.repo.CompanyByID(ctx, id)
}

// CreateCompany creates a new company profile
func (s *Service) CreateCompany(ctx context.Context, company model.Company) error {
	return s.repo.CompanyCreate(ctx, company)
}

// UpdateCompany updates a company profile
func (s *Service) UpdateCompany(ctx context.Context, company model.Company) error {
	return s.repo.CompanyUpdate(ctx, company)
}

// UpdateCompanyLocation updates a company location
func (s *Service) UpdateCompanyLocation(ctx context.Context, id int64, lat float64, lng float64) error {
	return s.repo.CompanyUpdateLocation(ctx, id, lat, lng)
}

// UpdateCompanyLogo updates a company logo
func (s *Service) UpdateCompanyLogo(ctx context.Context, id int64, logo string) error {
	return s.repo.CompanyUpdateLogo(ctx, id, logo)
}

// ActivateDeactivateCompany activates or deactivates a company
func (s *Service) ActivateDeactivateCompany(ctx context.Context, id int64, active bool) error {
	return s.repo.CompanyActivateDeactivate(ctx, id, active)
}
