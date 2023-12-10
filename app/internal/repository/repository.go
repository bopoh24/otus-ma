package repository

import (
	"github.com/bopoh24/ma_1/app/internal/model"
)

type UserRepository interface {
	UserCreate(user *model.User) error
	UserByID(id int64) (*model.User, error)
	UserByExternalID(externalId string) (model.User, error)
	UserUpdate(user *model.User) error
	UserUpdateByExternalID(user *model.User) error
	UserDelete(id int64) error
}

type OrderRepository interface {
	OrderCreate(order *model.Order) (int64, error)
	OrderByID(id int64) (*model.Order, error)
	OrderByUserID(id int64) ([]model.Order, error)
	OrderUpdate(order *model.Order) error
	OrderDelete(id int64) error
}
