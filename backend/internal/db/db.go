package db

import (
	"context"
	"fmt"
	"time"

	appconfig "pilipili-go/backend/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func New(cfg appconfig.DatabaseConfig) (*gorm.DB, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("database dsn is required")
	}

	gormDB, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open gorm db: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinute) * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ConnectTimeoutSecond)*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	return gormDB, nil
}
