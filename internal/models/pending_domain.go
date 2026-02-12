package models

import (
	"time"

	"gorm.io/gorm"
)

// PendingDomain 待激活域名模型（从FOSSBilling预同步）
type PendingDomain struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	RootDomainID          uint           `gorm:"not null;index" json:"root_domain_id"`
	Subdomain             string         `gorm:"size:63;not null" json:"subdomain"`
	FullDomain            string         `gorm:"size:255;not null;uniqueIndex" json:"full_domain"`
	FOSSBillingOrderID    int            `gorm:"column:fossbilling_order_id;not null" json:"fossbilling_order_id"`
	Status                string         `gorm:"size:20;default:pending" json:"status"` // pending/healthy/unhealthy
	RegisteredAt          time.Time      `gorm:"not null" json:"registered_at"`
	ExpiresAt             time.Time      `gorm:"not null" json:"expires_at"`
	FirstFailedAt         *time.Time     `gorm:"index" json:"first_failed_at,omitempty"` // 健康检查失败时间
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	RootDomain *RootDomain `gorm:"foreignKey:RootDomainID" json:"root_domain,omitempty"`
}

// TableName 指定表名
func (PendingDomain) TableName() string {
	return "pending_domains"
}
