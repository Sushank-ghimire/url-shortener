package migrations

import "database/sql"

func CreateUrlsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL,

		short_code VARCHAR(8) NOT NULL,
		long_url TEXT NOT NULL,
		title VARCHAR(50) DEFAULT NULL,
		custom_alias VARCHAR(8) DEFAULT NULL,
		is_active BOOLEAN DEFAULT TRUE,

		expires_at TIMESTAMPTZ DEFAULT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

		FOREIGN KEY (user_id) REFERENCES users(id),
		CONSTRAINT custom_alias_unique UNIQUE(custom_alias)
	);
	`

	_, err := db.Exec(query)
	return err
}
