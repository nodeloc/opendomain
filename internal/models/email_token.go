package models

import "time"

const (
	EmailTokenPurposeVerifyEmail   = "verify_email"
	EmailTokenPurposeResetPassword = "reset_password"
)

// EmailToken 邮件一次性令牌（邮箱验证 / 重置密码）
type EmailToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"not null;index:idx_email_tokens_user_purpose" json:"user_id"`
	Purpose   string     `gorm:"size:20;not null;index:idx_email_tokens_user_purpose" json:"purpose"`
	TokenHash string     `gorm:"size:64;not null;uniqueIndex" json:"-"`
	Email     string     `gorm:"size:255;not null" json:"email"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at"`
}

func (EmailToken) TableName() string {
	return "email_tokens"
}
