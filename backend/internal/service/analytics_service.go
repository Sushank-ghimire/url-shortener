package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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

func (s *AnalyticsService) GetAnalytics(
	ctx context.Context,
	userID int64,
	urlID int64,
	limit int64,
	offset int64,
) ([]model.Analytics, int, error) {

	url, err := s.urlRepo.GetURLByID(ctx, urlID)
	if err != nil {
		return nil, 0, err
	}

	if url.UserID != userID {
		return nil, 0, ErrAnalyticsUnauthorized
	}

	return s.analyticsRepo.Get(
		ctx,
		urlID,
		limit,
		offset,
	)
}

func (s *AnalyticsService) GetSummary(
	ctx context.Context,
	userID int64,
	urlID int64,
) (*model.AnalyticsSummary, error) {

	url, err := s.urlRepo.GetURLByID(ctx, urlID)
	if err != nil {
		return nil, err
	}

	if url.UserID != userID {
		return nil, ErrAnalyticsUnauthorized
	}

	return s.analyticsRepo.GetSummary(
		ctx,
		urlID,
	)
}

func (s *AnalyticsService) GetTimeline(
	ctx context.Context,
	userID int64,
	urlID int64,
	period string,
) ([]model.ClickTimeline, error) {
	url, err := s.urlRepo.GetURLByID(ctx, urlID)
	if err != nil {
		return nil, err
	}

	if url.UserID != userID {
		return nil, ErrAnalyticsUnauthorized
	}

	from, to, err := getTimelineRange(period)
	if err != nil {
		return nil, err
	}

	return s.analyticsRepo.GetTimeline(
		ctx,
		urlID,
		from,
		to,
	)
}

func getTimelineRange(period string) (time.Time, time.Time, error) {
	now := time.Now().UTC()

	switch period {
	case "7d":
		return now.AddDate(0, 0, -7), now, nil

	case "30d":
		return now.AddDate(0, 0, -30), now, nil

	case "90d":
		return now.AddDate(0, 0, -90), now, nil

	case "1y":
		return now.AddDate(-1, 0, 0), now, nil

	default:
		return time.Time{}, time.Time{},
			fmt.Errorf("invalid period: %s", period)
	}
}
