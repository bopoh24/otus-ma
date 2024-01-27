package memory

import (
	"github.com/bopoh24/ma_1/customer/internal/model"
	"github.com/bopoh24/ma_1/customer/internal/repository"
)

type Memory struct {
	Users map[int64]*model.User
}

func New() *Memory {
	return &Memory{
		Users: make(map[int64]*model.User),
	}
}

func (m *Memory) UserCreate(user *model.User) error {
	user.ID = int64(len(m.Users) + 1)
	m.Users[user.ID] = user
	return nil
}

func (m *Memory) UserByID(id int64) (*model.User, error) {
	user, ok := m.Users[id]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *Memory) UserUpdate(user *model.User) error {
	_, ok := m.Users[user.ID]
	if !ok {
		return repository.ErrUserNotFound
	}
	m.Users[user.ID] = user
	return nil
}

func (m *Memory) UserDelete(id int64) error {
	_, ok := m.Users[id]
	if !ok {
		return repository.ErrUserNotFound
	}
	delete(m.Users, id)
	return nil
}
