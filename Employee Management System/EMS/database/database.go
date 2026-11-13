package database

import (
	"context"
	"fmt"
	"time"

	"EMS/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Connect establishes a PostgreSQL database connection.
func Connect(cfg config.Config, databaseName string) (*gorm.DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate configuration: %w", err)
	}

	if databaseName == "" {
		return nil, fmt.Errorf("database name cannot be empty")
	}

	// Connect to PostgreSQL using GORM.
	db, err := gorm.Open(
		postgres.Open(cfg.DatabaseURL(databaseName)),
		&gorm.Config{
			Logger: gormlogger.Default.LogMode(
				getGORMLogLevel(cfg.LogLevel),
			),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("connect to database %s: %w", databaseName, err)
	}

	// Get the underlying SQL database.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection: %w", err)
	}

	// Configure connection pooling.
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	// Verify that PostgreSQL is reachable.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

// getGORMLogLevel maps application log levels to GORM log levels.
func getGORMLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "DEBUG":
		return gormlogger.Info
	case "INFO":
		return gormlogger.Warn
	case "WARN":
		return gormlogger.Error
	case "ERROR":
		return gormlogger.Error
	default:
		return gormlogger.Silent
	}
}