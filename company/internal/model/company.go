package model

import "time"

type Company struct {
	ID          int64      `json:"id,omitempty"`
	Owner       string     `json:"owner,omitempty"`
	Logo        string     `json:"logo,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone,omitempty"`
	Location    []float64  `json:"location,omitempty"`
	Active      bool       `json:"active,omitempty"`
	Created     *time.Time `json:"created,omitempty"`
	Updated     *time.Time `json:"updated,omitempty"`
}
