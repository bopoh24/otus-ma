package repository

import "errors"

// ErrUserNotFound is returned when a user is not found
var ErrUserNotFound = errors.New("user not found")

// ErrOrderNotFound is returned when an order is not found
var ErrOrderNotFound = errors.New("order not found")
