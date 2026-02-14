package models

import (
	"time"
)

type DomainScan struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	DomainID       uint       `gorm:"not null" json:"domain_id"`
	ScanType       string     `gorm:"not null" json:"scan_type"` // http, dns, ssl
	Status         string     `gorm:"not null" json:"status"`    // success, failed, timeout
	ResponseTime   *int       `json:"response_time"`             // milliseconds
	HTTPStatusCode *int       `json:"http_status_code"`
	SSLValid       *bool      `json:"ssl_valid"`
	SSLExpiryDate  *time.Time `json:"ssl_expiry_date"`
	ErrorMessage   string     `json:"error_message"`
	ScanDetails    *string    `gorm:"type:jsonb" json:"scan_details"`
	ScannedAt      time.Time  `gorm:"not null" json:"scanned_at"`
	CreatedAt      time.Time  `json:"created_at"`

	Domain *Domain `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
}

type DomainScanSummary struct {
	DomainID           uint       `gorm:"primaryKey" json:"domain_id"`
	LastScannedAt      *time.Time `json:"last_scanned_at"`
	HTTPStatus         string     `json:"http_status"`                                       // online, offline, error
	DNSStatus          string     `json:"dns_status"`                                        // resolved, failed
	SSLStatus          string     `json:"ssl_status"`                                        // valid, invalid, none
	SafeBrowsingStatus string     `json:"safe_browsing_status"`                              // safe, unsafe, unknown
	VirusTotalStatus   string     `gorm:"column:virustotal_status" json:"virustotal_status"` // clean, malicious, suspicious, unknown
	OverallHealth      string     `json:"overall_health"`                                    // healthy, degraded, down
	TotalScans         int        `json:"total_scans"`
	SuccessfulScans    int        `json:"successful_scans"`
	UpdatedAt          time.Time  `json:"updated_at"`

	Domain *Domain `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
}

type DomainScanResponse struct {
	ID             uint       `json:"id"`
	DomainID       uint       `json:"domain_id"`
	DomainName     string     `json:"domain_name"`
	ScanType       string     `json:"scan_type"`
	Status         string     `json:"status"`
	ResponseTime   *int       `json:"response_time"`
	HTTPStatusCode *int       `json:"http_status_code"`
	SSLValid       *bool      `json:"ssl_valid"`
	SSLExpiryDate  *time.Time `json:"ssl_expiry_date"`
	ErrorMessage   string     `json:"error_message"`
	ScanDetails    *string    `json:"scan_details"`
	ScannedAt      time.Time  `json:"scanned_at"`
}

type DomainHealthResponse struct {
	DomainID           uint       `json:"domain_id"`
	DomainName         string     `json:"domain_name"`
	LastScannedAt      *time.Time `json:"last_scanned_at"`
	HTTPStatus         string     `json:"http_status"`
	DNSStatus          string     `json:"dns_status"`
	SSLStatus          string     `json:"ssl_status"`
	SafeBrowsingStatus string     `json:"safe_browsing_status"`
	VirusTotalStatus   string     `json:"virustotal_status"`
	OverallHealth      string     `json:"overall_health"`
	TotalScans         int        `json:"total_scans"`
	SuccessfulScans    int        `json:"successful_scans"`
	Uptime             float64    `json:"uptime_percentage"`
	FirstFailedAt      *time.Time `json:"first_failed_at,omitempty"`
	IsSuspended        bool       `json:"is_suspended"`
}

// AggregatedScanRecord 聚合的扫描记录（用于历史记录显示）
type AggregatedScanRecord struct {
	ScannedAt          time.Time `json:"scanned_at"`
	Status             string    `json:"status"` // completed, threat_detected, failed
	HTTPStatus         string    `json:"http_status"`
	DNSStatus          string    `json:"dns_status"`
	SSLStatus          string    `json:"ssl_status"`
	SafeBrowsingStatus string    `json:"safe_browsing_status"`
	VirusTotalStatus   string    `json:"virustotal_status"`
}

func (s *DomainScan) ToResponse() *DomainScanResponse {
	resp := &DomainScanResponse{
		ID:             s.ID,
		DomainID:       s.DomainID,
		ScanType:       s.ScanType,
		Status:         s.Status,
		ResponseTime:   s.ResponseTime,
		HTTPStatusCode: s.HTTPStatusCode,
		SSLValid:       s.SSLValid,
		SSLExpiryDate:  s.SSLExpiryDate,
		ErrorMessage:   s.ErrorMessage,
		ScanDetails:    s.ScanDetails,
		ScannedAt:      s.ScannedAt,
	}

	if s.Domain != nil {
		resp.DomainName = s.Domain.FullDomain
	}

	return resp
}

func (s *DomainScanSummary) ToResponse() *DomainHealthResponse {
	resp := &DomainHealthResponse{
		DomainID:           s.DomainID,
		LastScannedAt:      s.LastScannedAt,
		HTTPStatus:         s.HTTPStatus,
		DNSStatus:          s.DNSStatus,
		SSLStatus:          s.SSLStatus,
		SafeBrowsingStatus: s.SafeBrowsingStatus,
		VirusTotalStatus:   s.VirusTotalStatus,
		OverallHealth:      s.OverallHealth,
		TotalScans:         s.TotalScans,
		SuccessfulScans:    s.SuccessfulScans,
	}

	if s.Domain != nil {
		resp.DomainName = s.Domain.FullDomain
		resp.FirstFailedAt = s.Domain.FirstFailedAt
		resp.IsSuspended = (s.Domain.Status == "suspended")
	}

	if s.TotalScans > 0 {
		resp.Uptime = float64(s.SuccessfulScans) / float64(s.TotalScans) * 100
	}

	return resp
}
