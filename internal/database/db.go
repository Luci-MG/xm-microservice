package database

import (
	"database/sql"
	"fmt"

	"xm-microservice/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Connect establishes a connection to the PostgreSQL database and runs migrations
func Connect(databaseURL string, log *logger.Logger) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Error(err, "Failed to open database connection")
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error(err, "Database ping failed")
		return nil, err
	}

	log.Info("Successfully connected to the database")

	// Run database migrations
	if err := runMigrations(databaseURL, log); err != nil {
		log.Error(err, "Migration error")
		return nil, err
	}

	return db, nil
}

// runMigrations applies any pending database migrations
func runMigrations(databaseURL string, log *logger.Logger) error {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Error(err, "Failed to open database for migrations")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Error(err, "Failed to create migration driver")
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./internal/database/migrations", // Path to migration files
		"postgres",                              // Database name
		driver,                                  // PostgreSQL migration driver
	)
	if err != nil {
		log.Error(err, "Failed to initialize migration")
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Error(err, "Migration failed")
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Info("Database migrated successfully")
	return nil
}
