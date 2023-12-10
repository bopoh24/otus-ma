package order

import (
	"github.com/bopoh24/ma_1/app/internal/model"
	"github.com/bopoh24/ma_1/app/internal/repository"
	"sync"
)

type Service struct {
	repo            repository.OrderRepository
	idempotencyKeys map[string]int64
	mu              sync.Mutex
}

// New returns a new Service instance
func New(repo repository.OrderRepository) *Service {
	return &Service{
		repo:            repo,
		idempotencyKeys: make(map[string]int64),
	}
}

func (s *Service) Create(order *model.Order, idempotencyKey string) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// check idempotency key if order already exists return order id
	if orderID, ok := s.idempotencyKeys[idempotencyKey]; ok {
		return orderID, nil
	}
	// create order
	var err error
	order.ID, err = s.repo.OrderCreate(order)
	if err != nil {
		return 0, err
	}
	s.idempotencyKeys[idempotencyKey] = order.ID
	return order.ID, nil
}

func (s *Service) ByID(orderID int64) (*model.Order, error) {
	return s.repo.OrderByID(orderID)
}

func (s *Service) ByUserID(userID int64) ([]model.Order, error) {
	return s.repo.OrderByUserID(userID)
}

func (s *Service) Update(order *model.Order) error {
	return s.repo.OrderUpdate(order)
}

func (s *Service) Delete(orderID int64) error {
	return s.repo.OrderDelete(orderID)
}
