// femabras/backend/internal/config/config.go
package config

import (
	"log"
	"os"
)

type Config struct {
	Port               string
	DatabaseURL        string
	JWTSecret          string
	FrontendURL        string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	AppName            string
	SMTPHost           string
	SMTPPort           string
	SMTPUsername       string
	SMTPPassword       string
	FromEmail          string
	RedisURL           string `env:"REDIS_URL"`
}

func Load() Config {
	cfg := Config{
		Port:        getEnvOrDefault("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		FrontendURL: getEnvOrDefault("FRONTEND_URL", "http://localhost:3000"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required (32+ chars)")
	}

	cfg.GoogleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	cfg.GoogleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	cfg.GoogleRedirectURL = os.Getenv("GOOGLE_REDIRECT_URL")
	cfg.AppName = getEnvOrDefault("APP_NAME", "Femabras Daily Challenge")
	cfg.SMTPHost = getEnvOrDefault("SMTP_HOST", "smtp.gmail.com")
	cfg.SMTPPort = getEnvOrDefault("SMTP_PORT", "587")
	cfg.SMTPUsername = os.Getenv("SMTP_USERNAME")
	cfg.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	cfg.FromEmail = os.Getenv("FROM_EMAIL")
	cfg.RedisURL = os.Getenv("REDIS_URL")

	return cfg
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
