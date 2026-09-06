package handler

import (
	"errors"
	"net/http"

	"github.com/Sushank-ghimire/url-shortener/internal/dto"
	"github.com/Sushank-ghimire/url-shortener/internal/middleware"
	"github.com/Sushank-ghimire/url-shortener/internal/model"
	"github.com/Sushank-ghimire/url-shortener/internal/repository"
	"github.com/Sushank-ghimire/url-shortener/internal/service"
	"github.com/Sushank-ghimire/url-shortener/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService        *service.AuthService
	verficationService *service.VerificationService
}

var (
	InvalidRequest      = "invalid request (all fields are required)"
	ErrEmailNotVerified = "user email is not verified"
)

func NewAuthHandler(
	authService *service.AuthService,
	verificationService *service.VerificationService,
) *AuthHandler {
	return &AuthHandler{
		authService:        authService,
		verficationService: verificationService,
	}
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, InvalidRequest)
		return
	}

	user, err := h.authService.Register(
		c.Request.Context(),
		req,
	)

	if user != nil {
		h.verficationService.SendVerificationEmail(
			c.Request.Context(),
			&model.User{
				ID:    user.ID,
				Email: user.Email,
				Name:  user.Name,
			},
		)
	}

	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			utils.Error(c, http.StatusConflict, "email already exists")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "failed to register user")
		return
	}

	utils.Success(c, http.StatusCreated, "verification link is sent to the user's mail", map[string]any{
		"user": user,
	})
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, InvalidRequest)
		return
	}
	response, err := h.authService.Login(
		c.Request.Context(),
		req,
	)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			utils.Error(c, http.StatusUnauthorized, "invalid email or password")
			return
		}
		if errors.Is(err, service.ErrEmailUnverified) {
			utils.Error(c, http.StatusUnauthorized, "user email is not verified")
			return
		}
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.Error(c, http.StatusUnauthorized, "user not found")
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.Error(c, http.StatusUnauthorized, "invalid credentials")
			return
		}

		utils.Error(c, http.StatusInternalServerError, "failed to login user")
		return
	}
	utils.Success(c, http.StatusOK, "login successful", response)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, InvalidRequest)
		return
	}

	response, err := h.authService.Refresh(
		c.Request.Context(),
		req,
	)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRefreshToken),
			errors.Is(err, service.ErrRefreshTokenExpired),
			errors.Is(err, service.ErrRefreshTokenRevoked):
			utils.Error(c, http.StatusUnauthorized, "invalid or expired refresh token")
		default:
			utils.Error(c, http.StatusInternalServerError, "failed to refresht token")
		}
		return
	}

	utils.Success(c, http.StatusOK, "token refreshed", response)
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadGateway, "invalid request")
		return
	}

	err := h.authService.Logout(
		c.Request.Context(),
		req.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) {
			utils.Error(c, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		if errors.Is(err, service.ErrRefreshTokenRevoked) {
			utils.Error(c, http.StatusUnauthorized, "refresh token revoked")
			return
		}
		if errors.Is(err, service.ErrRefreshTokenExpired) {
			utils.Error(c, http.StatusUnauthorized, "refresh token expired")
			return
		}

		utils.Error(c, http.StatusInternalServerError, "failed to logout")
		return
	}
	utils.Success(c, http.StatusOK, "logged out successfully", nil)
}

// GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *gin.Context) {
	value, exists := c.Get(middleware.UserIDKey)

	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	userId, ok := value.(int64)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "invalid user context")
		return
	}

	user, err := h.authService.GetUserByID(
		c.Request.Context(),
		userId,
	)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.Error(c, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, repository.UserLoggedOut) {
			utils.Error(c, http.StatusNotFound, "user is revoked login to access profile")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "failed to get user")
		return
	}
	utils.Success(c, http.StatusOK, "user fetched successfully", user)
}

// POST /api/v1/auth/logout-all
func (h *AuthHandler) LogoutAll(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		utils.Error(
			c,
			http.StatusUnauthorized,
			"authentication required",
		)
		return
	}
	err := h.authService.LogoutAll(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to logout from all devices")
		return
	}
	utils.Success(c, http.StatusOK, "logged out from all devices successfully", nil)
}
