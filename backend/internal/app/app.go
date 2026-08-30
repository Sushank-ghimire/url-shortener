package app

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Sushank-ghimire/url-shortener/internal/auth"
	"github.com/Sushank-ghimire/url-shortener/internal/config"
	"github.com/Sushank-ghimire/url-shortener/internal/database"
	"github.com/Sushank-ghimire/url-shortener/internal/handler"
	"github.com/Sushank-ghimire/url-shortener/internal/middleware"
	"github.com/Sushank-ghimire/url-shortener/internal/repository"
	"github.com/Sushank-ghimire/url-shortener/internal/routes"
	"github.com/Sushank-ghimire/url-shortener/internal/service"
	"github.com/Sushank-ghimire/url-shortener/migrations"

	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
)

type App struct {
	config *config.Config
	db     *sql.DB
	geoDB  *geoip2.Reader
	router *gin.Engine
}

var (
	GeoLiteCityLoc    = "data/Geolite/GeoLite2-City.mmdb"
	GeoLiteCountryLoc = "data/Geolite/GeoLite2-Country.mmdb"
)

func New() (*App, error) {
	cfg := config.LoadConfig()

	db := database.Connect(cfg)

	if err := migrations.MigrateAll(db); err != nil {
		db.Close()
		return nil, err
	}

	var geoDB *geoip2.Reader

	geoDB, err := geoip2.Open(GeoLiteCityLoc)

	if err != nil {
		log.Printf("GeoIP database not found, GeoIP disabled: %v", err)
	}

	app := &App{
		config: &cfg,
		db:     db,
		geoDB:  geoDB,
	}

	if err := app.setup(); err != nil {
		app.Close()
		return nil, err
	}

	return app, nil
}

func (a *App) setup() error {
	jwtService := auth.NewJWTService(
		a.config.JWTSecret,
	)

	accessTokenDuration, err := time.ParseDuration(
		a.config.AccessTokenExpiry,
	)
	if err != nil {
		return err
	}

	refreshTokenDuration, err := time.ParseDuration(
		a.config.RefreshTokenExpiry,
	)
	if err != nil {
		return err
	}

	userRepo := repository.NewUserRepository(a.db)

	refreshTokenRepo := repository.NewRefreshTokenRepository(
		a.db,
	)

	verificationTokenRepo := repository.NewVerificationRepository(
		a.db,
	)

	urlRepo := repository.NewURLRepository(a.db)

	analyticsRepo := repository.NewAnalyticsRepository(a.db)

	authService := service.NewAuthService(
		userRepo,
		refreshTokenRepo,
		jwtService,
		accessTokenDuration,
		refreshTokenDuration,
	)

	verificationService := service.NewVerificationService(
		userRepo,
		verificationTokenRepo,
		a.config,
	)

	urlService := service.NewURLService(
		urlRepo,
		userRepo,
	)

	geoIPService := service.NewGeoIPService(
		a.geoDB,
	)

	analyticsService := service.NewAnalyticsService(
		analyticsRepo,
		urlRepo,
		geoIPService,
	)

	authHandler := handler.NewAuthHandler(
		authService,
		verificationService,
	)

	verificationHandler := handler.NewVerificationHandler(
		verificationService,
	)

	urlHandler := handler.NewURLHandler(
		urlService,
		analyticsService,
	)

	analyticsHandler := handler.NewAnalyticsHandler(
		analyticsService,
		urlService,
	)

	a.setupRouter(
		jwtService,
		authHandler,
		verificationHandler,
		urlHandler,
		analyticsHandler,
	)

	return nil
}

func (a *App) setupRouter(
	jwtService *auth.JWTService,
	authHandler *handler.AuthHandler,
	verificationHandler *handler.VerificationHandler,
	urlHandler *handler.URLHandler,
	analyticsHandler *handler.AnalyticsHandler,
) {
	router := gin.Default()

	router.Use(
		middleware.CORS(
			strings.Split(a.config.CORSOrigins, ","),
		),
	)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Url shortener backend server is live",
		})
	})

	routes.SetupRoutes(
		router,
		authHandler,
		verificationHandler,
		urlHandler,
		analyticsHandler,
		jwtService,
	)

	a.router = router
}

func (a *App) Run() error {
	log.Printf(
		"server running on %s",
		a.config.Port,
	)

	return a.router.Run(a.config.Port)
}

func (a *App) Close() {
	if a.geoDB != nil {
		a.geoDB.Close()
	}

	if a.db != nil {
		a.db.Close()
	}
}
