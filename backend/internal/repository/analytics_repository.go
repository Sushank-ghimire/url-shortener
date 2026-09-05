package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
) ([]model.Analytics, int, error) {
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

	countQuery := `
		SELECT COUNT(*) FROM analytics WHERE url_id = $1
	`

	var total int

	if err := r.db.QueryRowContext(ctx, countQuery, urlID).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		query,
		urlID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
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
			return nil, 0, err
		}

		analytics = append(analytics, a)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return analytics, total, nil
}

func (r *AnalyticsRepository) GetSummary(
	ctx context.Context,
	urlID int64,
) (*model.AnalyticsSummary, error) {
	summary := &model.AnalyticsSummary{}
	err := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*) FROM analytics WHERE url_id = $1`,
		urlID,
	).Scan(&summary.Clicks)

	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(
		ctx,
		`
			SELECT COUNT (DISTINCT ip_address) FROM analytics WHERE url_id = $1
		`,
		urlID,
	).Scan(&summary.Unique)
	if err != nil {
		return nil, err
	}
	// Countries
	summary.Countries, err = r.getDimension(
		ctx,
		urlID,
		"country",
	)
	if err != nil {
		return nil, err
	}

	// Devices
	summary.Devices, err = r.getDimension(
		ctx,
		urlID,
		"device",
	)
	if err != nil {
		return nil, err
	}

	// Browsers
	summary.Browsers, err = r.getDimension(
		ctx,
		urlID,
		"browser",
	)
	if err != nil {
		return nil, err
	}
	// Browsers
	summary.Referrers, err = r.getDimension(
		ctx,
		urlID,
		"referer",
	)
	if err != nil {
		return nil, err
	}
	return summary, nil
}

func (r *AnalyticsRepository) getDimension(
	ctx context.Context,
	urlID int64,
	column string,
) ([]model.AnalyticsDimension, error) {
	// Only allow known columns
	allowedColumns := map[string]bool{
		"country": true,
		"device":  true,
		"browser": true,
		"referer": true,
	}
	if !allowedColumns[column] {
		return nil, fmt.Errorf("invalid analytics dimension")
	}
	query := fmt.Sprintf(
		`
			SELECT COALESCE(NULLIF(%S, ''), 'Direct') AS name, COUNT(*) AS value
			FROM analytics WHERE url_id = $1 GROUP BY %s ORDER BY value DESC
		`,
		column,
		column,
	)

	rows, err := r.db.QueryContext(
		ctx,
		query,
		urlID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]model.AnalyticsDimension, 0)
	for rows.Next() {
		var item model.AnalyticsDimension
		if err := rows.Scan(
			&item.Name,
			&item.Value,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *AnalyticsRepository) GetTimeline(
	ctx context.Context,
	urlID int64,
	from time.Time,
	to time.Time,
) ([]model.ClickTimeline, error) {

	query := `
		SELECT
			d.date,
			COUNT(a.id) AS clicks
		FROM generate_series(
			$2::date,
			($3::date - INTERVAL '1 day')::date,
			INTERVAL '1 day'
		) AS d(date)

		LEFT JOIN analytics a
			ON a.url_id = $1
			AND a.clicked_at >= d.date
			AND a.clicked_at < d.date + INTERVAL '1 day'

		GROUP BY d.date
		ORDER BY d.date ASC;
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		urlID,
		from,
		to,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	timeline := make([]model.ClickTimeline, 0)

	for rows.Next() {
		var item model.ClickTimeline

		if err := rows.Scan(
			&item.Date,
			&item.Clicks,
		); err != nil {
			return nil, err
		}

		timeline = append(timeline, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return timeline, nil
}
