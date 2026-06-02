package models

import "time"

type URLRecord struct {
    ID          string     `json:"id"`
    Slug        string     `json:"slug"`
    LongURL     string     `json:"long_url"`
    UserID      string     `json:"user_id"`
    CreatedAt   time.Time  `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at"`
    IsActive    bool       `json:"is_active"`
    MaxClicks   *int       `json:"max_clicks"`
    TotalClicks int        `json:"total_clicks"`
    LinkType    string     `json:"link_type"`
}

type ClickEvent struct {
    ID              int64      `json:"id"`
    Slug            string     `json:"slug"`
    ClickedAt       time.Time  `json:"clicked_at"`
    IPAddress       string     `json:"ip_address"`
    UserAgent       string     `json:"user_agent"`
    WasValid        bool       `json:"was_valid"`
    RejectionReason string     `json:"rejection_reason"`
}
