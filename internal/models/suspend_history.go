package models

import "time"

// SuspendHistory 记录域名被挂起的历史
type SuspendHistory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	DomainID  uint      `json:"domain_id" gorm:"not null;index"`
	Domain    *Domain   `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Reason    string    `json:"reason" gorm:"type:text;not null"` // 挂起原因
	Details   string    `json:"details" gorm:"type:text"`         // 详细信息
	CreatedAt time.Time `json:"created_at"`

	// UI 相关字段
	DomainName string `json:"domain_name" gorm:"-"` // 域名全称（用于查询结果）
}

// SuspendHistoryResponse 用于API响应
type SuspendHistoryResponse struct {
	ID         uint      `json:"id"`
	DomainID   uint      `json:"domain_id"`
	DomainName string    `json:"domain_name"`
	Reason     string    `json:"reason"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}

func (h *SuspendHistory) ToResponse() *SuspendHistoryResponse {
	domainName := h.DomainName
	if domainName == "" && h.Domain != nil {
		domainName = h.Domain.FullDomain
	}

	return &SuspendHistoryResponse{
		ID:         h.ID,
		DomainID:   h.DomainID,
		DomainName: domainName,
		Reason:     h.Reason,
		Details:    h.Details,
		CreatedAt:  h.CreatedAt,
	}
}
