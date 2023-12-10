package memory

import (
	"github.com/bopoh24/ma_1/app/internal/model"
	"github.com/bopoh24/ma_1/app/internal/repository"
)

type Memory struct {
	Users  map[int64]*model.User
	Orders map[int64]*model.Order
}

var _ repository.UserRepository = (*Memory)(nil)
var _ repository.OrderRepository = (*Memory)(nil)

// New returns a new Memory instance
func New() *Memory {
	return &Memory{
		Users:  make(map[int64]*model.User),
		Orders: make(map[int64]*model.Order),
	}
}

// UserCreate creates a new user
func (m *Memory) UserCreate(user *model.User) error {
	user.ID = int64(len(m.Users) + 1)
	m.Users[user.ID] = user
	return nil
}

// UserByID returns a user by id
func (m *Memory) UserByID(id int64) (*model.User, error) {
	user, ok := m.Users[id]
	if !ok {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

// UserUpdate updates a user
func (m *Memory) UserUpdate(user *model.User) error {
	_, ok := m.Users[user.ID]
	if !ok {
		return repository.ErrUserNotFound
	}
	m.Users[user.ID] = user
	return nil
}

// UserDelete deletes a user
func (m *Memory) UserDelete(id int64) error {
	_, ok := m.Users[id]
	if !ok {
		return repository.ErrUserNotFound
	}
	delete(m.Users, id)
	return nil
}

// UserByExternalID returns a user by external id
func (m *Memory) UserByExternalID(externalId string) (model.User, error) {
	for _, user := range m.Users {
		if user.ExternalID == externalId {
			return *user, nil
		}
	}
	return model.User{}, repository.ErrUserNotFound
}

// UserUpdateByExternalID updates a user
func (m *Memory) UserUpdateByExternalID(user *model.User) error {
	for _, u := range m.Users {
		if u.ExternalID == user.ExternalID {
			u.FirstName = user.FirstName
			u.LastName = user.LastName
			u.Email = user.Email
			u.Phone = user.Phone
			u.Description = user.Description
			return nil
		}
	}
	return repository.ErrUserNotFound
}

// OrderCreate creates a new order
func (m *Memory) OrderCreate(order *model.Order) (int64, error) {
	//max ID search
	var maxID int64
	for _, o := range m.Orders {
		if o.ID > maxID {
			maxID = o.ID
		}
	}
	order.ID = maxID + 1
	m.Orders[order.ID] = order
	return order.ID, nil
}

// OrderByID returns a order by id
func (m *Memory) OrderByID(id int64) (*model.Order, error) {
	order, ok := m.Orders[id]
	if !ok {
		return nil, repository.ErrOrderNotFound
	}
	return order, nil

}

// OrderByUserID returns orders by user id
func (m *Memory) OrderByUserID(id int64) ([]model.Order, error) {
	orders := make([]model.Order, 0)
	for _, order := range m.Orders {
		if order.UserID == id {
			orders = append(orders, *order)
		}
	}
	return orders, nil
}

// OrderUpdate updates a order
func (m *Memory) OrderUpdate(order *model.Order) error {
	_, ok := m.Orders[order.ID]
	if !ok {
		return repository.ErrOrderNotFound
	}
	m.Orders[order.ID] = order
	return nil
}

// OrderDelete deletes a order
func (m *Memory) OrderDelete(id int64) error {
	_, ok := m.Orders[id]
	if !ok {
		return repository.ErrOrderNotFound
	}
	delete(m.Orders, id)
	return nil
}
