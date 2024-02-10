package repository

import "errors"

// ErrCompanyNotFound is returned when a user is not found
var ErrCompanyNotFound = errors.New("company not found")
var ErrManagerNotFound = errors.New("manager not found")
