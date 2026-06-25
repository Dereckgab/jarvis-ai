package database

import (
	"fmt"
	"log"
	"time"

	"jarvis/config"
	"jarvis/internal/domain/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMySQLConnection creates and returns a new GORM MySQL database connection.
func NewMySQLConnection(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Log all parameters
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(cfg.GetDSN()), &gorm.Config{
		Logger: newLogger,
		// Add other GORM configurations here, e.g., NamingStrategy
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

// RunMigrations performs database auto-migrations for all entities.
func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.SystemInfo{},
		&entity.Game{},
		&entity.GameRequirement{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}
	return nil
}
