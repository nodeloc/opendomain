package models

import (
	"strings"
	"time"
)

// FlexibleTime 支持多种时间格式的自定义时间类型
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON 自定义 JSON 解析，支持多种时间格式
func (ft *FlexibleTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// 尝试多种时间格式
	formats := []string{
		time.RFC3339,                // 2006-01-02T15:04:05Z07:00
		"2006-01-02T15:04:05",      // 2006-01-02T15:04:05
		"2006-01-02T15:04",         // 2006-01-02T15:04 (前端 datetime-local)
		"2006-01-02 15:04:05",      // 2006-01-02 15:04:05
		"2006-01-02 15:04",         // 2006-01-02 15:04
		"2006-01-02",               // 2006-01-02
	}

	var err error
	for _, format := range formats {
		ft.Time, err = time.Parse(format, s)
		if err == nil {
			return nil
		}
	}

	return err
}

type Coupon struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Code          string     `gorm:"unique;not null" json:"code"`
	Description   string     `json:"description"`
	DiscountType  string     `gorm:"not null" json:"discount_type"` // percentage, fixed, quota_increase
	DiscountValue *float64   `json:"discount_value"`
	QuotaIncrease int        `json:"quota_increase"`
	MaxUses       int        `json:"max_uses"` // 0 = unlimited
	UsedCount     int        `json:"used_count"`
	ValidFrom     time.Time  `gorm:"not null" json:"valid_from"`
	ValidUntil    *time.Time `json:"valid_until"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	IsReusable    bool       `gorm:"default:false" json:"is_reusable"` // 是否可以被同一用户多次使用
	CreatedBy     *uint      `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CouponUsage struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CouponID       uint      `gorm:"not null" json:"coupon_id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	DomainID       *uint     `json:"domain_id"`
	UsedAt         time.Time `gorm:"not null" json:"used_at"`
	BenefitApplied string    `json:"benefit_applied"`
	Coupon         *Coupon   `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
}

// TableName 指定表名为单数形式，匹配迁移文件中的表名
func (CouponUsage) TableName() string {
	return "coupon_usage"
}

// Request/Response models
type CouponCreateRequest struct {
	Code          string        `json:"code" binding:"required,min=3,max=50"`
	Description   string        `json:"description"`
	DiscountType  string        `json:"discount_type" binding:"required,oneof=percentage fixed quota_increase"`
	DiscountValue *float64      `json:"discount_value"`
	QuotaIncrease int           `json:"quota_increase"`
	MaxUses       int           `json:"max_uses"`
	ValidFrom     *FlexibleTime `json:"valid_from"`
	ValidUntil    *FlexibleTime `json:"valid_until"`
	IsReusable    bool          `json:"is_reusable"`
}

type CouponUpdateRequest struct {
	Description   *string       `json:"description"`
	DiscountValue *float64      `json:"discount_value"`
	QuotaIncrease *int          `json:"quota_increase"`
	MaxUses       *int          `json:"max_uses"`
	ValidFrom     *FlexibleTime `json:"valid_from"`
	ValidUntil    *FlexibleTime `json:"valid_until"`
	IsActive      *bool         `json:"is_active"`
	IsReusable    *bool         `json:"is_reusable"`
}

type CouponApplyRequest struct {
	Code string `json:"code" binding:"required"`
}

type CouponResponse struct {
	ID            uint       `json:"id"`
	Code          string     `json:"code"`
	Description   string     `json:"description"`
	DiscountType  string     `json:"discount_type"`
	DiscountValue *float64   `json:"discount_value"`
	QuotaIncrease int        `json:"quota_increase"`
	MaxUses       int        `json:"max_uses"`
	UsedCount     int        `json:"used_count"`
	ValidFrom     time.Time  `json:"valid_from"`
	ValidUntil    *time.Time `json:"valid_until"`
	IsActive      bool       `json:"is_active"`
	IsReusable    bool       `json:"is_reusable"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (c *Coupon) ToResponse() *CouponResponse {
	return &CouponResponse{
		ID:            c.ID,
		Code:          c.Code,
		Description:   c.Description,
		DiscountType:  c.DiscountType,
		DiscountValue: c.DiscountValue,
		QuotaIncrease: c.QuotaIncrease,
		MaxUses:       c.MaxUses,
		UsedCount:     c.UsedCount,
		ValidFrom:     c.ValidFrom,
		ValidUntil:    c.ValidUntil,
		IsActive:      c.IsActive,
		IsReusable:    c.IsReusable,
		CreatedAt:     c.CreatedAt,
	}
}
