package models

type UserDB struct {
	ID       *int64  `json:"id,omitempty"`
	Username string  `json:"username"`
	Password *string `json:"password,omitempty"`
	SchoolID *int64  `json:"school_id,omitempty"`
}

type User struct {
	ID       *int64 `json:"id,omitempty"`
	Username string `json:"username"`
	SchoolID *int64 `json:"school_id,omitempty"`
}
