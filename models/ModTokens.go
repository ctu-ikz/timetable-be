package models

import "time"

type RefreshTokenDB struct {
	ID        int64     `json:"id,omitempty"`
	Token     string    `json:"token"`
	UserID    int64     `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	IsValid   bool      `json:"is_valid"`
}
