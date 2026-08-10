package config

import "os"

type Config struct {
	AppEnv     string
	Port       string
	DBPort     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Security Keys
	JWTSecret          string
	AccessTokenExpiry  string
	RefreshTokenExpiry string

	// Cors
	CORSOrigins string
}

func LoadConfig() Config {
	return Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		Port:               getEnv("PORT", ":8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "url_shortener"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		JWTSecret:          getEnv("JWT_SECRET_KEY", "jwt-secret-key"),
		AccessTokenExpiry:  getEnv("ACCESS_TOKEN_EXPIRY", "15m"),
		RefreshTokenExpiry: getEnv("REFRESH_TOKEN_EXPIRY", "168h"),
		CORSOrigins:        getEnv("CORS_ORIGINS", "http://localhost:3000"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
