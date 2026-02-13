package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/middleware"
	"opendomain/internal/models"
	"opendomain/pkg/powerdns"
)

type DNSHandler struct {
	db   *gorm.DB
	cfg  *config.Config
	pdns *powerdns.Client
}

func NewDNSHandler(db *gorm.DB, cfg *config.Config) *DNSHandler {
	return &DNSHandler{
		db:   db,
		cfg:  cfg,
		pdns: powerdns.NewClient(cfg.PowerDNS.APIURL, cfg.PowerDNS.APIKey),
	}
}

// ListRecords 获取域名的 DNS 记录列表
func (h *DNSHandler) ListRecords(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("domainId")

	// 验证域名所有权
	var domain models.Domain
	if err := h.db.First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 获取 DNS 记录
	var records []models.DNSRecord
	if err := h.db.Where("domain_id = ?", domainID).Order("created_at DESC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DNS records"})
		return
	}

	responses := make([]*models.DNSRecordResponse, len(records))
	for i, record := range records {
		responses[i] = record.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{"records": responses})
}

// GetRecord 获取单个 DNS 记录
func (h *DNSHandler) GetRecord(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("domainId")
	recordID := c.Param("recordId")

	// 验证域名所有权
	var domain models.Domain
	if err := h.db.First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// 获取 DNS 记录
	var record models.DNSRecord
	if err := h.db.Where("id = ? AND domain_id = ?", recordID, domainID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS record not found"})
		return
	}

	c.JSON(http.StatusOK, record.ToResponse())
}

// CreateRecord 创建 DNS 记录
func (h *DNSHandler) CreateRecord(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("domainId")

	// 验证域名所有权
	var domain models.Domain
	if err := h.db.Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This domain has been suspended. All operations are disabled."})
		return
	}

	// 检查是否使用自定义 nameservers
	if !domain.UseDefaultNameservers {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot manage DNS records for domains using custom nameservers. Please manage DNS records on your custom nameserver.",
			"use_custom_ns": true,
		})
		return
	}

	var req models.DNSRecordCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置默认 TTL
	if req.TTL == 0 {
		req.TTL = 3600
	}

	// MX 记录默认优先级
	if req.Type == "MX" && req.Priority == nil {
		defaultPriority := 10
		req.Priority = &defaultPriority
	}

	// 验证记录内容
	if err := validateDNSRecord(req.Type, req.Content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// CNAME 冲突检查：CNAME 不能与同名的其他记录共存（仅检查活跃记录）
	var conflictCount int64
	if req.Type == "CNAME" {
		h.db.Model(&models.DNSRecord{}).Where("domain_id = ? AND name = ? AND type != ? AND is_active = ?", domain.ID, req.Name, "CNAME", true).Count(&conflictCount)
		if conflictCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "CNAME record cannot coexist with other record types at the same name"})
			return
		}
	} else {
		h.db.Model(&models.DNSRecord{}).Where("domain_id = ? AND name = ? AND type = ? AND is_active = ?", domain.ID, req.Name, "CNAME", true).Count(&conflictCount)
		if conflictCount > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot add this record because a CNAME record already exists at the same name"})
			return
		}
	}

	// 创建 DNS 记录
	record := &models.DNSRecord{
		DomainID:         domain.ID,
		Name:             req.Name,
		Type:             req.Type,
		Content:          req.Content,
		TTL:              req.TTL,
		Priority:         req.Priority,
		IsActive:         true,
		SyncedToPowerDNS: false,
	}

	if err := h.db.Create(record).Error; err != nil {
		fmt.Printf("Failed to create DNS record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create DNS record: %v", err)})
		return
	}

	// 同步到 PowerDNS
	go h.syncRecordSetToPowerDNS(record, &domain)

	c.JSON(http.StatusOK, gin.H{
		"message": "DNS record created successfully",
		"record":  record.ToResponse(),
	})
}

// UpdateRecord 更新 DNS 记录
func (h *DNSHandler) UpdateRecord(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("domainId")
	recordID := c.Param("recordId")

	// 验证域名所有权
	var domain models.Domain
	if err := h.db.First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This domain has been suspended. All operations are disabled."})
		return
	}

	// 获取 DNS 记录
	var record models.DNSRecord
	if err := h.db.Where("id = ? AND domain_id = ?", recordID, domainID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS record not found"})
		return
	}

	// 检查是否使用自定义 nameservers
	if !domain.UseDefaultNameservers {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot manage DNS records for domains using custom nameservers. Please manage DNS records on your custom nameserver.",
			"use_custom_ns": true,
		})
		return
	}

	var req models.DNSRecordUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if req.Name != nil {
		record.Name = *req.Name
	}
	if req.Type != nil {
		record.Type = *req.Type
	}
	if req.Content != nil {
		record.Content = *req.Content
	}
	if req.TTL != nil {
		record.TTL = *req.TTL
	}
	if req.Priority != nil {
		record.Priority = req.Priority
	}
	if req.IsActive != nil {
		record.IsActive = *req.IsActive
	}

	// 标记为未同步
	record.SyncedToPowerDNS = false

	if err := h.db.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DNS record"})
		return
	}

	// 同步到 PowerDNS
	h.db.Preload("RootDomain").First(&domain, domain.ID)
	go h.syncRecordSetToPowerDNS(&record, &domain)

	c.JSON(http.StatusOK, gin.H{
		"message": "DNS record updated successfully",
		"record":  record.ToResponse(),
	})
}

// DeleteRecord 删除 DNS 记录
func (h *DNSHandler) DeleteRecord(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("domainId")
	recordID := c.Param("recordId")

	// 验证域名所有权
	var domain models.Domain
	if err := h.db.First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This domain has been suspended. All operations are disabled."})
		return
	}

	// 获取记录信息（删除前）
	var record models.DNSRecord
	if err := h.db.Where("id = ? AND domain_id = ?", recordID, domainID).First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS record not found"})
		return
	}

	// 检查是否使用自定义 nameservers
	if !domain.UseDefaultNameservers {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot manage DNS records for domains using custom nameservers. Please manage DNS records on your custom nameserver.",
			"use_custom_ns": true,
		})
		return
	}

	// 删除 DNS 记录
	if err := h.db.Delete(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DNS record"})
		return
	}

	// 从 PowerDNS 删除记录
	h.db.Preload("RootDomain").First(&domain, domain.ID)
	go h.deleteRecordFromPowerDNS(&record, &domain)

	c.JSON(http.StatusOK, gin.H{"message": "DNS record deleted successfully"})
}

// validateDNSRecord 验证 DNS 记录内容
func validateDNSRecord(recordType, content string) error {
	// 基本验证，可以根据需要扩展
	switch recordType {
	case "A":
		// IPv4 地址验证
		if !isValidIPv4(content) {
			return fmt.Errorf("invalid IPv4 address")
		}
	case "AAAA":
		// IPv6 地址验证
		if !isValidIPv6(content) {
			return fmt.Errorf("invalid IPv6 address")
		}
	case "CNAME", "NS":
		// 域名验证
		if content == "" {
			return fmt.Errorf("content cannot be empty")
		}
	case "MX":
		// MX 记录验证
		if content == "" {
			return fmt.Errorf("content cannot be empty")
		}
	case "TXT":
		// TXT 记录可以是任何文本
		if content == "" {
			return fmt.Errorf("content cannot be empty")
		}
	}
	return nil
}

func isValidIPv4(ip string) bool {
	// 简单的 IPv4 验证
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

func isValidIPv6(ip string) bool {
	// 简单的 IPv6 验证
	return strings.Contains(ip, ":")
}

// buildRecordFQDN 构建 PowerDNS 所需的完整记录名
// Name="@" FullDomain="sub.example.com" → "sub.example.com"
// Name="www" FullDomain="sub.example.com" → "www.sub.example.com"
func buildRecordFQDN(name, fullDomain string) string {
	if name == "@" || name == "" {
		return fullDomain
	}
	return name + "." + fullDomain
}

// syncRecordSetToPowerDNS 同步同 name+type 的所有记录到 PowerDNS
func (h *DNSHandler) syncRecordSetToPowerDNS(record *models.DNSRecord, domain *models.Domain) {
	if domain.RootDomain == nil {
		return
	}
	zoneDomain := domain.RootDomain.Domain
	recordFQDN := buildRecordFQDN(record.Name, domain.FullDomain)

	// 查询同 domain+name+type 的所有活跃记录
	var allRecords []models.DNSRecord
	h.db.Where("domain_id = ? AND name = ? AND type = ? AND is_active = ?",
		record.DomainID, record.Name, record.Type, true).Find(&allRecords)

	entries := make([]powerdns.RecordEntry, 0, len(allRecords))
	for _, r := range allRecords {
		entries = append(entries, powerdns.RecordEntry{
			Content:  r.Content,
			Priority: r.Priority,
		})
	}

	if err := h.pdns.SetRecords(zoneDomain, recordFQDN, record.Type, entries, record.TTL); err != nil {
		syncErr := err.Error()
		for i := range allRecords {
			allRecords[i].SyncError = &syncErr
			allRecords[i].SyncedToPowerDNS = false
			h.db.Save(&allRecords[i])
		}
	} else {
		now := time.Now()
		for i := range allRecords {
			allRecords[i].SyncError = nil
			allRecords[i].SyncedToPowerDNS = true
			allRecords[i].LastSyncedAt = &now
			h.db.Save(&allRecords[i])
		}
	}

	h.updateDomainSyncStatus(record.DomainID)
}

// updateDomainSyncStatus 更新域名的 DNS 同步状态
func (h *DNSHandler) updateDomainSyncStatus(domainID uint) {
	var unsyncedCount int64
	h.db.Model(&models.DNSRecord{}).Where("domain_id = ? AND is_active = ? AND synced_to_powerdns = ?",
		domainID, true, false).Count(&unsyncedCount)

	h.db.Model(&models.Domain{}).Where("id = ?", domainID).
		Update("dns_synced", unsyncedCount == 0)
}

// deleteRecordFromPowerDNS 从 PowerDNS 删除记录（重新同步剩余记录）
func (h *DNSHandler) deleteRecordFromPowerDNS(record *models.DNSRecord, domain *models.Domain) {
	if domain.RootDomain == nil {
		return
	}
	zoneDomain := domain.RootDomain.Domain
	recordFQDN := buildRecordFQDN(record.Name, domain.FullDomain)

	// 查询同 name+type 的剩余活跃记录
	var remaining []models.DNSRecord
	h.db.Where("domain_id = ? AND name = ? AND type = ? AND is_active = ?",
		record.DomainID, record.Name, record.Type, true).Find(&remaining)

	if len(remaining) == 0 {
		// 没有剩余记录，删除整个 RRset
		if err := h.pdns.DeleteRRset(zoneDomain, recordFQDN, record.Type); err != nil {
			fmt.Printf("Warning: Failed to delete RRset from PowerDNS: %v\n", err)
		}
	} else {
		// 还有剩余记录，用剩余记录替换
		entries := make([]powerdns.RecordEntry, 0, len(remaining))
		for _, r := range remaining {
			entries = append(entries, powerdns.RecordEntry{
				Content:  r.Content,
				Priority: r.Priority,
			})
		}
		if err := h.pdns.SetRecords(zoneDomain, recordFQDN, record.Type, entries, remaining[0].TTL); err != nil {
			fmt.Printf("Warning: Failed to update RRset in PowerDNS: %v\n", err)
		}
	}

	h.updateDomainSyncStatus(record.DomainID)
}
// SyncFromPowerDNS 从 PowerDNS 同步所有 DNS 记录到数据库
func (h *DNSHandler) SyncFromPowerDNS(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("domainId")

	// 验证域名所有权
	var domain models.Domain
	if err := h.db.Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.RootDomain == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Root domain not found"})
		return
	}

	// 检查是否使用自定义 nameservers
	if !domain.UseDefaultNameservers {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot sync DNS records for domains using custom nameservers. DNS records are managed on your custom nameserver, not in our PowerDNS system.",
			"use_custom_ns": true,
		})
		return
	}

	// 从 PowerDNS 获取 zone 信息
	zone, err := h.pdns.GetZone(domain.RootDomain.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch zone from PowerDNS: %v", err)})
		return
	}

	// 统计信息
	syncStats := struct {
		Created int `json:"created"`
		Updated int `json:"updated"`
		Skipped int `json:"skipped"`
	}{}

	fullDomainWithDot := ensureTrailingDot(domain.FullDomain)

	// 遍历所有 RRsets
	for _, rrset := range zone.RRsets {
		// 跳过 SOA 和根域名的 NS 记录
		if rrset.Type == "SOA" {
			continue
		}
		if rrset.Type == "NS" && rrset.Name == fullDomainWithDot {
			continue
		}

		// 只处理属于当前域名的记录
		if rrset.Name != fullDomainWithDot && !strings.HasSuffix(rrset.Name, "."+fullDomainWithDot) {
			continue
		}

		// 解析记录名称（将 FQDN 转换为本地名称）
		recordName := extractRecordName(rrset.Name, domain.FullDomain)

		// 处理每条记录
		for _, record := range rrset.Records {
			if record.Disabled {
				continue // 跳过已禁用的记录
			}

			// 解析记录内容和优先级
			content, priority := parseRecordContent(rrset.Type, record.Content)

			// 检查是否已存在相同记录
			var existingRecord models.DNSRecord
			err := h.db.Where("domain_id = ? AND name = ? AND type = ? AND content = ?",
				domain.ID, recordName, rrset.Type, content).First(&existingRecord).Error

			now := time.Now()
			if err == gorm.ErrRecordNotFound {
				// 创建新记录
				newRecord := &models.DNSRecord{
					DomainID:         domain.ID,
					Name:             recordName,
					Type:             rrset.Type,
					Content:          content,
					TTL:              rrset.TTL,
					Priority:         priority,
					IsActive:         true,
					SyncedToPowerDNS: true,
					LastSyncedAt:     &now,
				}
				if err := h.db.Create(newRecord).Error; err != nil {
					fmt.Printf("Warning: Failed to create DNS record: %v\n", err)
					syncStats.Skipped++
				} else {
					syncStats.Created++
				}
			} else if err == nil {
				// 更新现有记录
				existingRecord.TTL = rrset.TTL
				existingRecord.Priority = priority
				existingRecord.IsActive = true
				existingRecord.SyncedToPowerDNS = true
				existingRecord.LastSyncedAt = &now
				existingRecord.SyncError = nil

				if err := h.db.Save(&existingRecord).Error; err != nil {
					fmt.Printf("Warning: Failed to update DNS record: %v\n", err)
					syncStats.Skipped++
				} else {
					syncStats.Updated++
				}
			} else {
				fmt.Printf("Warning: Failed to query DNS record: %v\n", err)
				syncStats.Skipped++
			}
		}
	}

	// 更新域名同步状态
	h.updateDomainSyncStatus(domain.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "DNS records synced from PowerDNS successfully",
		"stats":   syncStats,
	})
}

// extractRecordName 从 FQDN 提取记录名称
// "sub.example.com." + "example.com" -> "sub"
// "example.com." + "example.com" -> "@"
func extractRecordName(fqdn, fullDomain string) string {
	fqdn = strings.TrimSuffix(fqdn, ".")
	fullDomain = strings.TrimSuffix(fullDomain, ".")

	if fqdn == fullDomain {
		return "@"
	}

	if strings.HasSuffix(fqdn, "."+fullDomain) {
		return strings.TrimSuffix(fqdn, "."+fullDomain)
	}

	return fqdn
}

// parseRecordContent 解析记录内容，提取优先级（如果有）
func parseRecordContent(recordType, content string) (string, *int) {
	content = strings.TrimSuffix(content, ".")

	// MX 记录格式: "10 mail.example.com."
	if recordType == "MX" {
		parts := strings.Fields(content)
		if len(parts) >= 2 {
			if priority, err := strconv.Atoi(parts[0]); err == nil {
				return strings.TrimSuffix(parts[1], "."), &priority
			}
		}
	}

	return content, nil
}

// ensureTrailingDot 确保字符串以点结尾（用于 PowerDNS FQDN）
func ensureTrailingDot(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[len(s)-1] != '.' {
		return s + "."
	}
	return s
}
