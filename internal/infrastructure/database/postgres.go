package database

import (
	"fmt"
	"log"
	"ss-catalog-service/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewPostgresDB creates and validates a new GORM PostgreSQL connection.
// It does NOT perform migrations — migration is the responsibility of the caller (main.go).
func NewPostgresDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(gormpostgres.Open(cfg.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

// RunMigrations executes SQL migrations using golang-migrate
func RunMigrations(db *gorm.DB, migrationsPath string) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for migration: %w", err)
	}

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver for migration: %w", err)
	}

	// Source path for file://
	sourceURL := fmt.Sprintf("file://%s", migrationsPath)
	
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	log.Printf("Running migrations from %s", migrationsPath)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
