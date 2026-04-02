package repository

import (
	"fmt"

	"github.com/cinagroup/cinaseek/backend/internal/model"
)

// SystemRepo provides operations for system-level tables.
type SystemRepo struct {
}

// NewSystemRepo creates a new SystemRepo.
func NewSystemRepo() *SystemRepo {
	return &SystemRepo{}
}

// GetSetting retrieves a system setting by key.
func (r *SystemRepo) GetSetting(key string) (*model.SystemSetting, error) {
	db := GetDB()
	var setting model.SystemSetting
	if err := db.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return nil, fmt.Errorf("get setting %q: %w", key, err)
	}
	return &setting, nil
}

// SetSetting creates or updates a system setting.
func (r *SystemRepo) SetSetting(key, value string) error {
	db := GetDB()
	var setting model.SystemSetting
	result := db.Where("`key` = ?", key).First(&setting)

	setting.Key = key
	setting.Value = value

	if result.Error != nil {
		if err := db.Create(&setting).Error; err != nil {
			return fmt.Errorf("create setting: %w", err)
		}
	} else {
		if err := db.Save(&setting).Error; err != nil {
			return fmt.Errorf("save setting: %w", err)
		}
	}
	return nil
}

// GetAllSettings returns all system settings as a map.
func (r *SystemRepo) GetAllSettings() (map[string]string, error) {
	db := GetDB()
	var settings []model.SystemSetting
	if err := db.Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("get all settings: %w", err)
	}

	m := make(map[string]string, len(settings))
	for _, s := range settings {
		m[s.Key] = s.Value
	}
	return m, nil
}
