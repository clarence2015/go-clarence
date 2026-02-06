package config

import (
	"errors"
	"os"
)

// Config holds application runtime configuration.
type Config struct {
	Env          string
	HTTPPort     string
	DatabaseDSN  string
}

// Load reads configuration from environment variables with sensible defaults.
//
// Required:
//   - DATABASE_DSN: PostgreSQL DSN, e.g. postgres://user:pass@localhost:5432/goblog?sslmode=disable
//
// Optional:
//   - APP_ENV: "production" | "development" | "test" (default: "development")
//   - HTTP_PORT: HTTP listen port (default: "8080")
func Load() (*Config, error) {
	env := getEnv("APP_ENV", "development")
	httpPort := getEnv("HTTP_PORT", "8080")
	dsn := os.Getenv("DATABASE_DSN")

	if dsn == "" {
		return nil, errors.New("DATABASE_DSN is required")
	}

	return &Config{
		Env:         env,
		HTTPPort:    httpPort,
		DatabaseDSN: dsn,
	}, nil
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

