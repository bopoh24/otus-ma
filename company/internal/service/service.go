package service

import (
	"context"
	"github.com/bopoh24/ma_1/company/internal/config"
	"github.com/bopoh24/ma_1/company/internal/model"
)

//go:generate mockgen -source service.go -destination ../../mocks/repository.go -package mock Repository
type Repository interface {
	CompanyCreate(ctx context.Context, userId, email, firstName, lastName string, company model.Company) error
	CompanyUpdate(ctx context.Context, company model.Company) error
	CompanyByID(ctx context.Context, id int64) (model.Company, error)
	CompanyActivateDeactivate(ctx context.Context, id int64, active bool) error
	CompanyUpdateLocation(ctx context.Context, id int64, lat float64, lng float64) error
	CompanyUpdateLogo(ctx context.Context, id int64, logo string) error

	CompanyManagers(ctx context.Context, companyID int64) ([]model.Manager, error)
	ManagerByID(ctx context.Context, id int64) (model.Manager, error)
	ManagerByUserID(ctx context.Context, userId string) (model.Manager, error)
	ManagerByEmail(ctx context.Context, email string) (model.Manager, error)

	ManagerInvite(ctx context.Context, companyID int64, email string, role model.MangerRole) error
	ManagerActivateDeactivate(ctx context.Context, id int64, active bool) error
	ManagerSetRole(ctx context.Context, id int64, role model.MangerRole) error
	ManagerDelete(ctx context.Context, id int64) error
	Close(ctx context.Context) error
}

type Service struct {
	repo Repository
	conf *config.Config
}

// New returns a new Service instance
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CompanyByID returns a company profile by its ID
func (s *Service) CompanyByID(ctx context.Context, id int64) (model.Company, error) {
	return s.repo.CompanyByID(ctx, id)
}

// CreateCompany creates a new company profile
func (s *Service) CreateCompany(ctx context.Context, userId, email, firstName, lastName string, company model.Company) error {
	return s.repo.CompanyCreate(ctx, userId, email, firstName, lastName, company)
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

func (s *Service) CompanyManagers(ctx context.Context, companyID int64) ([]model.Manager, error) {
	return s.repo.CompanyManagers(ctx, companyID)
}

func (s *Service) ManagerByID(ctx context.Context, id int64) (model.Manager, error) {
	return s.repo.ManagerByID(ctx, id)
}

func (s *Service) ManagerByUserID(ctx context.Context, userId string) (model.Manager, error) {
	return s.repo.ManagerByUserID(ctx, userId)
}

func (s *Service) ManagerByEmail(ctx context.Context, email string) (model.Manager, error) {
	return s.repo.ManagerByEmail(ctx, email)
}

func (s *Service) ManagerInvite(ctx context.Context, companyID int64, email string, role model.MangerRole) error {
	return s.repo.ManagerInvite(ctx, companyID, email, role)
}

func (s *Service) ManagerActivateDeactivate(ctx context.Context, id int64, active bool) error {
	return s.repo.ManagerActivateDeactivate(ctx, id, active)
}

func (s *Service) ManagerSetRole(ctx context.Context, id int64, role model.MangerRole) error {
	return s.repo.ManagerSetRole(ctx, id, role)
}

func (s *Service) ManagerDelete(ctx context.Context, id int64) error {
	return s.repo.ManagerDelete(ctx, id)
}

func (s *Service) Close(ctx context.Context) error {
	return s.repo.Close(ctx)
}
