package model

import (
	"github.com/bopoh24/ma_1/pkg/http/helper"
	"time"
)

type Service struct {
	ID       int    `json:"id,omitempty"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
}

type OrderStatus string

const (
	OrderStatusOpen               OrderStatus = "open"
	OrderStatusReserved           OrderStatus = "reserved"
	OrderStatusPaid               OrderStatus = "paid"
	OrderStatusSubmitted          OrderStatus = "submitted"
	OrderStatusCanceledByCustomer OrderStatus = "canceled_by_customer"
	OrderStatusCanceledByCompany  OrderStatus = "canceled_by_company"
	OrderStatusCompleted          OrderStatus = "completed"
)

type Offer struct {
	ID           int64           `json:"id,omitempty"`
	ServiceID    int             `json:"service_id"`
	ServiceName  string          `json:"service_name,omitempty"`
	Customer     string          `json:"customer"`
	CompanyID    int64           `json:"company_id"`
	CompanyName  string          `json:"company_name"`
	Location     []float64       `json:"location"`
	Datetime     time.Time       `json:"datetime"`
	Description  string          `json:"description"`
	Price        float64         `json:"price"`
	Status       OrderStatus     `json:"status"`
	CancelReason string          `json:"cancel_reason,omitempty"`
	CreatedBy    string          `json:"created_by,omitempty"`
	UpdatedBy    string          `json:"updated_by,omitempty"`
	CreatedAt    helper.JSONTime `json:"created_at,omitempty"`
	UpdatedAt    helper.JSONTime `json:"updated_at,omitempty"`
}
