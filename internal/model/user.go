package model

type User struct {
	ID          int64  `json:"id,omitempty"`
	ExternalID  string `json:"externalId"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"LastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Description string `json:"description,omitempty"`
}
