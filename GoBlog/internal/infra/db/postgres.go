package db

import (
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres creates a new gorm.DB connection to PostgreSQL using the provided DSN.
func NewPostgres(dsn string, logger *slog.Logger) (*gorm.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	logger.Info("connected to PostgreSQL database")

	return gormDB, nil
}

