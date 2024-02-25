package repository

import "errors"

var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
