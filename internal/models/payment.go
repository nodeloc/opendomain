package models

import (
	"time"
)

// Payment 支付记录模型
type Payment struct {
	ID                 uint    `gorm:"primarykey" json:"id"`
	OrderID            uint    `gorm:"not null;index" json:"order_id"`

	// NodeLoc payment details
	TransactionID      *string `gorm:"size:100;unique" json:"transaction_id,omitempty"`
	NodelocPaymentID   string  `gorm:"size:50;not null" json:"nodeloc_payment_id"`

	// Payment information
	Amount   float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency string  `gorm:"size:3;default:CNY" json:"currency"`

	// Payment status
	Status string `gorm:"size:20;default:pending" json:"status"` // pending/processing/completed/failed/refunded

	// Gateway response
	GatewayResponse *string `gorm:"type:text" json:"gateway_response,omitempty"`
	Signature       *string `gorm:"size:255" json:"signature,omitempty"`

	// Callback information
	CallbackReceivedAt *time.Time `json:"callback_received_at,omitempty"`
	CallbackIP         *string    `gorm:"size:45" json:"callback_ip,omitempty"`

	// Timestamps
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// Relations
	Order *Order `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

// TableName 指定表名
func (Payment) TableName() string {
	return "payments"
}

// PaymentConfig 支付配置模型
type PaymentConfig struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Provider    string `gorm:"size:20;default:nodeloc" json:"provider"`
	PaymentID   string `gorm:"size:50;not null" json:"payment_id"`
	SecretKey   string `gorm:"size:255;not null" json:"-"` // Never expose in JSON
	CallbackURL string `gorm:"size:255" json:"callback_url"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	IsTestMode  bool   `gorm:"default:false" json:"is_test_mode"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (PaymentConfig) TableName() string {
	return "payment_configs"
}

// NodelocCallbackRequest NodeLoc 回调请求
type NodelocCallbackRequest struct {
	TransactionID      string  `form:"transaction_id" binding:"required"`
	ExternalReference  string  `form:"external_reference"` // Our order_number
	Amount             float64 `form:"amount" binding:"required"`
	PlatformFee        float64 `form:"platform_fee"`
	MerchantPoints     float64 `form:"merchant_points"`
	Status             string  `form:"status" binding:"required"` // completed/failed/cancelled
	PaidAt             string  `form:"paid_at"`
	Signature          string  `form:"signature" binding:"required"`
}

// PaymentInitiateResponse 支付发起响应
type PaymentInitiateResponse struct {
	PaymentID   uint   `json:"payment_id"`
	RedirectURL string `json:"redirect_url"`
}
