package model

import (
	"database/sql"
	"time"
)

type EmailVerificationToken struct {
	ID     int64 `json:"id"`
	UserID int64 `json:"user_id"`

	TokenHash string       `json:"-"`
	Purpose   string       `json:"purpose"`
	ExpiresAt time.Time    `json:"expires_at"`
	UsedAt    sql.NullTime `json:"used_at"`

	CreatedAt time.Time `json:"created_at"`
}
