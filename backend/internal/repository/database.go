package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/cinagroup/cinaseek/backend/internal/config"
	"github.com/cinagroup/cinaseek/backend/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initialises the database connection with connection pooling and runs migrations.
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s search_path=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.Schema)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// Run auto-migration (non-fatal: tables may already exist from SQL migration)
	if err := autoMigrate(); err != nil {
		slog.Warn("auto-migration skipped (tables likely exist from SQL migration)", "error", err)
	}

	slog.Info("database connection established",
		"host", cfg.Host,
		"port", cfg.Port,
		"database", cfg.DBName,
		"schema", cfg.Schema,
	)
	return nil
}

// autoMigrate creates or updates all table schemas.
func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.VM{},
		&model.VMSnapshot{},
		&model.VMLog{},
		&model.Mount{},
		&model.OpenClawConfig{},
		&model.RemoteAccess{},
		&model.IPWhitelist{},
		&model.RemoteLog{},
		&model.LoginLog{},
		&model.Session{},
		&model.SystemSetting{},
		&model.VMMetric{},
	)
}

// GetDB returns the global database instance.
func GetDB() *gorm.DB {
	return DB
}

// HealthCheck pings the database to verify connectivity.
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialised")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	return nil
}

// CloseDB gracefully closes the database connection pool.
func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB for closing: %w", err)
	}

	// Give pending queries up to 5 seconds to finish.
	maxWait := time.After(5 * time.Second)
	done := make(chan error, 1)
	go func() {
		done <- sqlDB.Close()
	}()

	select {
	case err := <-done:
		slog.Info("database connection closed")
		return err
	case <-maxWait:
		slog.Warn("database close timed out, forcing shutdown")
		return sqlDB.Close()
	}
}
