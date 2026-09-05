package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/Sushank-ghimire/url-shortener/internal/auth"
	"github.com/Sushank-ghimire/url-shortener/internal/dto"
	"github.com/Sushank-ghimire/url-shortener/internal/model"
	"github.com/Sushank-ghimire/url-shortener/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrEmailExists         = errors.New("email already exists")
	ErrEmailUnverified     = errors.New("email is not verified")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrRefreshTokenRevoked = errors.New("refresh token revoked")
)

type AuthService struct {
	userRepo         *repository.UserRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	jwtService       *auth.JWTService

	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewAuthService(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	jwtService *auth.JWTService,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		refreshTokenRepo:     refreshTokenRepo,
		jwtService:           jwtService,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (*dto.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// check whether user already exists
	_, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return nil, ErrEmailExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Hash password.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:       strings.TrimSpace(req.Name),
		Email:      email,
		PassHash:   string(passwordHash),
		IsVerified: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return &dto.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*dto.AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	user, err := s.userRepo.FindByEmail(ctx, email)

	if user == nil {
		return nil, repository.ErrUserNotFound
	}

	if !user.IsVerified {
		return nil, ErrEmailUnverified
	}

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	// Compare plain password against stored bcrypt hash.
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PassHash),
		[]byte(req.Password),
	); err != nil {
		return nil, ErrInvalidCredentials
	}

	accessToken, err := s.jwtService.GenerateAccessToken(
		user.ID,
		s.accessTokenDuration,
	)

	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	tokenHash := hashRefreshToken(refreshToken)

	refreshTokenModel := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenModel); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			IsVerified: user.IsVerified,
		},
	}, nil
}

func (s *AuthService) Refresh(
	ctx context.Context,
	req dto.RefreshRequest,
) (*dto.AuthResponse, error) {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return nil, ErrInvalidRefreshToken
	}

	// Hash the token supplied by the client.
	tokenHash := hashRefreshToken(req.RefreshToken)

	storedToken, err := s.refreshTokenRepo.FindByHash(
		ctx,
		tokenHash,
	)

	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Check whether token has already been revoked.
	if storedToken.RevokedAt != nil {
		return nil, ErrRefreshTokenRevoked
	}

	// Check expiration.
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	// Get associated user.
	user, err := s.userRepo.FindByID(
		ctx,
		storedToken.UserID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidRefreshToken
		}

		return nil, err
	}

	// Generate new access token.
	accessToken, err := s.jwtService.GenerateAccessToken(
		user.ID,
		s.accessTokenDuration,
	)

	if err != nil {
		return nil, err
	}

	// Rotate refresh token.
	newRefreshToken, err := generateRefreshToken()

	if err != nil {
		return nil, err
	}

	newRefreshTokenHash := hashRefreshToken(newRefreshToken)

	newToken := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(s.refreshTokenDuration),
	}

	// Revoke old refresh token.
	if err := s.refreshTokenRepo.Revoke(
		ctx,
		storedToken.ID,
	); err != nil {
		return nil, err
	}

	// Store new refresh token.
	if err := s.refreshTokenRepo.Create(
		ctx,
		newToken,
	); err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: dto.UserResponse{
			ID:         user.ID,
			Name:       user.Name,
			Email:      user.Email,
			IsVerified: user.IsVerified,
		},
	}, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	refreshToken = strings.TrimSpace(refreshToken)

	if refreshToken == "" {
		return ErrInvalidRefreshToken
	}

	tokenHash := hashRefreshToken(refreshToken)
	token, err := s.refreshTokenRepo.FindByHash(
		ctx,
		tokenHash,
	)

	if err != nil {
		return ErrInvalidRefreshToken
	}
	if token.RevokedAt != nil {
		return ErrRefreshTokenRevoked
	}

	return s.refreshTokenRepo.Revoke(
		ctx,
		token.ID,
	)
}

func (s *AuthService) LogoutAll(
	ctx context.Context,
	userID int64,
) error {
	return s.refreshTokenRepo.RevokeAllByUserID(
		ctx,
		userID,
	)
}

func generateRefreshToken() (string, error) {
	// 32 random bytes = 256 bits of entropy.
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(tokenBytes), nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) GetUserByID(
	ctx context.Context,
	userID int64,
) (*dto.UserResponse, error) {

	isLoggedIn, err := s.userRepo.CheckUserValidity(ctx, userID)

	if err != nil {
		return nil, err
	}

	if !isLoggedIn {
		return nil, repository.UserLoggedOut
	}

	user, err := s.userRepo.FindByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}, nil
}
