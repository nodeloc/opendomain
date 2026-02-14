package models

import "time"

// APIQuota 存储 API 配额使用情况
type APIQuota struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	APIName    string    `gorm:"uniqueIndex;not null" json:"api_name"` // "google_safe_browsing" or "virustotal"
	Date       string    `gorm:"not null" json:"date"`                 // YYYY-MM-DD
	UsedCount  int       `gorm:"not null;default:0" json:"used_count"`
	DailyLimit int       `gorm:"not null" json:"daily_limit"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (APIQuota) TableName() string {
	return "api_quotas"
}
