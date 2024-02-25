package model

import (
	"encoding/json"
	"github.com/bopoh24/ma_1/booking/pkg/model"
)

type CustomerContacts struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
}

type CompanyContacts struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type CompanyManagerContacts struct {
	Email string `json:"email"`
}

type NotificationType string

const (
	BookingPaid      NotificationType = "booking_paid"
	BookingFailed    NotificationType = "booking_failed"
	BookingSubmitted NotificationType = "booking_submitted"
	BookingCompleted NotificationType = "booking_completed"

	BookingCancelledByCustomer NotificationType = "booking_cancelled_by_customer"
	BookingCancelledByCompany  NotificationType = "booking_cancelled_by_company"
)

type BookingNotification struct {
	Offer                  model.Offer              `json:"offer"`
	Type                   NotificationType         `json:"type"`
	Status                 string                   `json:"status,omitempty"`
	FailReason             string                   `json:"fail_reason,omitempty"`
	CompanyContacts        CompanyContacts          `json:"company_contacts,omitempty"`
	CustomerContacts       CustomerContacts         `json:"customer_contacts"`
	CompanyManagerContacts []CompanyManagerContacts `json:"company_manager_contacts"`
}

func (b *BookingNotification) ToMessage() string {
	data, _ := json.Marshal(b)
	return string(data)
}
