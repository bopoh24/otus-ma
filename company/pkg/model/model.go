package model

import "time"

type Company struct {
	ID          int64      `json:"id,omitempty"`
	Logo        string     `json:"logo"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	Location    []float64  `json:"location"`
	Active      bool       `json:"active"`
	Created     *time.Time `json:"created"`
	Updated     *time.Time `json:"updated"`
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
