package memory

import (
	"context"
	"github.com/bopoh24/ma_1/company/internal/model"
	"github.com/bopoh24/ma_1/company/internal/repository"
)

type Memory struct {
	Companies map[int64]model.Company
}

// New returns a new Memory instance
func New() *Memory {
	return &Memory{
		Companies: make(map[int64]model.Company),
	}
}

func (m *Memory) CompanyCreate(_ context.Context, company model.Company) error {
	m.Companies[int64(len(m.Companies)+1)] = company
	return nil
}

func (m *Memory) CompanyUpdate(_ context.Context, company model.Company) error {
	if _, ok := m.Companies[company.ID]; !ok {
		return repository.ErrCompanyNotFound
	}
	m.Companies[company.ID] = company
	return nil
}

func (m *Memory) CompanyByID(_ context.Context, id int64) (model.Company, error) {
	company, ok := m.Companies[id]
	if !ok {
		return model.Company{}, repository.ErrCompanyNotFound
	}
	return company, nil
}

func (m *Memory) CompanyActivateDeactivate(_ context.Context, id int64, active bool) error {
	company, ok := m.Companies[id]
	if !ok {
		return repository.ErrCompanyNotFound
	}
	company.Active = active
	m.Companies[id] = company
	return nil
}

func (m Memory) CompanyUpdateLocation(_ context.Context, id int64, lat float64, lng float64) error {
	company, ok := m.Companies[id]
	if !ok {
		return repository.ErrCompanyNotFound
	}
	company.Location = []float64{lat, lng}
	m.Companies[id] = company
	return nil
}

func (m Memory) CompanyUpdateLogo(_ context.Context, id int64, logo string) error {
	company, ok := m.Companies[id]
	if !ok {
		return repository.ErrCompanyNotFound
	}
	company.Logo = logo
	m.Companies[id] = company
	return nil
}
