package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/models"
)

type DomainScanHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewDomainScanHandler(db *gorm.DB, cfg *config.Config) *DomainScanHandler {
	return &DomainScanHandler{db: db, cfg: cfg}
}

// ListDomainHealth 获取所有域名健康状态 (公开)
func (h *DomainScanHandler) ListDomainHealth(c *gin.Context) {
	var summaries []models.DomainScanSummary
	if err := h.db.Preload("Domain").Order("last_scanned_at DESC").Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domain health"})
		return
	}

	responses := make([]*models.DomainHealthResponse, len(summaries))
	for i, summary := range summaries {
		responses[i] = summary.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"health_reports": responses})
}

// GetDomainHealth 获取单个域名健康状态
func (h *DomainScanHandler) GetDomainHealth(c *gin.Context) {
	domainID := c.Param("domainId")

	var summary models.DomainScanSummary
	if err := h.db.Preload("Domain").Where("domain_id = ?", domainID).First(&summary).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain health not found"})
		return
	}

	c.JSON(http.StatusOK, summary.ToResponse())
}

// GetDomainScans 获取域名扫描历史
func (h *DomainScanHandler) GetDomainScans(c *gin.Context) {
	domainID := c.Param("domainId")

	var scans []models.DomainScan
	if err := h.db.Preload("Domain").
		Where("domain_id = ?", domainID).
		Order("scanned_at DESC").
		Limit(50).
		Find(&scans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scans"})
		return
	}

	responses := make([]*models.DomainScanResponse, len(scans))
	for i, scan := range scans {
		responses[i] = scan.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"scans": responses})
}

// GetHealthStatistics 获取健康统计信息
func (h *DomainScanHandler) GetHealthStatistics(c *gin.Context) {
	var stats struct {
		TotalDomains   int64 `json:"total_domains"`
		HealthyDomains int64 `json:"healthy_domains"`
		DegradedDomains int64 `json:"degraded_domains"`
		DownDomains    int64 `json:"down_domains"`
	}

	h.db.Model(&models.DomainScanSummary{}).Count(&stats.TotalDomains)
	h.db.Model(&models.DomainScanSummary{}).Where("overall_health = ?", "healthy").Count(&stats.HealthyDomains)
	h.db.Model(&models.DomainScanSummary{}).Where("overall_health = ?", "degraded").Count(&stats.DegradedDomains)
	h.db.Model(&models.DomainScanSummary{}).Where("overall_health = ?", "down").Count(&stats.DownDomains)

	c.JSON(http.StatusOK, stats)
}
