package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Sushank-ghimire/url-shortener/internal/model"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExist = errors.New("email already exists")
	UserLoggedOut        = errors.New("user must be logged in to get profile")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	query := `
		INSERT INTO users (
			name,
			email,
			pass_hash,
			is_verified
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, created_at, updated_at
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PassHash,
		user.IsVerified,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	query := `
		SELECT id, name, email, pass_hash, is_verified, created_at, updated_at
		FROM users WHERE email = $1
	`
	var user model.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PassHash,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	query := `
		SELECT id, name, email, pass_hash, is_verified, created_at, updated_at FROM users
		WHERE id = $1
	`

	var user model.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PassHash,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) MarkEmailVerified(
	ctx context.Context,
	userID int64,
) error {
	query := `
		UPDATE users SET is_verified = TRUE, updated_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return fmt.Errorf("mark email as verified: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check email verification update: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %d", userID)
	}
	return nil
}

func (r *UserRepository) CheckUserValidity(
	ctx context.Context,
	userId int64,
) (isLoggedIn bool, error error) {
	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 FROM refresh_tokens WHERE user_id = $1 AND revoked_at IS NULL
			AND expires_at > NOW()
		)
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&exists,
	)

	if err != nil {
		return false, err
	}

	if !exists {
		return false, UserLoggedOut
	}
	return true, nil
}
