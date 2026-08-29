package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Sushank-ghimire/url-shortener/internal/dto"
	"github.com/Sushank-ghimire/url-shortener/internal/model"
	"github.com/Sushank-ghimire/url-shortener/internal/service"
	"github.com/Sushank-ghimire/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	urlService       *service.URLService
	analyticsService *service.AnalyticsService
}

func NewURLHandler(
	urlService *service.URLService,
	analyticsService *service.AnalyticsService,
) *URLHandler {
	return &URLHandler{
		urlService:       urlService,
		analyticsService: analyticsService,
	}
}

func (h *URLHandler) ShortURL(c *gin.Context) {
	var req dto.ShortURLRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		fmt.Println(req)
		utils.Error(
			c,
			http.StatusBadRequest,
			"invalid request all fields are required",
		)
		return
	}

	userID, ok := getUserID(c)

	if !ok {
		utils.Error(
			c,
			http.StatusUnauthorized,
			"user not authenticated",
		)
		return
	}

	url, err := h.urlService.CreateURL(
		c.Request.Context(),
		userID,
		&model.Url{
			CustomAlias: req.CustomAlias,
			Title:       req.Title,
			LongURL:     req.LongURL,
			ExpiresAt:   req.ExpiresAt,
			UserID:      userID,
		},
	)

	if err != nil || url == nil {
		fmt.Print("failed to short url: ", err.Error())
		utils.Error(
			c,
			http.StatusInternalServerError,
			"failed to short the url",
		)
		return
	}

	utils.Success(
		c,
		http.StatusCreated,
		"url shorted successfully",
		url,
	)
}

func (h *URLHandler) GetURLByID(c *gin.Context) {
	shortCode := c.Param("id")

	if shortCode == "" {
		utils.Error(
			c,
			http.StatusBadRequest,
			"short code is not provided",
		)
		return
	}

	url, err := h.urlService.RedirectUser(
		c.Request.Context(),
		shortCode,
	)

	if err != nil {
		if errors.Is(err, service.ErrURLExpired) {
			utils.Error(
				c,
				http.StatusNotFound,
				"url is expired",
			)
			return
		}
		if errors.Is(err, service.ErrURLInactive) {
			utils.Error(
				c,
				http.StatusNotFound,
				"url is deactivated",
			)
			return
		}

		utils.Error(
			c,
			http.StatusBadRequest,
			"internal server error",
		)
		return
	}

	analytics := buildAnalytics(
		c,
		url.ID,
	)

	_ = h.analyticsService.Create(
		c.Request.Context(),
		analytics,
	)

	c.Redirect(http.StatusFound, url.LongURL)
}

func (h *URLHandler) ResolveURL(c *gin.Context) {
	shortCode := c.Param("shortcode")
	if shortCode == "" {
		utils.Error(
			c,
			http.StatusBadRequest,
			"shortcode is required",
		)
		return
	}
	url, err := h.urlService.ResolveURL(
		c,
		shortCode,
	)
	if err != nil {
		utils.Error(
			c,
			http.StatusNotFound,
			"URL not found",
		)
		return
	}

	if !url.IsActive {
		utils.Error(
			c,
			http.StatusGone,
			"URL is inactive",
		)
		return
	}
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		utils.Error(
			c,
			http.StatusGone,
			"URL is expired",
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"exists": true,
		},
	)
}
