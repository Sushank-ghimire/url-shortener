package dto

import "time"

type URLResponse struct {
	ID          int64      `json:"id"`
	ShortCode   string     `json:"short_code"`
	LongURL     string     `json:"long_url"`
	Title       string     `json:"title"`
	CustomAlias string     `json:"custom_alias"`
	IsActive    bool       `json:"is_active"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ShortURLRequest struct {
	Title       string     `json:"title"`
	CustomAlias string     `json:"custom_alias"`
	LongURL     string     `json:"long_url" binding:"required,url"`
	ExpiresAt   *time.Time `json:"expires_at"`
}
