package handler

import (
	"errors"
	"net/http"

	"github.com/Sushank-ghimire/url-shortener/internal/dto"
	"github.com/Sushank-ghimire/url-shortener/internal/service"
	"github.com/Sushank-ghimire/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	verificationService *service.VerificationService
}

func NewVerificationHandler(
	verificationService *service.VerificationService,
) *VerificationHandler {
	return &VerificationHandler{
		verificationService: verificationService,
	}
}

// POST /api/v1/email/verify?token=token
func (h *VerificationHandler) VerifyEmail(
	c *gin.Context,
) {
	var req dto.VerificationRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.Error(
			c,
			http.StatusBadRequest,
			"verification token is required",
		)
		return
	}
	err := h.verificationService.VerifyEmail(
		c.Request.Context(),
		req.Token,
	)

	if err != nil {
		if errors.Is(err, service.TokenExpired) {
			utils.Error(
				c,
				http.StatusBadRequest,
				"verification token is expired",
			)
			return
		}
		if errors.Is(err, service.TokenAlreadyUsed) {
			utils.Error(
				c,
				http.StatusBadRequest,
				"verification token is already used",
			)
			return
		}

		if errors.Is(err, service.ErrInvalidToken) {
			utils.Error(
				c,
				http.StatusBadRequest,
				"invalid token provided",
			)
			return
		}

		utils.Error(
			c,
			http.StatusBadRequest,
			"invalid token provided",
		)
		return
	}
	utils.Success(
		c,
		http.StatusOK,
		"email is successfully verified",
		nil,
	)
}
