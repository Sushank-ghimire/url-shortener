package service

import (
	"context"
	"errors"
	"log"

	"github.com/Sushank-ghimire/url-shortener/internal/model"
	"github.com/Sushank-ghimire/url-shortener/internal/repository"
)

type AnalyticsService struct {
	analyticsRepo *repository.AnalyticsRepository
	urlRepo       *repository.URLRepository
	geoIPService  *GeoIPService
}

var (
	ErrAnalyticsUnauthorized = errors.New("user does not own this url")
	ErrInvalidAnalyticsData  = errors.New("analytics data is required")
)

func NewAnalyticsService(
	analyticsRepo *repository.AnalyticsRepository,
	urlRepo *repository.URLRepository,
	geoIPService *GeoIPService,
) *AnalyticsService {
	return &AnalyticsService{
		analyticsRepo: analyticsRepo,
		urlRepo:       urlRepo,
		geoIPService:  geoIPService,
	}
}

func (s *AnalyticsService) Create(
	ctx context.Context,
	analytics *model.Analytics,
) error {
	if analytics == nil {
		return ErrInvalidAnalyticsData
	}

	country, city, err := s.geoIPService.Lookup(
		analytics.IpAddress,
	)

	log.Print(s.geoIPService.Lookup(analytics.IpAddress))

	if err != nil {
		analytics.Country = country
		analytics.City = city
	}

	return s.analyticsRepo.Create(ctx, analytics)
}

func (s *AnalyticsService) Get(
	ctx context.Context,
	userID int64,
	urlID int64,
	limit int64,
	offset int64,
) ([]model.Analytics, error) {

	url, err := s.urlRepo.GetURLByID(ctx, urlID)
	if err != nil {
		return nil, err
	}

	if url.UserID != userID {
		return nil, ErrAnalyticsUnauthorized
	}

	return s.analyticsRepo.Get(
		ctx,
		urlID,
		limit,
		offset,
	)
}
