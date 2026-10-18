package models

import "time"

const (
	AbuseReportStatusOpen      = "open"
	AbuseReportStatusReviewing = "reviewing"
	AbuseReportStatusResolved  = "resolved"
	AbuseReportStatusDismissed = "dismissed"
)

var AbuseReportStatuses = []string{
	AbuseReportStatusOpen,
	AbuseReportStatusReviewing,
	AbuseReportStatusResolved,
	AbuseReportStatusDismissed,
}

var AbuseReportCategories = []string{"phishing", "malware", "scam", "spam", "impersonation", "other"}

// AbuseReport 滥用举报
type AbuseReport struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Subdomain     string     `gorm:"size:255;not null;index" json:"subdomain"`
	Category      string     `gorm:"size:20;not null" json:"category"`
	URL           string     `gorm:"column:url;size:2048;not null;default:''" json:"url"`
	Details       string     `gorm:"type:text;not null" json:"details"`
	ReporterEmail string     `gorm:"size:255;not null;default:''" json:"reporter_email"`
	ReporterIP    string     `gorm:"column:reporter_ip;size:64;not null;default:''" json:"reporter_ip"`
	UserAgent     string     `gorm:"size:512;not null;default:''" json:"-"`
	Status        string     `gorm:"size:20;not null;default:open" json:"status"`
	AdminNote     string     `gorm:"type:text;not null;default:''" json:"admin_note"`
	DomainID      *uint      `gorm:"index" json:"domain_id"`
	Locale        string     `gorm:"size:10;not null;default:''" json:"-"`
	ResolvedAt    *time.Time `json:"resolved_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (AbuseReport) TableName() string {
	return "abuse_reports"
}

// AbuseReportCreateRequest 公开举报提交请求
type AbuseReportCreateRequest struct {
	Subdomain     string `json:"subdomain"`
	Category      string `json:"category"`
	URL           string `json:"url"`
	Details       string `json:"details"`
	ReporterEmail string `json:"reporter_email"`
}

// AbuseReportUpdateRequest 管理员更新举报请求
type AbuseReportUpdateRequest struct {
	Status    *string `json:"status"`
	AdminNote *string `json:"admin_note"`
}

// AbuseReportResponse 举报响应（管理员）
type AbuseReportResponse struct {
	ID            uint       `json:"id"`
	Subdomain     string     `json:"subdomain"`
	Category      string     `json:"category"`
	URL           string     `json:"url"`
	Details       string     `json:"details"`
	ReporterEmail string     `json:"reporter_email"`
	ReporterIP    string     `json:"reporter_ip"`
	Status        string     `json:"status"`
	AdminNote     string     `json:"admin_note"`
	DomainID      *uint      `json:"domain_id"`
	DomainStatus  *string    `json:"domain_status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ResolvedAt    *time.Time `json:"resolved_at"`
}

func (r *AbuseReport) ToResponse(domainStatus *string) *AbuseReportResponse {
	return &AbuseReportResponse{
		ID:            r.ID,
		Subdomain:     r.Subdomain,
		Category:      r.Category,
		URL:           r.URL,
		Details:       r.Details,
		ReporterEmail: r.ReporterEmail,
		ReporterIP:    r.ReporterIP,
		Status:        r.Status,
		AdminNote:     r.AdminNote,
		DomainID:      r.DomainID,
		DomainStatus:  domainStatus,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		ResolvedAt:    r.ResolvedAt,
	}
}
