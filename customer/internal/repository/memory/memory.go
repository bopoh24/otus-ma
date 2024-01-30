package memory

import (
	"context"
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/customer/internal/repository"
)

type Memory struct {
	Customers map[string]model.Customer
}

// New returns a new Memory instance
func New() *Memory {
	return &Memory{
		Customers: make(map[string]model.Customer),
	}
}

func (m *Memory) CustomerCreate(_ context.Context, user model.Customer) error {
	m.Customers[user.ID] = user
	return nil
}

func (m *Memory) CustomerUpdate(_ context.Context, customer model.Customer) error {
	_, ok := m.Customers[customer.ID]
	if !ok {
		return repository.ErrCustomerNotFound
	}
	m.Customers[customer.ID] = customer
	return nil
}

func (m *Memory) CustomerByID(_ context.Context, id string) (model.Customer, error) {
	customer, ok := m.Customers[id]
	if !ok {
		return model.Customer{}, repository.ErrCustomerNotFound
	}
	return customer, nil
}

func (m *Memory) CustomerUpdatePhone(_ context.Context, id string, phone string) error {
	customer, ok := m.Customers[id]
	if !ok {
		return repository.ErrCustomerNotFound
	}
	customer.Phone = phone
	m.Customers[id] = customer
	return nil
}

func (m *Memory) CustomerUpdateLocation(_ context.Context, id string, lat float64, lng float64) error {
	customer, ok := m.Customers[id]
	if !ok {
		return repository.ErrCustomerNotFound
	}
	customer.Location = [2]float64{lat, lng}
	m.Customers[id] = customer
	return nil
}
