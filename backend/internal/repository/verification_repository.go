package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sushank-ghimire/url-shortener/internal/model"
)

type VerificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(
	db *sql.DB,
) *VerificationRepository {
	return &VerificationRepository{
		db: db,
	}
}

func (r *VerificationRepository) AddUserToken(
	ctx context.Context,
	token model.EmailVerificationToken,
) error {
	query := `
		INSERT INTO email_verification_tokens (
			user_id,
			token_hash,
			purpose,
			expires_at
		)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		token.UserID,
		token.TokenHash,
		token.Purpose,
		token.ExpiresAt,
	)

	if err != nil {
		return fmt.Errorf("create verification token: %2", err)
	}
	return nil
}

func (r *VerificationRepository) GetValidTokenByHash(
	ctx context.Context,
	hashToken string,
) (token *model.EmailVerificationToken, error error) {
	var verificationToken model.EmailVerificationToken
	query := `
		SELECT id, user_id, purpose, expires_at, used_at, created_at
	 	FROM email_verification_tokens WHERE token_hash = $1 AND used_at IS NULL
		AND expires_at > NOW()
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		hashToken,
	).Scan(
		&verificationToken.ID,
		&verificationToken.UserID,
		&verificationToken.Purpose,
		&verificationToken.ExpiresAt,
		&verificationToken.UsedAt,
		&verificationToken.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid token: %w", err)
		}
		return nil, fmt.Errorf("get token by hash: %w", err)
	}

	return &verificationToken, nil
}

func (r *VerificationRepository) MarkUsed(
	ctx context.Context,
	tokenID int64,
) error {
	query := `
		UPDATE email_verification_tokens SET used_at = NOW() WHERE id = $1 AND used_at IS NULL
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		tokenID,
	)
	if err != nil {
		return fmt.Errorf("mark verification token used: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *VerificationRepository) InvalidateUserTokens(
	ctx context.Context,
	userID int64,
	purpose string,
) error {
	query := `
		UPDATE email_verification_tokens SET used_at = NOW()
		WHERE user_id = $1 AND purpose = $2 AND used_at IS NULL
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		purpose,
	)
	if err != nil {
		return fmt.Errorf("invalidate verification tokens: %w", err)
	}

	return nil
}
