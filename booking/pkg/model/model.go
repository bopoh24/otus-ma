package model

import (
	"time"
)

type Service struct {
	ID       int    `json:"id,omitempty"`
	ParentID int    `json:"parent_id"`
	Name     string `json:"name"`
}

type OfferStatus string

const (
	OfferStatusOpen               OfferStatus = "open"
	OfferStatusFailed             OfferStatus = "failed"
	OfferStatusReserved           OfferStatus = "reserved"
	OfferStatusPaid               OfferStatus = "paid"
	OfferStatusSubmitted          OfferStatus = "submitted"
	OfferStatusCanceledByCompany  OfferStatus = "canceled_by_company"
	OfferStatusCanceledByCustomer OfferStatus = "canceled_by_customer"
	OfferStatusCompleted          OfferStatus = "completed"
)

type Offer struct {
	ID           int64       `json:"id,omitempty"`
	ServiceID    int         `json:"service_id"`
	ServiceName  string      `json:"service_name,omitempty"`
	Customer     string      `json:"customer,omitempty"`
	CompanyID    int64       `json:"company_id,omitempty"`
	CompanyName  string      `json:"company_name"`
	Location     []float64   `json:"location,omitempty"`
	Datetime     time.Time   `json:"datetime"`
	Description  string      `json:"description"`
	Price        float64     `json:"price"`
	Status       OfferStatus `json:"status"`
	CancelReason string      `json:"cancel_reason,omitempty"`
	CreatedBy    string      `json:"created_by,omitempty"`
	UpdatedBy    string      `json:"updated_by,omitempty"`
	CreatedAt    *time.Time  `json:"created_at,omitempty"`
	UpdatedAt    *time.Time  `json:"updated_at,omitempty"`
}
