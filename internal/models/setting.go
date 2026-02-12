package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemSetting 系统设置模型
type SystemSetting struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	SettingKey   string `gorm:"size:100;not null;uniqueIndex" json:"setting_key"`
	SettingValue string `gorm:"type:text;not null" json:"setting_value"`
	Description  string `gorm:"size:255" json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}

// GetSettingValue retrieves a setting value by key with a fallback default
func GetSettingValue(db *gorm.DB, key string, defaultValue string) string {
	var setting SystemSetting
	if err := db.Where("setting_key = ?", key).First(&setting).Error; err != nil {
		return defaultValue
	}
	return setting.SettingValue
}
