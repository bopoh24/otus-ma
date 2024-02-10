package model

import "time"

type Company struct {
	ID          int64      `json:"id,omitempty"`
	Logo        string     `json:"logo,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone,omitempty"`
	Email       string     `json:"email,omitempty"`
	Location    []float64  `json:"location,omitempty"`
	Active      bool       `json:"active,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
}

type MangerRole string

const (
	MangerRoleAdmin   MangerRole = "admin"
	MangerRoleManager MangerRole = "manager"
)

type Manager struct {
	ID        int64      `json:"id,omitempty"`
	CompanyID int64      `json:"company_id"`
	UserID    string     `json:"user_id"`
	Email     string     `json:"email"`
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
	Role      MangerRole `json:"role"`
	Active    bool       `json:"active,omitempty"`
	Created   *time.Time `json:"created,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
}
