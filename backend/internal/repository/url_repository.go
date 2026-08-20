package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Sushank-ghimire/url-shortener/internal/model"
)

type URLRepository struct {
	db *sql.DB
}

var (
	ErrURLNotFound         = errors.New("requested url with id not found")
	ErrFailedToDeactiveURL = errors.New("failed to deactivate url")
	ErrFailedToActiveURL   = errors.New("failed to activate url")
	ErrCustomAlias         = errors.New("custom alias is required")
)

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{
		db: db,
	}
}

func (r *URLRepository) Create(
	ctx context.Context,
	url *model.Url,
) error {
	query := `
		INSERT INTO urls (
			id,
			user_id,
			short_code,
			long_url,
			title,
			expires_at,
			custom_alias
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at, updated_at, expires_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		url.ID,
		url.UserID,
		url.ShortCode,
		url.LongURL,
		url.Title,
		url.ExpiresAt,
		url.CustomAlias,
	).Scan(
		&url.CreatedAt,
		&url.UpdatedAt,
		&url.ExpiresAt,
	)
}

func (r *URLRepository) GetURLByShortCode(
	ctx context.Context,
	code string,
) (urlData *model.Url, error error) {

	var url model.Url

	query := `
			SELECT id, user_id, custom_alias, is_active, long_url, short_code FROM urls WHERE short_code = $1
		`
	err := r.db.QueryRowContext(
		ctx,
		query,
		code,
	).Scan(
		&url.ID,
		&url.UserID,
		&url.CustomAlias,
		&url.IsActive,
		&url.LongURL,
		&url.ShortCode,
	)
	return &url, err
}

func (r *URLRepository) GetURLByID(
	ctx context.Context,
	id int64,
) (urlData *model.Url, error error) {

	var url model.Url

	query := `
			SELECT id, user_id, custom_alias, is_active, long_url, short_code FROM urls WHERE id = $1
		`
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&url.ID,
		&url.UserID,
		&url.CustomAlias,
		&url.IsActive,
		&url.LongURL,
		&url.ShortCode,
	)
	return &url, err
}

func (r *URLRepository) GetURLsByUserId(
	ctx context.Context,
	userId int64,
) ([]model.Url, error) {
	query := `
		SELECT * FROM urls WHERE user_id = $1 AND is_active = TRUE ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		userId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []model.Url

	for rows.Next() {
		var url model.Url

		err := rows.Scan(
			&url.ID,
			&url.ShortCode,
			&url.LongURL,
			&url.Title,
			&url.CustomAlias,
			&url.IsActive,
			&url.ExpiresAt,
			&url.CreatedAt,
			&url.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (r *URLRepository) DeactivateURL(
	ctx context.Context,
	urlID int64,
	userID int64,
) error {
	query := `
		UPDATE urls SET is_active = FALSE WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.Exec(
		query,
		urlID,
		userID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrURLNotFound
		}

		return ErrFailedToDeactiveURL
	}
	rows, err := result.RowsAffected()

	if rows == 0 {
		return ErrURLNotFound
	}

	return nil
}

func (r *URLRepository) ActivateURL(
	ctx context.Context,
	urlID int64,
) error {
	query := `
		UPDATE urls SET is_active = TRUE WHERE id = $1
	`

	result, err := r.db.Exec(
		query,
		urlID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrURLNotFound
		}

		return ErrFailedToActiveURL
	}
	rows, err := result.RowsAffected()

	if rows == 0 {
		return ErrURLNotFound
	}

	return nil
}

func (r *URLRepository) CheckURLUniqueness(
	ctx context.Context,
	customAlias string,
) (bool, error) {
	query := `
		SELECT EXISTS (SELECT 1 FROM urls WHERE custom_alias = $1)
	`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		customAlias,
	).Scan(&exists)

	if err != nil {
		return false, err
	}

	return !exists, nil
}

func (r *URLRepository) NextID(
	ctx context.Context,
) (int64, error) {
	var id int64

	err := r.db.QueryRowContext(
		ctx,
		`SELECT nextval('urls_id_seq')`,
	).Scan(&id)
	return id, err
}
