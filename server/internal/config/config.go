package config

import (
	"os"
	"strings"
)

// Config holds application configuration
type Config struct {
	DatabaseURL  string
	JWTSecret    string
	Environment  string
	AllowOrigins []string
	RateLimitRPS int
}

// Load loads configuration from environment variables
func Load() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Generate a warning if no JWT secret is set
		jwtSecret = "change-this-secret-in-production-use-a-long-random-string"
	}

	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	origins := []string{"http://localhost:5173", "http://localhost:3000"}
	if allowOrigins != "" {
		origins = strings.Split(allowOrigins, ",")
	}

	return &Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    jwtSecret,
		Environment:  getEnv("ENVIRONMENT", "development"),
		AllowOrigins: origins,
		RateLimitRPS: 10, // 10 requests per second per IP
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
