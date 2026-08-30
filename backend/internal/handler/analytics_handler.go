package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Sushank-ghimire/url-shortener/internal/dto"
	"github.com/Sushank-ghimire/url-shortener/internal/service"
	"github.com/Sushank-ghimire/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
	urlService       *service.URLService
}

func NewAnalyticsHandler(
	analyticsService *service.AnalyticsService,
	urlService *service.URLService,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		urlService:       urlService,
	}
}

func (h *AnalyticsHandler) GetAnalytics(
	c *gin.Context,
) {

	var pagination dto.PaginationRequest

	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.Error(
			c,
			http.StatusBadRequest,
			"invalid pagination parameters",
		)
		return
	}

	if pagination.Limit <= 0 {
		pagination.Limit = 20
	}

	if pagination.Limit > 100 {
		pagination.Limit = 100
	}

	if pagination.Offset < 0 {
		pagination.Offset = 0
	}

	urlIDParam := c.Param("id")

	if urlIDParam == "" {
		utils.Error(
			c,
			http.StatusBadRequest,
			"url id is missing",
		)
		return
	}

	urlID, err := strconv.ParseInt(urlIDParam, 10, 64)

	if err != nil {
		utils.Error(
			c,
			http.StatusBadRequest,
			"invalid url id",
		)
		return
	}

	userID, ok := getUserID(c)
	if !ok {
		utils.Error(
			c,
			http.StatusUnauthorized,
			"analytics access denied",
		)
		return
	}
	analytics, err := h.analyticsService.Get(
		c.Request.Context(),
		userID,
		urlID,
		pagination.Limit,
		pagination.Offset,
	)

	if err != nil {
		if errors.Is(err, service.ErrAnalyticsUnauthorized) {
			utils.Error(
				c,
				http.StatusForbidden,
				"you don't own this url",
			)
			return
		}
		utils.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve analytics",
		)
		return
	}
	utils.Success(
		c,
		http.StatusOK,
		"analytics retrieved successfully",
		gin.H{
			"analytics": analytics,
			"pagination": gin.H{
				"limit":  pagination.Limit,
				"offset": pagination.Offset,
			},
		},
	)
}
