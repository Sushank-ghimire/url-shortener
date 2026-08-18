package migrations

import (
	"database/sql"
	"fmt"
)

func MigrateAll(db *sql.DB) error {
	if err := CreateUsersTable(db); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	if err := CreateUrlsTable(db); err != nil {
		return fmt.Errorf("create urls table: %w", err)
	}
	if err := CreateAnalyticsTable(db); err != nil {
		return fmt.Errorf("create analytics table: %w", err)
	}
	if err := CreateRefreshTokenTable(db); err != nil {
		return fmt.Errorf("create refresh token table: %w", err)
	}
	if err := CreateVerificationTokenTable(db); err != nil {
		return fmt.Errorf("create email verification token table: %w", err)
	}
	return nil
}
