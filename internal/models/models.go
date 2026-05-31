package models

import "time"

 type URLRecord struct {
	ID string `json:"id"`
	Slug string `json:"slug"`
	LongURL string `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Clicks int `json:"clicks"`
	UserID string `json:"user_id"`
}

type User struct {
	ID string `json:"id"`
Email string `json:"email"`
PasswordHash string `json:"-"`
CreatedAt time.Time `json:"created_at"`
}