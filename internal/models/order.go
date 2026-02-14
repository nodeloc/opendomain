package models

import (
	"time"
)

// Order 订单模型
type Order struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	OrderNumber string `gorm:"size:32;unique;not null" json:"order_number"`
	UserID      uint   `gorm:"not null;index" json:"user_id"`

	// Domain information
	Subdomain    string `gorm:"size:63;not null" json:"subdomain"`
	RootDomainID uint   `gorm:"not null;index" json:"root_domain_id"`
	FullDomain   string `gorm:"size:255;not null" json:"full_domain"`

	// Pricing details
	Years          int     `gorm:"not null" json:"years"`
	IsLifetime     bool    `gorm:"default:false" json:"is_lifetime"`
	BasePrice      float64 `gorm:"type:decimal(10,2);not null" json:"base_price"`
	DiscountAmount float64 `gorm:"type:decimal(10,2);default:0" json:"discount_amount"`
	FinalPrice     float64 `gorm:"type:decimal(10,2);not null" json:"final_price"`

	// Coupon information
	CouponID   *uint   `json:"coupon_id,omitempty"`
	CouponCode *string `gorm:"size:50" json:"coupon_code,omitempty"`

	// Order status
	Status   string `gorm:"size:20;default:pending" json:"status"` // pending/paid/cancelled/refunded/expired
	DomainID *uint  `json:"domain_id,omitempty"`

	// Timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	PaidAt    *time.Time `json:"paid_at,omitempty"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`

	Notes *string `gorm:"type:text" json:"notes,omitempty"`

	// Relations
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	RootDomain *RootDomain `gorm:"foreignKey:RootDomainID" json:"root_domain,omitempty"`
	Domain     *Domain     `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
	Coupon     *Coupon     `gorm:"foreignKey:CouponID" json:"coupon,omitempty"`
	Payment    *Payment    `gorm:"foreignKey:OrderID" json:"payment,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// OrderCreateRequest 创建订单请求
type OrderCreateRequest struct {
	Subdomain    string  `json:"subdomain" binding:"required,min=3,max=63"`
	RootDomainID uint    `json:"root_domain_id" binding:"required"`
	Years        int     `json:"years" binding:"omitempty,min=0,max=10"`
	IsLifetime   bool    `json:"is_lifetime"`
	CouponCode   *string `json:"coupon_code"`
}

// OrderCalculateRequest 计算价格请求
type OrderCalculateRequest struct {
	RootDomainID uint    `json:"root_domain_id" binding:"required"`
	Years        int     `json:"years" binding:"omitempty,min=0,max=10"`
	IsLifetime   bool    `json:"is_lifetime"`
	CouponCode   *string `json:"coupon_code"`
}

// OrderResponse 订单响应
type OrderResponse struct {
	ID             uint        `json:"id"`
	OrderNumber    string      `json:"order_number"`
	FullDomain     string      `json:"full_domain"`
	Years          int         `json:"years"`
	IsLifetime     bool        `json:"is_lifetime"`
	BasePrice      float64     `json:"base_price"`
	DiscountAmount float64     `json:"discount_amount"`
	FinalPrice     float64     `json:"final_price"`
	Status         string      `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	ExpiresAt      time.Time   `json:"expires_at"`
	PaidAt         *time.Time  `json:"paid_at,omitempty"`
	RootDomain     *RootDomain `json:"root_domain,omitempty"`
	Payment        *Payment    `json:"payment,omitempty"`
}

// PriceCalculationResponse 价格计算响应
type PriceCalculationResponse struct {
	BasePrice      float64 `json:"base_price"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalPrice     float64 `json:"final_price"`
	CouponApplied  bool    `json:"coupon_applied"`
	CouponCode     *string `json:"coupon_code,omitempty"`
	CouponType     *string `json:"coupon_type,omitempty"`
	CouponError    *string `json:"coupon_error,omitempty"` // 优惠券验证失败的详细原因
}

// ToResponse 转换为响应格式
func (o *Order) ToResponse() *OrderResponse {
	return &OrderResponse{
		ID:             o.ID,
		OrderNumber:    o.OrderNumber,
		FullDomain:     o.FullDomain,
		Years:          o.Years,
		IsLifetime:     o.IsLifetime,
		BasePrice:      o.BasePrice,
		DiscountAmount: o.DiscountAmount,
		FinalPrice:     o.FinalPrice,
		Status:         o.Status,
		CreatedAt:      o.CreatedAt,
		ExpiresAt:      o.ExpiresAt,
		PaidAt:         o.PaidAt,
		RootDomain:     o.RootDomain,
		Payment:        o.Payment,
	}
}
