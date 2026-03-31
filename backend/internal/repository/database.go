package repository

import (
	"fmt"
	"log"

	"multipass-backend/internal/config"
	"multipass-backend/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 自动迁移模型
	err = autoMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

// autoMigrate 自动创建表结构
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

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
