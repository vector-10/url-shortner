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