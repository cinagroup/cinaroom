package pool

import (
	"time"

	"gorm.io/gorm"
)

// DBPoolConfig holds database connection pool settings.
type DBPoolConfig struct {
	MaxOpen     int           // Maximum open connections (default 50)
	MaxIdle     int           // Maximum idle connections (default 10)
	MaxLifetime time.Duration // Maximum connection lifetime (default 1h)
	MaxIdleTime time.Duration // Maximum idle connection lifetime (default 10min)
}

// DefaultDBPoolConfig returns the recommended pool configuration.
func DefaultDBPoolConfig() DBPoolConfig {
	return DBPoolConfig{
		MaxOpen:     50,
		MaxIdle:     10,
		MaxLifetime: 1 * time.Hour,
		MaxIdleTime: 10 * time.Minute,
	}
}

// ApplyToDB applies the pool configuration to a gorm.DB instance.
func ApplyToDB(db *gorm.DB, cfg DBPoolConfig) {
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	if cfg.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpen)
	}
	if cfg.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	}
	if cfg.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)
	}
	if cfg.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(cfg.MaxIdleTime)
	}
}
