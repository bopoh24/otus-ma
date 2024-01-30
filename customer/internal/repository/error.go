package repository

import "errors"

// ErrCustomerNotFound is returned when a user is not found
var ErrCustomerNotFound = errors.New("customer not found")
