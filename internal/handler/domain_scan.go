package handler

import (
	"fmt"
	"net/http"
	"time"

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
	query := h.db.Model(&models.DomainScanSummary{})

	// 如果有搜索关键词，通过域名名称搜索
	if search != "" {
		query = query.Joins("JOIN domains ON domains.id = domain_scan_summaries.domain_id").
			Where("domains.full_domain LIKE ?", "%"+search+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count domain health"})
		return
	}

	// 重新构建查询用于分页（因为 Count 可能会修改查询状态）
	query = h.db.Model(&models.DomainScanSummary{}).Preload("Domain")
	if search != "" {
		query = query.Joins("JOIN domains ON domains.id = domain_scan_summaries.domain_id").
			Where("domains.full_domain LIKE ?", "%"+search+"%")
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
		TotalDomains    int64 `json:"total_domains"`
		HealthyDomains  int64 `json:"healthy_domains"`
		DegradedDomains int64 `json:"degraded_domains"`
		DownDomains     int64 `json:"down_domains"`
	}

	h.db.Model(&models.DomainScanSummary{}).Count(&stats.TotalDomains)
	h.db.Model(&models.DomainScanSummary{}).Where("overall_health = ?", "healthy").Count(&stats.HealthyDomains)
	h.db.Model(&models.DomainScanSummary{}).Where("overall_health = ?", "degraded").Count(&stats.DegradedDomains)
	h.db.Model(&models.DomainScanSummary{}).Where("overall_health = ?", "down").Count(&stats.DownDomains)

	c.JSON(http.StatusOK, stats)
}

// GetAPIQuotaStatus 获取 API 配额使用状态（管理员）
func (h *DomainScanHandler) GetAPIQuotaStatus(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	
	var quotas []models.APIQuota
	if err := h.db.Where("date = ?", today).Find(&quotas).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"google_safe_browsing": gin.H{"used": 0, "limit": 10000, "remaining": 10000},
			"virustotal":           gin.H{"used": 0, "limit": 500, "remaining": 500},
		})
		return
	}
	
	result := gin.H{}
	for _, quota := range quotas {
		result[quota.APIName] = gin.H{
			"used":      quota.UsedCount,
			"limit":     quota.DailyLimit,
			"remaining": quota.DailyLimit - quota.UsedCount,
			"date":      quota.Date,
		}
	}
	
	// 如果没有记录，返回默认值
	if _, ok := result["google_safe_browsing"]; !ok {
		result["google_safe_browsing"] = gin.H{"used": 0, "limit": 10000, "remaining": 10000}
	}
	if _, ok := result["virustotal"]; !ok {
		result["virustotal"] = gin.H{"used": 0, "limit": 500, "remaining": 500}
	}
	
	c.JSON(http.StatusOK, result)
}

// ListDomainScans 获取域名扫描记录列表（管理员）
func (h *DomainScanHandler) ListDomainScans(c *gin.Context) {
	domainID := c.Query("domain_id")
	scanType := c.Query("scan_type")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "50")
	
	var pageInt, pageSizeInt int
	if _, err := fmt.Sscanf(page, "%d", &pageInt); err != nil || pageInt < 1 {
		pageInt = 1
	}
	if _, err := fmt.Sscanf(pageSize, "%d", &pageSizeInt); err != nil || pageSizeInt < 1 {
		pageSizeInt = 50
	}
	if pageSizeInt > 200 {
		pageSizeInt = 200
	}
	
	query := h.db.Model(&models.DomainScan{}).Preload("Domain")
	
	if domainID != "" {
		query = query.Where("domain_id = ?", domainID)
	}
	if scanType != "" {
		query = query.Where("scan_type = ?", scanType)
	}
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count scans"})
		return
	}
	
	var scans []models.DomainScan
	offset := (pageInt - 1) * pageSizeInt
	if err := query.Order("scanned_at DESC").
		Offset(offset).
		Limit(pageSizeInt).
		Find(&scans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scans"})
		return
	}
	
	responses := make([]*models.DomainScanResponse, len(scans))
	for i, scan := range scans {
		responses[i] = scan.ToResponse()
	}
	
	totalPages := int((total + int64(pageSizeInt) - 1) / int64(pageSizeInt))
	
	c.JSON(http.StatusOK, gin.H{
		"scans": responses,
		"pagination": gin.H{
			"page":        pageInt,
			"page_size":   pageSizeInt,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// GetDomainScanSummaries 获取域名扫描摘要列表（管理员）
func (h *DomainScanHandler) GetDomainScanSummaries(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "50")
	search := c.Query("search")
	status := c.Query("status") // overall_health
	
	var pageInt, pageSizeInt int
	if _, err := fmt.Sscanf(page, "%d", &pageInt); err != nil || pageInt < 1 {
		pageInt = 1
	}
	if _, err := fmt.Sscanf(pageSize, "%d", &pageSizeInt); err != nil || pageSizeInt < 1 {
		pageSizeInt = 50
	}
	if pageSizeInt > 200 {
		pageSizeInt = 200
	}
	
	query := h.db.Model(&models.DomainScanSummary{}).Preload("Domain")
	
	if search != "" {
		query = query.Joins("JOIN domains ON domains.id = domain_scan_summaries.domain_id").
			Where("domains.full_domain LIKE ?", "%"+search+"%")
	}
	
	if status != "" {
		query = query.Where("overall_health = ?", status)
	}
	
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count summaries"})
		return
	}
	
	var summaries []models.DomainScanSummary
	offset := (pageInt - 1) * pageSizeInt
	if err := query.Order("last_scanned_at DESC NULLS LAST").
		Offset(offset).
		Limit(pageSizeInt).
		Find(&summaries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summaries"})
		return
	}
	
	responses := make([]*models.DomainHealthResponse, len(summaries))
	for i, summary := range summaries {
		responses[i] = summary.ToResponse()
	}
	
	totalPages := int((total + int64(pageSizeInt) - 1) / int64(pageSizeInt))
	
	c.JSON(http.StatusOK, gin.H{
		"summaries": responses,
		"pagination": gin.H{
			"page":        pageInt,
			"page_size":   pageSizeInt,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
