package migrations

import "database/sql"

func CreateAnalyticsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS analytics (
		id BIGSERIAL PRIMARY KEY,
		url_id BIGINT NOT NULL,

		ip_address INET DEFAULT NULL,
		country VARCHAR(40) DEFAULT NULL,
		city VARCHAR(40) DEFAULT NULL,
		browser VARCHAR(40) DEFAULT NULL,
		os VARCHAR(40) DEFAULT NULL,
		device VARCHAR(40) DEFAULT NULL,
		referer TEXT NOT NULL DEFAULT 'share',

		clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);
	`

	_, err := db.Exec(query)
	return err
}
