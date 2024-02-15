package model

import "time"

type Account struct {
	Customer string    `json:"customer"`
	Balance  int64     `json:"balance"`
	Created  time.Time `json:"created_at"`
}

type TransactionType string

const (
	TransactionTypeTopUp   TransactionType = "top-up"
	TransactionTypePayment TransactionType = "payment"
	TransactionTypeRefund  TransactionType = "refund"
)

type Transaction struct {
	ID       int64     `json:"id"`
	Type     string    `json:"transaction_type"`
	Customer string    `json:"customer"`
	OrderId  int64     `json:"order_id,omitempty"`
	Amount   int64     `json:"amount"`
	Created  time.Time `json:"created_at"`
}
