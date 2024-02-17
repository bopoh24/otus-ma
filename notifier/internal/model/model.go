package model

import "github.com/bopoh24/ma_1/booking/pkg/model"

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

type BookingNotification struct {
	Offer                  model.Offer
	Status                 string                   `json:"status,omitempty"`
	CompanyContacts        CompanyContacts          `json:"company_contacts,omitempty"`
	CustomerContacts       CustomerContacts         `json:"customer_contacts"`
	CompanyManagerContacts []CompanyManagerContacts `json:"company_manager_contacts"`
}
