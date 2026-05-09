package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Name        string
	Port        string
	Environment string
	DefaultLang string
}

type DatabaseConfig struct {
	DSN string
}

type JWTConfig struct {
	Issuer         string
	Audience       string
	PublicKeyPath  string
}

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// Load reads configuration from .env and environment variables.
func Load() *Config {
	// Load .env if exists (ignore error as env may be set by OS)
	_ = godotenv.Load()

	cfg := &Config{}

	// App Config
	cfg.App.Name = getEnv("APP_NAME", "SS-CatalogService")
	cfg.App.Port = getEnv("APP_PORT", "8081")
	cfg.App.Environment = getEnv("APP_ENV", "development")
	cfg.App.DefaultLang = getEnv("DEFAULT_LANG", "id-ID")

	// Database Config
	cfg.Database.DSN = os.Getenv("DB_DSN")
	if cfg.Database.DSN == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	// JWT Config
	cfg.JWT.Issuer = getEnv("JWT_ISSUER", "https://auth.internal.smallmarket.com")
	cfg.JWT.Audience = getEnv("JWT_AUDIENCE", "smallmarket-api")
	cfg.JWT.PublicKeyPath = getEnv("JWT_PUBLIC_KEY_PATH", "../../secrets/jwt_public_key.pem")

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
