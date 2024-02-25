package model

import "time"

type Customer struct {
	ID        string     `json:"id,omitempty"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Photo     string     `json:"photo,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Location  []float64  `json:"location,omitempty"`
	Created   *time.Time `json:"created,omitempty"`
	Updated   *time.Time `json:"updated,omitempty"`
}
