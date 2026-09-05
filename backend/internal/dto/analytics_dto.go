package dto

import (
	"github.com/Sushank-ghimire/url-shortener/internal/model"
)

type TimelineResponse struct {
	Timeline []model.ClickTimeline `json:"timeline"`
}

type PaginationRequest struct {
	Limit  int64 `form:"limit"`
	Offset int64 `form:"offset"`
}
