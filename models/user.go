package models

type User struct {
	ID       *int64  `json:"id,omitempty"`
	Username string  `json:"username"`
	Password *string `json:"password,omitempty"`
}
