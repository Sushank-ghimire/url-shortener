package migrations

import "database/sql"

func CreateUsersTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(55) NOT NULL,
		email VARCHAR(255) NOT NULL,
		pass_hash TEXT NOT NULL,
		is_verified BOOLEAN NOT NULL DEFAULT FALSE,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		CONSTRAINT users_email_unique UNIQUE(email)
	);
	`

	_, err := db.Exec(query)
	return err
}
