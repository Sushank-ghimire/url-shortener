package routes

import (
	"github.com/Sushank-ghimire/url-shortener/internal/auth"
	"github.com/Sushank-ghimire/url-shortener/internal/handler"
	"github.com/Sushank-ghimire/url-shortener/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	router *gin.Engine,
	authHandler *handler.AuthHandler,
	verificationHandler *handler.VerificationHandler,
	urlHandler *handler.URLHandler,
	analyticsHandler *handler.AnalyticsHandler,
	jwtService *auth.JWTService,
) {
	api := router.Group("/api/v1")

	setupAuthRoutes(
		api,
		authHandler,
		jwtService,
	)

	setupVerificationRoutes(
		api,
		verificationHandler,
	)

	setupURLRoutes(
		api,
		urlHandler,
		analyticsHandler,
		jwtService,
	)
}

func setupAuthRoutes(
	api *gin.RouterGroup,
	h *handler.AuthHandler,
	jwtService *auth.JWTService,
) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
	protected := auth.Group("")
	protected.Use(middleware.AuthRequired(jwtService))
	{
		protected.POST("/refresh", h.Refresh)
		protected.POST("/logout", h.Logout)
		protected.POST("/logout-all", h.LogoutAll)
		protected.GET("/me", h.GetMe)
	}
}

func setupVerificationRoutes(
	api *gin.RouterGroup,
	h *handler.VerificationHandler,
) {
	email := api.Group("/email")
	{

		email.POST("/verify", h.VerifyEmail)
	}
}

func setupURLRoutes(
	api *gin.RouterGroup,
	urlHandler *handler.URLHandler,
	analyticsHandler *handler.AnalyticsHandler,
	jwtService *auth.JWTService,
) {
	url := api.Group("/url")

	// Public
	url.GET("/:id", urlHandler.GetURLByID)
	url.GET("/:shortcode/resolve", urlHandler.ResolveURL)

	// Protected
	protected := url.Group("")
	protected.Use(middleware.AuthRequired(jwtService))
	{
		protected.POST("/", urlHandler.ShortURL)

		analytics := protected.Group("/:id/analytics")
		{
			analytics.GET("", analyticsHandler.GetAnalytics)
			analytics.GET("/summary", analyticsHandler.GetSummary)
			analytics.GET("/timeline", analyticsHandler.GetTimeline)
		}
	}
}
