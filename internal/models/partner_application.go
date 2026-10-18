package models

import "time"

const (
	PartnerApplicationStatusNew       = "new"
	PartnerApplicationStatusContacted = "contacted"
	PartnerApplicationStatusAccepted  = "accepted"
	PartnerApplicationStatusDeclined  = "declined"
)

var PartnerApplicationStatuses = []string{
	PartnerApplicationStatusNew,
	PartnerApplicationStatusContacted,
	PartnerApplicationStatusAccepted,
	PartnerApplicationStatusDeclined,
}

var PartnershipTypes = []string{"sponsorship", "resources", "business", "community", "other"}

// PartnerApplication 合作申请
type PartnerApplication struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"size:100;not null" json:"name"`
	Email           string    `gorm:"size:254;not null" json:"email"`
	Organization    string    `gorm:"size:200;not null;default:''" json:"organization"`
	Website         string    `gorm:"size:500;not null;default:''" json:"website"`
	PartnershipType string    `gorm:"size:20;not null" json:"partnership_type"`
	Message         string    `gorm:"type:text;not null" json:"message"`
	Contact         string    `gorm:"size:200;not null;default:''" json:"contact"`
	Status          string    `gorm:"size:20;not null;default:new" json:"status"`
	AdminNote       string    `gorm:"type:text;not null;default:''" json:"admin_note"`
	IP              string    `gorm:"column:ip;size:64;not null;default:''" json:"ip"`
	UserAgent       string    `gorm:"size:512;not null;default:''" json:"-"`
	Locale          string    `gorm:"size:10;not null;default:''" json:"locale"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (PartnerApplication) TableName() string {
	return "partner_applications"
}

// PartnerApplicationCreateRequest 公开合作申请提交请求
type PartnerApplicationCreateRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Organization    string `json:"organization"`
	Website         string `json:"website"`
	PartnershipType string `json:"partnership_type"`
	Message         string `json:"message"`
	Contact         string `json:"contact"`
	// Fax is a honeypot field that is hidden from people; bots tend to fill it.
	Fax string `json:"fax"`
}

// PartnerApplicationUpdateRequest 管理员更新合作申请请求
type PartnerApplicationUpdateRequest struct {
	Status    *string `json:"status"`
	AdminNote *string `json:"admin_note"`
}
