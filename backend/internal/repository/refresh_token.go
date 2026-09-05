package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Sushank-ghimire/url-shortener/internal/model"
)

type RefreshTokenRepository struct {
	db *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: db,
	}
}

func (r *RefreshTokenRepository) Create(
	ctx context.Context,
	token *model.RefreshToken,
) error {
	query := `
		INSERT INTO refresh_tokens (
			user_id,
			token_hash,
			expires_at
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRowContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt).Scan(
		&token.ID,
		&token.CreatedAt,
	)
}

func (r *RefreshTokenRepository) FindByHash(
	ctx context.Context,
	tokenHash string,
) (*model.RefreshToken, error) {

	query := `
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			created_at,
			revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`

	var token model.RefreshToken

	err := r.db.QueryRowContext(
		ctx,
		query,
		tokenHash,
	).Scan(
		&token.ID,
		&token.UserID,
		&token.TokenHash,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("refresh token not found")
	}

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *RefreshTokenRepository) Revoke(
	ctx context.Context,
	id int64,
) error {

	now := time.Now()

	query := `
		UPDATE refresh_tokens SET revoked_at = $1
		WHERE id = $2 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		now,
		id,
	)

	return err
}

func (r *RefreshTokenRepository) RevokeAllByUserID(
	ctx context.Context,
	userID int64,
) error {
	_, err := r.db.ExecContext(
		ctx,
		`
			UPDATE refresh_tokens SET revoked_at = NOW() WHERE user_id = $1 AND revoked_at IS NULL
		`,
		userID,
	)
	return err
}
