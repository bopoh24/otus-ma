package user

import (
	"github.com/bopoh24/ma_1/app/internal/model"
	"github.com/bopoh24/ma_1/app/internal/repository"
)

type Service struct {
	repo repository.UserRepository
}

// New returns a new Service instance
func New(repo repository.UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// Create creates a new user
func (s *Service) Create(user *model.User) error {
	return s.repo.UserCreate(user)
}

// ByID returns a user by id
func (s *Service) ByID(userID int64) (*model.User, error) {
	return s.repo.UserByID(userID)
}

// ByExternalID returns a user by external id
func (s *Service) ByExternalID(userID string) (model.User, error) {
	return s.repo.UserByExternalID(userID)
}

// Update updates a user
func (s *Service) Update(user *model.User) error {
	return s.repo.UserUpdate(user)
}

func (s *Service) UpdateByExternalID(user *model.User) error {
	return s.repo.UserUpdateByExternalID(user)
}

// Delete deletes a user
func (s *Service) Delete(userID int64) error {
	return s.repo.UserDelete(userID)
}
