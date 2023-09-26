package model

type User struct {
	ID        int64  `json:"id,omitempty"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
