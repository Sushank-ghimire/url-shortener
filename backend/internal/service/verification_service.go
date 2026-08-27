package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/Sushank-ghimire/url-shortener/internal/auth"
	"github.com/Sushank-ghimire/url-shortener/internal/config"
	"github.com/Sushank-ghimire/url-shortener/internal/encoder"
	"github.com/Sushank-ghimire/url-shortener/internal/model"
	"github.com/Sushank-ghimire/url-shortener/internal/repository"
)

type VerificationService struct {
	userRepo  *repository.UserRepository
	tokenRepo *repository.VerificationRepository
	config    *config.Config
}

var (
	TokenExpired       = errors.New("token is already expired")
	TokenAlreadyUsed   = errors.New("token is already used")
	ErrFailedEmailSent = errors.New("failed to send email")
	ErrFailedToAdd     = errors.New("failed to insert token in db")
	ErrTokenGeneration = errors.New("failed to generate verification token")
	ErrInvalidToken    = errors.New("invalid token provided")
)

func NewVerificationService(
	userRepo *repository.UserRepository,
	tokenRepo *repository.VerificationRepository,
	config *config.Config,
) *VerificationService {
	return &VerificationService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		config:    config,
	}
}

func (r *VerificationService) SendVerificationEmail(
	ctx context.Context,
	user *model.User,
) error {
	token, err := encoder.GenerateVerificationToken()

	if err != nil {
		return ErrTokenGeneration
	}
	hashedToken := encoder.HashToken(token)
	expiresAt := time.Now().Add(10 * time.Minute)

	verificationToken := model.EmailVerificationToken{
		UserID:    user.ID,
		TokenHash: hashedToken,
		Purpose:   "email_verification",
		ExpiresAt: expiresAt,
	}

	err = r.tokenRepo.AddUserToken(ctx, verificationToken)
	if err != nil {
		return ErrFailedToAdd
	}

	verificationURL := fmt.Sprintf(
		"%s/auth/verify-email?token=%s",
		r.config.FrontendURL,
		url.QueryEscape(token),
	)

	sentId, err := auth.SendVerificationEmail(
		user.Email,
		user.Name,
		verificationURL,
		r.config,
	)
	if err != nil || sentId == "" {
		return ErrFailedEmailSent
	}
	return nil
}

func (r *VerificationService) VerifyEmail(
	ctx context.Context,
	token string,
) error {
	hashedToken := encoder.HashToken(token)

	dbToken, err := r.tokenRepo.GetValidTokenByHash(
		ctx,
		hashedToken,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrInvalidToken
		}
		return fmt.Errorf(
			"get verification token: %w",
			err,
		)
	}

	now := time.Now()

	if dbToken.ExpiresAt.Before(now) {
		return TokenExpired
	}

	if dbToken.UsedAt.Valid {
		return TokenAlreadyUsed
	}

	user, err := r.userRepo.FindByID(
		ctx,
		dbToken.UserID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return repository.ErrUserNotFound
		}

		return fmt.Errorf("find user: %w", err)
	}

	if err := r.userRepo.MarkEmailVerified(ctx, user.ID); err != nil {
		return fmt.Errorf(
			"mark email as verified: %w",
			err,
		)
	}

	if err := r.tokenRepo.MarkUsed(ctx, dbToken.ID); err != nil {
		return fmt.Errorf(
			"mark verification token as used: %w",
			err,
		)
	}

	return nil
}
