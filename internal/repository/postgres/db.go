package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/Shivanand-hulikatti/digit-go-services/land-services-go/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultPingTimeout = 5 * time.Second

// Open opens a GORM PostgreSQL connection for read-only land search.
func Open(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("underlying sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// Ping checks database connectivity.
func Ping(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is not initialized (check configs/app.yaml or LAND_DATABASE_* environment variables)")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("underlying sql db: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// CloseDB closes the underlying connection pool.
func CloseDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
