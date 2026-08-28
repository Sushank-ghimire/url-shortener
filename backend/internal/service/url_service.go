package service

import (
	"context"
	"errors"
	"time"

	"github.com/Sushank-ghimire/url-shortener/internal/dto"
	"github.com/Sushank-ghimire/url-shortener/internal/encoder"
	"github.com/Sushank-ghimire/url-shortener/internal/model"
	"github.com/Sushank-ghimire/url-shortener/internal/repository"
)

type URLService struct {
	urlRepo  *repository.URLRepository
	userRepo *repository.UserRepository
}

var (
	ErrURLInactive  = errors.New("provided url is not active")
	ErrURLExpired   = errors.New("provided url is expired")
	ErrUnauthorized = errors.New("user does not own this url")
	ErrCustomAlias  = errors.New("custom alias is not supported")
)

func NewURLService(
	urlRepo *repository.URLRepository,
	userRepo *repository.UserRepository,
) *URLService {
	return &URLService{
		urlRepo:  urlRepo,
		userRepo: userRepo,
	}
}

func (s *URLService) CreateURL(
	ctx context.Context,
	userID int64,
	url *model.Url,
) (*dto.URLResponse, error) {

	if url.CustomAlias != "" {
		url.UserID = userID
		url.ShortCode = url.CustomAlias

		exists, err := s.urlRepo.CheckURLUniqueness(ctx, url.CustomAlias)

		if err != nil || exists {
			return nil, ErrCustomAlias
		}

		if err := s.urlRepo.Create(ctx, url); err != nil {
			return nil, ErrCustomAlias
		}

		return toURLResponse(url), nil
	}

	if url == nil {
		return nil, errors.New("url data is required")
	}

	nextID, err := s.urlRepo.NextID(ctx)
	if err != nil {
		return nil, err
	}

	url.UserID = userID
	url.ShortCode = encoder.Encode(nextID)

	err = s.urlRepo.Create(ctx, url)
	if err != nil {
		return nil, err
	}

	return toURLResponse(url), nil
}

func (s *URLService) RedirectUser(
	ctx context.Context,
	shortCode string,
) (*model.Url, error) {
	url, err := s.urlRepo.GetURLByShortCode(ctx, shortCode)
	if err != nil {
		return nil, err
	}

	if !url.IsActive {
		return nil, ErrURLInactive
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return nil, ErrURLExpired
	}

	return url, nil
}

func (s *URLService) GetURLByID(
	ctx context.Context,
	id int64,
) (*dto.URLResponse, error) {
	url, err := s.urlRepo.GetURLByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toURLResponse(url), nil
}

func (s *URLService) GetUserURLs(
	ctx context.Context,
	userID int64,
) ([]dto.URLResponse, error) {
	urls, err := s.urlRepo.GetURLsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.URLResponse, 0, len(urls))

	for _, url := range urls {
		responses = append(responses, *toURLResponse(&url))
	}

	return responses, nil
}

func (s *URLService) DeleteURL(
	ctx context.Context,
	userID int64,
	urlID int64,
) error {
	return s.urlRepo.DeactivateURL(ctx, urlID, userID)
}

func (s *URLService) ResolveURL(
	ctx context.Context,
	shortcode string,
) (*model.Url, error) {
	url, err := s.urlRepo.GetURLByShortCode(ctx, shortcode)
	if err != nil {
		return nil, err
	}
	return url, err
}

func toURLResponse(url *model.Url) *dto.URLResponse {
	return &dto.URLResponse{
		ID:          url.ID,
		ShortCode:   url.ShortCode,
		LongURL:     url.LongURL,
		Title:       url.Title,
		CustomAlias: url.CustomAlias,
		IsActive:    url.IsActive,
		ExpiresAt:   url.ExpiresAt,
		CreatedAt:   url.CreatedAt,
		UpdatedAt:   url.UpdatedAt,
	}
}
