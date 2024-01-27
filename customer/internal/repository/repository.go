package repository

import (
	"github.com/bopoh24/ma_1/customer/internal/model"
)

type Repository interface {
	UserCreate(user *model.User) error
	UserByID(id int64) (*model.User, error)
	UserByExternalID(externalId string) (model.User, error)
	UserUpdate(user *model.User) error
	UserUpdateByExternalID(user *model.User) error
	UserDelete(id int64) error
}
