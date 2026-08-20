package model

import "time"

type Url struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	ShortCode   string     `json:"short_code"`
	LongURL     string     `json:"long_url"`
	Title       string     `json:"title"`
	CustomAlias string     `json:"custom_alias"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
