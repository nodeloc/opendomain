package models

import (
	"time"

	"gorm.io/gorm"
)

// RootDomain 根域名模型
type RootDomain struct {
	ID                    uint      `gorm:"primarykey" json:"id"`
	Domain                string    `gorm:"size:100;not null;uniqueIndex" json:"domain"`
	Description           *string   `json:"description,omitempty"`
	Nameservers           string    `gorm:"type:text;not null" json:"nameservers"` // JSON array as text
	UseDefaultNameservers bool      `gorm:"default:true" json:"use_default_nameservers"`
	IsActive              bool      `gorm:"default:true" json:"is_active"`
	IsHot                 bool      `gorm:"default:false" json:"is_hot"`
	IsNew                 bool      `gorm:"default:false" json:"is_new"`
	Priority              int       `gorm:"default:0" json:"priority"`
	MinLength             int       `gorm:"default:3" json:"min_length"`
	MaxLength             int       `gorm:"default:63" json:"max_length"`
	RegistrationCount     int       `gorm:"default:0" json:"registration_count"`
	PricePerYear          *float64  `gorm:"type:decimal(10,2)" json:"price_per_year,omitempty"`
	LifetimePrice         *float64  `gorm:"type:decimal(10,2)" json:"lifetime_price,omitempty"`
	IsFree                bool      `gorm:"default:true" json:"is_free"`
	MinUserLevel          string    `gorm:"size:20;default:normal" json:"min_user_level"` // normal/basic/member/regular/leader
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// TableName 指定表名
func (RootDomain) TableName() string {
	return "root_domains"
}

// Domain 域名模型
type Domain struct {
	ID                    uint           `gorm:"primarykey" json:"id"`
	UserID                uint           `gorm:"not null;index" json:"user_id"`
	RootDomainID          uint           `gorm:"not null;index" json:"root_domain_id"`
	Subdomain             string         `gorm:"size:63;not null" json:"subdomain"`
	FullDomain            string         `gorm:"size:255;not null;uniqueIndex" json:"full_domain"`
	Status                string         `gorm:"size:20;default:active" json:"status"` // active/expired/suspended/deleted
	RegisteredAt          time.Time      `gorm:"not null" json:"registered_at"`
	ExpiresAt             time.Time      `gorm:"not null" json:"expires_at"`
	AutoRenew             bool           `gorm:"default:false" json:"auto_renew"`
	Nameservers           string         `gorm:"type:text" json:"nameservers"`
	UseDefaultNameservers bool           `gorm:"default:true" json:"use_default_nameservers"`
	ReminderSent30d       bool           `gorm:"column:reminder_sent_30d;default:false" json:"-"`
	ReminderSent7d        bool           `gorm:"column:reminder_sent_7d;default:false" json:"-"`
	DNSSynced             bool           `gorm:"default:false" json:"dns_synced"`
	DNSSyncError          *string        `gorm:"type:text" json:"dns_sync_error,omitempty"`
	FirstFailedAt         *time.Time     `gorm:"index" json:"first_failed_at,omitempty"`    // Tracks when domain first failed health check
	SuspendedAt           *time.Time     `json:"suspended_at,omitempty"`                    // When domain was suspended
	SuspendReason         *string        `gorm:"type:text" json:"suspend_reason,omitempty"` // Reason for suspension
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`

	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RootDomain *RootDomain `gorm:"foreignKey:RootDomainID" json:"root_domain,omitempty"`
}

// TableName 指定表名
func (Domain) TableName() string {
	return "domains"
}

// DomainSearchRequest 域名搜索请求
type DomainSearchRequest struct {
	Subdomain    string `form:"subdomain" binding:"required,min=3,max=63"`
	RootDomainID uint   `form:"root_domain_id" binding:"required"`
}

// DomainRegisterRequest 域名注册请求
type DomainRegisterRequest struct {
	Subdomain    string  `json:"subdomain" binding:"required,min=3,max=63"`
	RootDomainID uint    `json:"root_domain_id" binding:"required"`
	CouponCode   *string `json:"coupon_code,omitempty"`
}

// AdminCreateDomainRequest 管理员为指定用户创建域名请求
type AdminCreateDomainRequest struct {
	UserID       uint   `json:"user_id" binding:"required"`
	Subdomain    string `json:"subdomain" binding:"required,min=1,max=63"`
	RootDomainID uint   `json:"root_domain_id" binding:"required"`
	Years        int    `json:"years" binding:"omitempty,min=1,max=10"`
}

// DomainResponse 域名响应
type DomainResponse struct {
	ID                    uint          `json:"id"`
	UserID                uint          `json:"user_id"`
	RootDomainID          uint          `json:"root_domain_id"`
	Subdomain             string        `json:"subdomain"`
	FullDomain            string        `json:"full_domain"`
	Status                string        `json:"status"`
	RegisteredAt          time.Time     `json:"registered_at"`
	ExpiresAt             time.Time     `json:"expires_at"`
	AutoRenew             bool          `json:"auto_renew"`
	Nameservers           string        `json:"nameservers"`
	UseDefaultNameservers bool          `json:"use_default_nameservers"`
	DNSSynced             bool          `json:"dns_synced"`
	FirstFailedAt         *time.Time    `json:"first_failed_at,omitempty"`
	SuspendedAt           *time.Time    `json:"suspended_at,omitempty"`
	SuspendReason         *string       `json:"suspend_reason,omitempty"`
	RootDomain            *RootDomain   `json:"root_domain,omitempty"`
	User                  *UserResponse `json:"user,omitempty"`
}

var maxJSONTime = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)

// clampJSONTime keeps a time within the range encoding/json can marshal; a
// single out-of-range expiry would otherwise make the whole list response
// fail to encode and arrive empty.
func clampJSONTime(t time.Time) time.Time {
	if t.After(maxJSONTime) {
		return maxJSONTime
	}
	return t
}

// ToResponse 转换为响应格式
func (d *Domain) ToResponse() *DomainResponse {
	resp := &DomainResponse{
		ID:                    d.ID,
		UserID:                d.UserID,
		RootDomainID:          d.RootDomainID,
		Subdomain:             d.Subdomain,
		FullDomain:            d.FullDomain,
		Status:                d.Status,
		RegisteredAt:          clampJSONTime(d.RegisteredAt),
		ExpiresAt:             clampJSONTime(d.ExpiresAt),
		AutoRenew:             d.AutoRenew,
		Nameservers:           d.Nameservers,
		UseDefaultNameservers: d.UseDefaultNameservers,
		DNSSynced:             d.DNSSynced,
		FirstFailedAt:         d.FirstFailedAt,
		SuspendedAt:           d.SuspendedAt,
		SuspendReason:         d.SuspendReason,
		RootDomain:            d.RootDomain,
	}
	if d.User != nil {
		resp.User = d.User.ToResponse()
	}
	return resp
}
