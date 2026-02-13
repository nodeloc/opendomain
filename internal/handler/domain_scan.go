package handler

import (
	"fmt"
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
	// 获取查询参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")
	search := c.Query("search")

	// 解析分页参数
	var pageInt, pageSizeInt int
	if _, err := fmt.Sscanf(page, "%d", &pageInt); err != nil || pageInt < 1 {
		pageInt = 1
	}
	if _, err := fmt.Sscanf(pageSize, "%d", &pageSizeInt); err != nil || pageSizeInt < 1 {
		pageSizeInt = 20
	}
	if pageSizeInt > 100 {
		pageSizeInt = 100
	}

	// 构建查询
	query := h.db.Model(&models.DomainScanSummary{}).Preload("Domain")

	// 如果有搜索关键词，通过域名名称搜索
	if search != "" {
		query = query.Joins("JOIN domains ON domains.id = domain_scan_summaries.domain_id").
			Where("domains.name LIKE ?", "%"+search+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count domain health"})
		return
	}

	// 分页查询
	var summaries []models.DomainScanSummary
	offset := (pageInt - 1) * pageSizeInt
	if err := query.Order("last_scanned_at DESC").
		Offset(offset).
		Limit(pageSizeInt).
		Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domain health"})
		return
	}

	responses := make([]*models.DomainHealthResponse, len(summaries))
	for i, summary := range summaries {
		responses[i] = summary.ToResponse()
	}

	// 计算总页数
	totalPages := int((total + int64(pageSizeInt) - 1) / int64(pageSizeInt))

	c.JSON(http.StatusOK, gin.H{
		"health_reports": responses,
		"pagination": gin.H{
			"page":        pageInt,
			"page_size":   pageSizeInt,
			"total":       total,
			"total_pages": totalPages,
		},
	})
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
