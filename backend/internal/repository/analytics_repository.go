package repository

import (
	"context"
	"database/sql"

	"github.com/Sushank-ghimire/url-shortener/internal/model"
)

type AnalyticsRepository struct {
	db *sql.DB
}

func NewAnalyticsRepository(db *sql.DB) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

/*
INSERT/UPDATE/DELETE without RETURNING
ExecContext()

SELECT one row
QueryRowContext()

SELECT multiple rows
QueryContext()

INSERT ... RETURNING ...
QueryRowContext()
*/

func (r *AnalyticsRepository) Create(
	ctx context.Context,
	analytics *model.Analytics,
) error {
	query := `
		INSERT INTO analytics (url_id, ip_address, country, city, browser, os, device, referer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		analytics.URLID,
		analytics.IpAddress,
		analytics.Country,
		analytics.City,
		analytics.Browser,
		analytics.Os,
		analytics.Device,
		analytics.Referer,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *AnalyticsRepository) Get(
	ctx context.Context,
	urlID int64,
	limit int64,
	offset int64,
) ([]model.Analytics, error) {
	query := `
		SELECT
			id,
			url_id,
			ip_address,
			country,
			city,
			browser,
			os,
			device,
			referer,
			clicked_at
		FROM analytics
		WHERE url_id = $1
		ORDER BY clicked_at DESC
		LIMIT $2 OFFSET $3;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		urlID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []model.Analytics

	for rows.Next() {
		var a model.Analytics

		err := rows.Scan(
			&a.ID,
			&a.URLID,
			&a.IpAddress,
			&a.Country,
			&a.City,
			&a.Browser,
			&a.Os,
			&a.Device,
			&a.Referer,
			&a.ClickedAt,
		)
		if err != nil {
			return nil, err
		}

		analytics = append(analytics, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return analytics, nil
}
