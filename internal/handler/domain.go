package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
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

type DomainHandler struct {
	db   *gorm.DB
	cfg  *config.Config
	pdns *powerdns.Client
}

func NewDomainHandler(db *gorm.DB, cfg *config.Config) *DomainHandler {
	return &DomainHandler{
		db:   db,
		cfg:  cfg,
		pdns: powerdns.NewClient(cfg.PowerDNS.APIURL, cfg.PowerDNS.APIKey),
	}
}

// ListRootDomains 获取根域名列表
// @Summary 获取根域名列表
// @Tags Public
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/public/root-domains [get]
func (h *DomainHandler) ListRootDomains(c *gin.Context) {
	var rootDomains []models.RootDomain
	if err := h.db.Where("is_active = ?", true).Order("priority DESC, id ASC").Find(&rootDomains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": middleware.T(c, "error.internal_server")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"root_domains": rootDomains})
}

// SearchDomain 搜索域名可用性
// @Summary 搜索域名可用性
// @Tags Domain
// @Produce json
// @Param subdomain query string true "子域名"
// @Param root_domain_id query int true "根域名ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/domains/search [get]
// @Security Bearer
func (h *DomainHandler) SearchDomain(c *gin.Context) {
	var req models.DomainSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": middleware.T(c, "error.validation")})
		return
	}

	// 获取根域名
	var rootDomain models.RootDomain
	if err := h.db.First(&rootDomain, req.RootDomainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": middleware.T(c, "error.root_domain_not_found")})
		return
	}

	if !rootDomain.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": middleware.T(c, "error.root_domain_inactive")})
		return
	}

	// 验证子域名格式
	if !isValidSubdomain(req.Subdomain) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     middleware.T(c, "error.domain_invalid"),
			"available": false,
		})
		return
	}

	// 检查长度
	if len(req.Subdomain) < rootDomain.MinLength || len(req.Subdomain) > rootDomain.MaxLength {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Subdomain length must be between %d and %d characters",
				rootDomain.MinLength, rootDomain.MaxLength),
			"available": false,
		})
		return
	}

	// 检查黑名单
	if isBlacklisted(req.Subdomain) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":     "This subdomain is reserved",
			"available": false,
		})
		return
	}

	fullDomain := fmt.Sprintf("%s.%s", req.Subdomain, rootDomain.Domain)

	// 检查是否已被注册
	var existingDomain models.Domain
	err := h.db.Where("full_domain = ?", fullDomain).First(&existingDomain).Error

	available := err == gorm.ErrRecordNotFound

	// 检查是否在待激活列表中
	reserved := false
	if available {
		var pendingDomain models.PendingDomain
		if err := h.db.Where("full_domain = ? AND deleted_at IS NULL", fullDomain).First(&pendingDomain).Error; err == nil {
			available = false
			reserved = true
		}
	}

	response := gin.H{
		"available":   available,
		"subdomain":   req.Subdomain,
		"root_domain": rootDomain.Domain,
		"full_domain": fullDomain,
	}

	if reserved {
		response["reserved"] = true
		response["message"] = "This domain is reserved and can only be activated by syncing from FOSSBilling"
	}

	c.JSON(http.StatusOK, response)
}

// RegisterDomain 注册域名（仅适用于免费域名）
// @Summary 注册域名
// @Tags Domain
// @Accept json
// @Produce json
// @Param request body models.DomainRegisterRequest true "注册信息"
// @Success 200 {object} models.DomainResponse
// @Router /api/domains [post]
// @Security Bearer
func (h *DomainHandler) RegisterDomain(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.DomainRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户信息
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 获取根域名
	var rootDomain models.RootDomain
	if err := h.db.First(&rootDomain, req.RootDomainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Root domain not found"})
		return
	}

	if !rootDomain.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Root domain is not active"})
		return
	}

	// 检查是否为付费域名
	if !rootDomain.IsFree {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":            "This domain requires payment. Please create an order first.",
			"requires_payment": true,
			"price_per_year":   rootDomain.PricePerYear,
			"lifetime_price":   rootDomain.LifetimePrice,
		})
		return
	}

	// 验证域名
	if !isValidSubdomain(req.Subdomain) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subdomain format"})
		return
	}

	if isBlacklisted(req.Subdomain) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This subdomain is reserved"})
		return
	}

	fullDomain := fmt.Sprintf("%s.%s", req.Subdomain, rootDomain.Domain)

	// 检查是否已存在
	var existingDomain models.Domain
	if err := h.db.Where("full_domain = ?", fullDomain).First(&existingDomain).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Domain already registered"})
		return
	}

	// 检查是否在待激活列表中（从FOSSBilling预同步的域名）
	var pendingDomain models.PendingDomain
	if err := h.db.Where("full_domain = ? AND deleted_at IS NULL", fullDomain).First(&pendingDomain).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":    "This domain is reserved and can only be activated by syncing from FOSSBilling",
			"reserved": true,
		})
		return
	}

	// 检查用户配额
	var userDomainCount int64
	h.db.Model(&models.Domain{}).Where("user_id = ? AND status = ?", userID, "active").Count(&userDomainCount)

	var usedCoupon *models.Coupon
	quotaExceeded := int(userDomainCount) >= user.DomainQuota

	// 如果配额超限，检查优惠券
	if quotaExceeded {
		if req.CouponCode == nil || *req.CouponCode == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Domain quota exceeded. Please use a coupon or purchase a domain.",
				"quota_exceeded": true,
			})
			return
		}

		// 查找优惠券
		var coupon models.Coupon
		if err := h.db.Where("UPPER(code) = UPPER(?)", *req.CouponCode).First(&coupon).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
			return
		}

		// 验证优惠券
		if err := h.validateCouponForFreeRegistration(&coupon, userID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 只接受 quota_increase 类型
		if coupon.DiscountType != "quota_increase" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "This coupon type cannot be used for free domain registration",
			})
			return
		}

		usedCoupon = &coupon
	}

	// 开始事务注册域名
	err := h.db.Transaction(func(tx *gorm.DB) error {
		// 创建域名
		now := time.Now()
		domain := &models.Domain{
			UserID:                userID,
			RootDomainID:          req.RootDomainID,
			Subdomain:             req.Subdomain,
			FullDomain:            fullDomain,
			Status:                "active",
			RegisteredAt:          now,
			ExpiresAt:             now.AddDate(1, 0, 0), // 1 年后过期
			AutoRenew:             false,
			Nameservers:           rootDomain.Nameservers,
			UseDefaultNameservers: rootDomain.UseDefaultNameservers,
			DNSSynced:             false,
		}

		if err := tx.Create(domain).Error; err != nil {
			return err
		}

		// 如果使用了优惠券，记录使用
		if usedCoupon != nil {
			usage := &models.CouponUsage{
				CouponID: usedCoupon.ID,
				UserID:   userID,
				DomainID: &domain.ID,
			}
			if err := tx.Create(usage).Error; err != nil {
				return err
			}

			// 更新优惠券使用次数
			if err := tx.Model(&models.Coupon{}).
				Where("id = ?", usedCoupon.ID).
				UpdateColumn("used_count", gorm.Expr("used_count + ?", 1)).Error; err != nil {
				return err
			}

			// 如果是 quota_increase 类型，更新用户配额
			if usedCoupon.DiscountType == "quota_increase" && usedCoupon.QuotaIncrease > 0 {
				if err := tx.Model(&user).
					UpdateColumn("domain_quota", gorm.Expr("domain_quota + ?", usedCoupon.QuotaIncrease)).Error; err != nil {
					return err
				}
			}
		}

		// 更新根域名注册计数
		if err := tx.Model(&rootDomain).
			UpdateColumn("registration_count", gorm.Expr("registration_count + ?", 1)).Error; err != nil {
			return err
		}

		// 重新加载域名关联数据
		return tx.Preload("RootDomain").First(domain, domain.ID).Error
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register domain"})
		return
	}

	// 获取最新的域名数据
	var domain models.Domain
	h.db.Preload("RootDomain").Where("full_domain = ?", fullDomain).First(&domain)

	// 在 PowerDNS 中配置域名
	if domain.RootDomain != nil {
		var nameservers []string
		if err := json.Unmarshal([]byte(domain.Nameservers), &nameservers); err == nil {
			// 如果使用默认 NS 或自定义 NS，都调用此函数进行配置
			go h.updateDomainNSRecordsInPowerDNS(&domain, nameservers, domain.UseDefaultNameservers)
		}
	}

	response := gin.H{
		"message": "Domain registered successfully",
		"domain":  domain.ToResponse(),
	}

	if usedCoupon != nil {
		response["coupon_applied"] = true
		response["coupon_code"] = usedCoupon.Code
	}

	c.JSON(http.StatusOK, response)
}

// validateCouponForFreeRegistration 验证免费注册时的优惠券
func (h *DomainHandler) validateCouponForFreeRegistration(coupon *models.Coupon, userID uint) error {
	now := time.Now()

	// 检查是否激活
	if !coupon.IsActive {
		return fmt.Errorf("coupon is not active")
	}

	// 检查有效期
	if now.Before(coupon.ValidFrom) {
		return fmt.Errorf("coupon not yet valid")
	}
	if coupon.ValidUntil != nil && now.After(*coupon.ValidUntil) {
		return fmt.Errorf("coupon has expired")
	}

	// 检查使用次数
	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return fmt.Errorf("coupon usage limit reached")
	}

	// 检查用户是否已使用（quota 类型的优惠券每人只能用一次）
	var usage models.CouponUsage
	if err := h.db.Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).
		First(&usage).Error; err == nil {
		return fmt.Errorf("you have already used this coupon")
	}

	return nil
}

// ListMyDomains 获取我的域名列表
// @Summary 获取我的域名列表
// @Tags Domain
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/domains [get]
// @Security Bearer
func (h *DomainHandler) ListMyDomains(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var domains []models.Domain
	if err := h.db.Preload("RootDomain").Where("user_id = ?", userID).
		Order("created_at DESC").Find(&domains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains"})
		return
	}

	// 获取域名的扫描摘要
	domainIDs := make([]uint, len(domains))
	for i, domain := range domains {
		domainIDs[i] = domain.ID
	}

	var summaries []models.DomainScanSummary
	if len(domainIDs) > 0 {
		h.db.Where("domain_id IN ?", domainIDs).Find(&summaries)
	}

	// 创建域名ID到扫描摘要的映射
	summaryMap := make(map[uint]*models.DomainScanSummary)
	for i := range summaries {
		summaryMap[summaries[i].DomainID] = &summaries[i]
	}

	responses := make([]map[string]interface{}, len(domains))
	for i, domain := range domains {
		resp := domain.ToResponse()
		domainMap := map[string]interface{}{
			"id":                      resp.ID,
			"user_id":                 resp.UserID,
			"root_domain_id":          resp.RootDomainID,
			"subdomain":               resp.Subdomain,
			"full_domain":             resp.FullDomain,
			"status":                  resp.Status,
			"registered_at":           resp.RegisteredAt,
			"expires_at":              resp.ExpiresAt,
			"auto_renew":              resp.AutoRenew,
			"nameservers":             resp.Nameservers,
			"use_default_nameservers": resp.UseDefaultNameservers,
			"dns_synced":              resp.DNSSynced,
			"root_domain":             resp.RootDomain,
		}

		// 添加扫描状态
		if summary, ok := summaryMap[domain.ID]; ok {
			domainMap["scan_summary"] = map[string]interface{}{
				"overall_health":       summary.OverallHealth,
				"http_status":          summary.HTTPStatus,
				"dns_status":           summary.DNSStatus,
				"ssl_status":           summary.SSLStatus,
				"safe_browsing_status": summary.SafeBrowsingStatus,
				"virustotal_status":    summary.VirusTotalStatus,
				"last_scanned_at":      summary.LastScannedAt,
				"total_scans":          summary.TotalScans,
				"successful_scans":     summary.SuccessfulScans,
			}
			if summary.TotalScans > 0 {
				uptime := float64(summary.SuccessfulScans) / float64(summary.TotalScans) * 100
				domainMap["scan_summary"].(map[string]interface{})["uptime_percentage"] = uptime
			}
		} else {
			domainMap["scan_summary"] = nil
		}

		responses[i] = domainMap
	}

	c.JSON(http.StatusOK, gin.H{"domains": responses})
}

// GetDomain 获取域名详情
func (h *DomainHandler) GetDomain(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("id")

	var domain models.Domain
	if err := h.db.Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	// 验证所有权
	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, domain.ToResponse())
}

// DeleteDomain 删除域名
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("id")

	var domain models.Domain
	if err := h.db.Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	// 验证所有权
	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This domain has been suspended. All operations are disabled."})
		return
	}

	// 删除所有 DNS 记录（从数据库和 PowerDNS）
	if err := h.deleteAllDNSRecordsForDomain(&domain); err != nil {
		fmt.Printf("Warning: Failed to delete DNS records for domain %s: %v\n", domain.FullDomain, err)
		// 继续删除域名，即使 DNS 清理失败
	}

	// 软删除域名
	if err := h.db.Delete(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Domain deleted successfully"})
}

// ModifyNameservers 修改域名 Nameservers
func (h *DomainHandler) ModifyNameservers(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("id")

	var domain models.Domain
	if err := h.db.Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	// 验证所有权
	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This domain has been suspended. All operations are disabled."})
		return
	}

	var req struct {
		Nameservers []string `json:"nameservers" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证 nameservers 格式
	for _, ns := range req.Nameservers {
		if ns == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nameserver cannot be empty"})
			return
		}
	}

	// 更新域名的 nameservers（存储为 JSON）
	nameserversJSON, err := json.Marshal(req.Nameservers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode nameservers"})
		return
	}

	// 判断是否为默认 NS
	defaultNS := []string{h.cfg.DNS.DefaultNS1, h.cfg.DNS.DefaultNS2}
	isDefault := len(req.Nameservers) == 2 &&
		req.Nameservers[0] == defaultNS[0] &&
		req.Nameservers[1] == defaultNS[1]

	if err := h.db.Model(&domain).Updates(map[string]interface{}{
		"nameservers":             string(nameserversJSON),
		"use_default_nameservers": isDefault,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update nameservers"})
		return
	}

	// 在 PowerDNS 中更新 NS 记录
	if domain.RootDomain != nil {
		go h.updateDomainNSRecordsInPowerDNS(&domain, req.Nameservers, isDefault)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Nameservers updated successfully",
		"nameservers": req.Nameservers,
	})
}

// RenewDomain 续费域名
func (h *DomainHandler) RenewDomain(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("id")

	var domain models.Domain
	if err := h.db.Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	// 验证所有权
	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This domain has been suspended. All operations are disabled."})
		return
	}

	var req struct {
		Years      int     `json:"years"`
		IsLifetime bool    `json:"is_lifetime"`
		CouponCode *string `json:"coupon_code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证参数
	if !req.IsLifetime && (req.Years < 1 || req.Years > 10) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Years must be between 1 and 10"})
		return
	}

	// 如果是付费域名，需要创建续费订单
	if !domain.RootDomain.IsFree {
		// 计算价格
		var basePrice float64
		if req.IsLifetime {
			if domain.RootDomain.LifetimePrice == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Lifetime pricing not available"})
				return
			}
			basePrice = *domain.RootDomain.LifetimePrice
		} else {
			if domain.RootDomain.PricePerYear == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Pricing not configured"})
				return
			}
			basePrice = *domain.RootDomain.PricePerYear * float64(req.Years)
		}
		
		// 应用优惠券
		var discountAmount float64
		var couponID *uint
		var couponCode *string

		if req.CouponCode != nil && *req.CouponCode != "" {
			var coupon models.Coupon
			if err := h.db.Where("UPPER(code) = UPPER(?)", *req.CouponCode).First(&coupon).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
				return
			}

			// 验证优惠券
			if err := h.validateCouponForOrder(&coupon, userID); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// 应用折扣
			if coupon.DiscountType == "percentage" && coupon.DiscountValue != nil {
				discountAmount = basePrice * (*coupon.DiscountValue / 100.0)
			} else if coupon.DiscountType == "fixed" && coupon.DiscountValue != nil {
				discountAmount = math.Min(*coupon.DiscountValue, basePrice)
			} else if coupon.DiscountType == "quota_increase" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "This coupon type cannot be used for renewals",
				})
				return
			}

			couponID = &coupon.ID
			couponCode = &coupon.Code
		}

		finalPrice := math.Max(0, basePrice-discountAmount)

		// 生成订单号
		orderNumber := h.generateOrderNumber()

		// 创建续费订单（15分钟后过期）
		order := &models.Order{
			OrderNumber:    orderNumber,
			UserID:         userID,
			RootDomainID:   domain.RootDomainID,
			Subdomain:      domain.Subdomain,
			FullDomain:     domain.FullDomain,
			DomainID:       &domain.ID,
			Years:          req.Years,
			IsLifetime:     req.IsLifetime,
			BasePrice:      basePrice,
			DiscountAmount: discountAmount,
			FinalPrice:     finalPrice,
			Status:         "pending",
			CouponID:       couponID,
			CouponCode:     couponCode,
			ExpiresAt:      time.Now().Add(15 * time.Minute),
		}

		if err := h.db.Create(order).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":          "Renewal order created successfully",
			"requires_payment": true,
			"order_id":         order.ID,
			"order_number":     order.OrderNumber,
			"base_price":       basePrice,
			"discount_amount":  discountAmount,
			"final_price":      finalPrice,
			"coupon_applied":   couponID != nil,
			"years":            req.Years,
		})
		return
	}

	// 免费域名直接续费
	newExpiry := domain.ExpiresAt.AddDate(req.Years, 0, 0)
	if err := h.db.Model(&domain).Update("expires_at", newExpiry).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to renew domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Domain renewed successfully",
		"new_expiry": newExpiry,
		"years":      req.Years,
	})
}

// validateCouponForOrder 验证订单优惠券
func (h *DomainHandler) validateCouponForOrder(coupon *models.Coupon, userID uint) error {
	now := time.Now()

	// 检查是否激活
	if !coupon.IsActive {
		return fmt.Errorf("coupon is not active")
	}

	// 检查有效期
	if now.Before(coupon.ValidFrom) {
		return fmt.Errorf("coupon not yet valid")
	}
	if coupon.ValidUntil != nil && now.After(*coupon.ValidUntil) {
		return fmt.Errorf("coupon has expired")
	}

	// 检查使用次数
	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return fmt.Errorf("coupon usage limit reached")
	}

	// 如果优惠券不可重复使用，检查用户是否已使用
	if !coupon.IsReusable {
		var usage models.CouponUsage
		if err := h.db.Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).
			First(&usage).Error; err == nil {
			return fmt.Errorf("you have already used this coupon")
		}
	}

	return nil
}

// generateOrderNumber 生成订单号
func (h *DomainHandler) generateOrderNumber() string {
	// 格式: ORD + 时间戳 + 随机字符
	timestamp := time.Now().Unix()
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("ORD%d%s", timestamp, randomStr[:8])
}


// TransferDomain 转移域名（站内转移）
func (h *DomainHandler) TransferDomain(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	domainID := c.Param("id")

	var domain models.Domain
	if err := h.db.First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	// 验证所有权
	if domain.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if domain.Status == "suspended" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This domain has been suspended. All operations are disabled."})
		return
	}

	var req struct {
		Target string `json:"target" binding:"required"` // Email or username
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 查找目标用户（通过 email 或 username）
	var targetUser models.User
	if err := h.db.Where("email = ? OR username = ?", req.Target, req.Target).First(&targetUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Target user not found"})
		return
	}

	// 不能转移给自己
	if targetUser.ID == userID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot transfer domain to yourself"})
		return
	}

	// 执行转移
	if err := h.db.Model(&domain).Update("user_id", targetUser.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to transfer domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Domain transferred successfully",
		"new_owner": targetUser.Email,
		"domain":    domain.FullDomain,
	})
}

// ListAllDomains 管理员：获取所有域名列表
func (h *DomainHandler) ListAllDomains(c *gin.Context) {
	search := c.Query("search")
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	query := h.db.Model(&models.Domain{})

	// 搜索功能：支持按域名、子域名、用户名、邮箱搜索
	if search != "" {
		query = query.Joins("LEFT JOIN users ON users.id = domains.user_id").
			Where("domains.full_domain LIKE ? OR domains.subdomain LIKE ? OR users.username LIKE ? OR users.email LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count domains"})
		return
	}

	// 分页查询（重新构建query以包含Preload）
	var domains []models.Domain
	offset := (page - 1) * pageSize
	domainQuery := h.db.Preload("RootDomain").Preload("User")
	if search != "" {
		domainQuery = domainQuery.Joins("LEFT JOIN users ON users.id = domains.user_id").
			Where("domains.full_domain LIKE ? OR domains.subdomain LIKE ? OR users.username LIKE ? OR users.email LIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := domainQuery.Order("domains.created_at DESC").Offset(offset).Limit(pageSize).Find(&domains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains"})
		return
	}

	responses := make([]*models.DomainResponse, len(domains))
	for i, domain := range domains {
		responses[i] = domain.ToResponse()
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"domains": responses,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// isValidSubdomain 验证子域名格式
func isValidSubdomain(subdomain string) bool {
	// 长度检查
	if len(subdomain) < 3 || len(subdomain) > 63 {
		return false
	}

	// 格式检查：字母、数字、连字符，不能以连字符开头或结尾
	pattern := `^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`
	matched, _ := regexp.MatchString(pattern, subdomain)
	return matched
}

// isBlacklisted 检查是否在黑名单中
func isBlacklisted(subdomain string) bool {
	blacklist := []string{
		"admin", "root", "api", "www", "mail", "smtp", "ftp", "ssh", "dns",
		"test", "demo", "dev", "stage", "prod", "blog", "forum", "shop",
		"status", "support", "help", "docs", "cdn", "static", "assets",
	}

	for _, word := range blacklist {
		if subdomain == word {
			return true
		}
	}
	return false
}

// Admin Root Domain Management

// ListAllRootDomains 管理员：获取所有根域名列表
func (h *DomainHandler) ListAllRootDomains(c *gin.Context) {
	var rootDomains []models.RootDomain
	if err := h.db.Order("priority DESC, id ASC").Find(&rootDomains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch root domains"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"root_domains": rootDomains})
}

// CreateRootDomain 管理员：创建根域名
func (h *DomainHandler) CreateRootDomain(c *gin.Context) {
	var req struct {
		Domain                string   `json:"domain" binding:"required"`
		Description           *string  `json:"description"`
		Priority              int      `json:"priority"`
		IsActive              bool     `json:"is_active"`
		IsHot                 bool     `json:"is_hot"`
		IsNew                 bool     `json:"is_new"`
		IsFree                bool     `json:"is_free"`
		PricePerYear          *float64 `json:"price_per_year"`
		LifetimePrice         *float64 `json:"lifetime_price"`
		UseDefaultNameservers bool     `json:"use_default_nameservers"`
		Nameservers           []string `json:"nameservers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查域名是否已存在
	var existing models.RootDomain
	if err := h.db.Where("domain = ?", req.Domain).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Root domain already exists"})
		return
	}

	// 设置 nameservers
	var nameserversJSON string
	if req.UseDefaultNameservers {
		// 使用默认 NS
		defaultNS := []string{h.cfg.DNS.DefaultNS1, h.cfg.DNS.DefaultNS2}
		nsBytes, err := json.Marshal(defaultNS)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode default nameservers"})
			return
		}
		nameserversJSON = string(nsBytes)
	} else {
		// 使用自定义 NS
		if len(req.Nameservers) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Nameservers required when not using default"})
			return
		}
		nsBytes, err := json.Marshal(req.Nameservers)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode nameservers"})
			return
		}
		nameserversJSON = string(nsBytes)
	}

	rootDomain := &models.RootDomain{
		Domain:                req.Domain,
		Description:           req.Description,
		Priority:              req.Priority,
		IsActive:              req.IsActive,
		IsHot:                 req.IsHot,
		IsNew:                 req.IsNew,
		IsFree:                req.IsFree,
		PricePerYear:          req.PricePerYear,
		LifetimePrice:         req.LifetimePrice,
		UseDefaultNameservers: req.UseDefaultNameservers,
		Nameservers:           nameserversJSON,
	}

	if err := h.db.Create(rootDomain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create root domain"})
		return
	}

	// 在 PowerDNS 中创建 Zone
	var nameservers []string
	if err := json.Unmarshal([]byte(nameserversJSON), &nameservers); err == nil {
		if err := h.pdns.CreateZone(req.Domain, nameservers); err != nil {
			fmt.Printf("Warning: Failed to create PowerDNS zone for %s: %v\n", req.Domain, err)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Root domain created successfully",
		"root_domain": rootDomain,
	})
}

// UpdateRootDomain 管理员：更新根域名
func (h *DomainHandler) UpdateRootDomain(c *gin.Context) {
	id := c.Param("id")

	var rootDomain models.RootDomain
	if err := h.db.First(&rootDomain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Root domain not found"})
		return
	}

	var req struct {
		Description           *string  `json:"description"`
		Priority              *int     `json:"priority"`
		IsActive              *bool    `json:"is_active"`
		IsHot                 *bool    `json:"is_hot"`
		IsNew                 *bool    `json:"is_new"`
		IsFree                *bool    `json:"is_free"`
		PricePerYear          *float64 `json:"price_per_year"`
		LifetimePrice         *float64 `json:"lifetime_price"`
		UseDefaultNameservers *bool    `json:"use_default_nameservers"`
		Nameservers           []string `json:"nameservers"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	updates := make(map[string]interface{})
	selectFields := []string{}

	if req.Description != nil {
		updates["description"] = *req.Description
		selectFields = append(selectFields, "description")
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
		selectFields = append(selectFields, "priority")
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
		selectFields = append(selectFields, "is_active")
	}
	if req.IsHot != nil {
		updates["is_hot"] = *req.IsHot
		selectFields = append(selectFields, "is_hot")
	}
	if req.IsNew != nil {
		updates["is_new"] = *req.IsNew
		selectFields = append(selectFields, "is_new")
	}
	if req.IsFree != nil {
		updates["is_free"] = *req.IsFree
		selectFields = append(selectFields, "is_free")
	}
	if req.PricePerYear != nil {
		updates["price_per_year"] = *req.PricePerYear
		selectFields = append(selectFields, "price_per_year")
	}
	if req.LifetimePrice != nil {
		updates["lifetime_price"] = *req.LifetimePrice
		selectFields = append(selectFields, "lifetime_price")
	}
	if req.UseDefaultNameservers != nil {
		updates["use_default_nameservers"] = *req.UseDefaultNameservers
		selectFields = append(selectFields, "use_default_nameservers")

		// 如果切换到使用默认 NS，更新 nameservers
		if *req.UseDefaultNameservers {
			defaultNS := []string{h.cfg.DNS.DefaultNS1, h.cfg.DNS.DefaultNS2}
			nsBytes, err := json.Marshal(defaultNS)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode default nameservers"})
				return
			}
			updates["nameservers"] = string(nsBytes)
			selectFields = append(selectFields, "nameservers")
		} else {
			// 如果切换到自定义 NS，需要提供 nameservers
			if len(req.Nameservers) > 0 {
				nsBytes, err := json.Marshal(req.Nameservers)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encode nameservers"})
					return
				}
				updates["nameservers"] = string(nsBytes)
				selectFields = append(selectFields, "nameservers")
			}
		}
	}

	if err := h.db.Model(&rootDomain).Select(selectFields).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update root domain"})
		return
	}

	// Reload the root domain to get fresh data
	if err := h.db.First(&rootDomain, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reload root domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Root domain updated successfully",
		"root_domain": rootDomain,
	})
}

// DeleteRootDomain 管理员：删除根域名
func (h *DomainHandler) DeleteRootDomain(c *gin.Context) {
	id := c.Param("id")

	var rootDomain models.RootDomain
	if err := h.db.First(&rootDomain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Root domain not found"})
		return
	}

	// 检查是否有子域名在使用
	var count int64
	h.db.Model(&models.Domain{}).Where("root_domain_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Cannot delete: %d domains are using this root domain", count)})
		return
	}

	if err := h.db.Delete(&rootDomain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete root domain"})
		return
	}

	// 从 PowerDNS 删除 Zone
	if err := h.pdns.DeleteZone(rootDomain.Domain); err != nil {
		fmt.Printf("Warning: Failed to delete PowerDNS zone for %s: %v\n", rootDomain.Domain, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Root domain deleted successfully"})
}

// ListDomainsByRootDomain 管理员：获取某个根域名下的所有域名
func (h *DomainHandler) ListDomainsByRootDomain(c *gin.Context) {
	rootDomainID := c.Param("id")
	search := c.Query("search")

	// 验证根域名存在
	var rootDomain models.RootDomain
	if err := h.db.First(&rootDomain, rootDomainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Root domain not found"})
		return
	}

	query := h.db.Preload("User").Preload("RootDomain").
		Where("root_domain_id = ?", rootDomainID)

	if search != "" {
		query = query.Where("subdomain ILIKE ? OR full_domain ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	var domains []models.Domain
	if err := query.Order("created_at DESC").Find(&domains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains"})
		return
	}

	responses := make([]*models.DomainResponse, len(domains))
	for i, domain := range domains {
		responses[i] = domain.ToResponse()
	}

	c.JSON(http.StatusOK, gin.H{
		"domains":     responses,
		"root_domain": rootDomain,
	})
}

// AdminUpdateDomainStatus 管理员：更新域名状态
func (h *DomainHandler) AdminUpdateDomainStatus(c *gin.Context) {
	domainID := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required,oneof=active suspended"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var domain models.Domain
	if err := h.db.Preload("User").Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	if err := h.db.Model(&domain).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update domain status"})
		return
	}

	// 同步 PowerDNS：暂停时 disable 所有记录，激活时 enable
	if domain.RootDomain != nil {
		disabled := req.Status == "suspended"
		go func() {
			if err := h.pdns.SetSubdomainDisabled(domain.RootDomain.Domain, domain.FullDomain, disabled); err != nil {
				fmt.Printf("Warning: Failed to %s DNS records in PowerDNS for %s: %v\n",
					map[bool]string{true: "disable", false: "enable"}[disabled], domain.FullDomain, err)
			}
		}()
	}

	domain.Status = req.Status
	c.JSON(http.StatusOK, gin.H{
		"message": "Domain status updated",
		"domain":  domain.ToResponse(),
	})
}

// deleteAllDNSRecordsForDomain 删除域名的所有 DNS 记录（从数据库和 PowerDNS）
func (h *DomainHandler) deleteAllDNSRecordsForDomain(domain *models.Domain) error {
	// 查询所有 DNS 记录
	var records []models.DNSRecord
	if err := h.db.Where("domain_id = ?", domain.ID).Find(&records).Error; err != nil {
		return fmt.Errorf("failed to query DNS records: %w", err)
	}

	// 如果没有记录，直接返回
	if len(records) == 0 {
		return nil
	}

	// 如果有 RootDomain，从 PowerDNS 删除记录
	if domain.RootDomain != nil {
		zoneDomain := domain.RootDomain.Domain

		// 按 name+type 分组记录，以便从 PowerDNS 删除
		recordGroups := make(map[string]models.DNSRecord)
		for _, record := range records {
			key := fmt.Sprintf("%s|%s", record.Name, record.Type)
			if _, exists := recordGroups[key]; !exists {
				recordGroups[key] = record
			}
		}

		// 删除每个 RRset
		for _, record := range recordGroups {
			recordFQDN := buildRecordFQDN(record.Name, domain.FullDomain)
			if err := h.pdns.DeleteRRset(zoneDomain, recordFQDN, record.Type); err != nil {
				fmt.Printf("Warning: Failed to delete RRset %s/%s from PowerDNS: %v\n", recordFQDN, record.Type, err)
				// 继续删除其他记录
			}
		}
	}

	// 从数据库删除所有 DNS 记录
	if err := h.db.Where("domain_id = ?", domain.ID).Delete(&models.DNSRecord{}).Error; err != nil {
		return fmt.Errorf("failed to delete DNS records from database: %w", err)
	}

	return nil
}

// AdminDeleteDomain 管理员：删除域名
func (h *DomainHandler) AdminDeleteDomain(c *gin.Context) {
	domainID := c.Param("id")

	var domain models.Domain
	if err := h.db.Preload("RootDomain").First(&domain, domainID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Domain not found"})
		return
	}

	// 删除所有 DNS 记录（从数据库和 PowerDNS）
	if err := h.deleteAllDNSRecordsForDomain(&domain); err != nil {
		fmt.Printf("Warning: Failed to delete DNS records for domain %s: %v\n", domain.FullDomain, err)
		// 继续删除域名，即使 DNS 清理失败
	}

	// 软删除域名
	if err := h.db.Delete(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete domain"})
		return
	}

	// 减少根域名注册计数
	h.db.Model(&models.RootDomain{}).Where("id = ?", domain.RootDomainID).
		UpdateColumn("registration_count", gorm.Expr("GREATEST(registration_count - 1, 0)"))

	c.JSON(http.StatusOK, gin.H{"message": "Domain deleted successfully"})
}

// CleanupExpiredDomains 自动清理过期域名
// 删除过期超过指定天数的域名及其 DNS 记录
func (h *DomainHandler) CleanupExpiredDomains(daysAfterExpiry int) {
	if daysAfterExpiry <= 0 {
		daysAfterExpiry = 30 // 默认过期 30 天后删除
	}

	cutoffTime := time.Now().AddDate(0, 0, -daysAfterExpiry)

	fmt.Printf("Starting cleanup of domains expired before %s...\n", cutoffTime.Format("2006-01-02 15:04:05"))

	// 查找过期超过指定天数的域名
	var expiredDomains []models.Domain
	if err := h.db.Preload("RootDomain").
		Where("expires_at < ? AND status != ?", cutoffTime, "suspended").
		Find(&expiredDomains).Error; err != nil {
		fmt.Printf("Error querying expired domains: %v\n", err)
		return
	}

	if len(expiredDomains) == 0 {
		fmt.Println("No expired domains to clean up.")
		return
	}

	fmt.Printf("Found %d expired domains to clean up.\n", len(expiredDomains))

	successCount := 0
	failCount := 0

	for _, domain := range expiredDomains {
		fmt.Printf("Cleaning up domain: %s (expired at %s)\n",
			domain.FullDomain, domain.ExpiresAt.Format("2006-01-02"))

		// 删除所有 DNS 记录
		if err := h.deleteAllDNSRecordsForDomain(&domain); err != nil {
			fmt.Printf("Warning: Failed to delete DNS records for domain %s: %v\n", domain.FullDomain, err)
			// 继续删除域名，即使 DNS 清理失败
		}

		// 软删除域名
		if err := h.db.Delete(&domain).Error; err != nil {
			fmt.Printf("Error: Failed to delete domain %s: %v\n", domain.FullDomain, err)
			failCount++
			continue
		}

		// 减少根域名注册计数
		if domain.RootDomainID > 0 {
			h.db.Model(&models.RootDomain{}).Where("id = ?", domain.RootDomainID).
				UpdateColumn("registration_count", gorm.Expr("GREATEST(registration_count - 1, 0)"))
		}

		successCount++
	}

	fmt.Printf("Cleanup completed: %d domains deleted, %d failed.\n", successCount, failCount)
}

// updateDomainNSRecordsInPowerDNS 更新域名在 PowerDNS root zone 中的 NS 记录
// 当用户使用自定义 nameservers 时，需要在 root domain 的 zone 中添加该子域名的 NS 记录
// 当用户切换回默认 NS 时，删除这些自定义 NS 记录，并为该子域名创建独立的 zone
func (h *DomainHandler) updateDomainNSRecordsInPowerDNS(domain *models.Domain, nameservers []string, isDefault bool) {
	if domain.RootDomain == nil {
		fmt.Printf("Warning: Cannot update NS records for domain %s: root domain not loaded\n", domain.FullDomain)
		return
	}

	rootDomain := domain.RootDomain.Domain
	subdomainFQDN := domain.FullDomain

	// 如果使用默认 NS，删除子域名的 NS 记录，并创建子域名的独立 zone
	if isDefault {
		// 1. 删除在 root zone 中的 NS 记录
		if err := h.pdns.DeleteRRset(rootDomain, subdomainFQDN, "NS"); err != nil {
			fmt.Printf("Warning: Failed to delete NS records for %s in PowerDNS: %v\n", subdomainFQDN, err)
		} else {
			fmt.Printf("Deleted custom NS records for %s (using default NS)\n", subdomainFQDN)
		}

		// 2. 为子域名创建独立的 zone
		defaultNS := []string{h.cfg.DNS.DefaultNS1, h.cfg.DNS.DefaultNS2}
		if err := h.pdns.CreateZone(subdomainFQDN, defaultNS); err != nil {
			// 如果 zone 已经存在，不是错误
			if !strings.Contains(err.Error(), "Conflict") && !strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Warning: Failed to create zone for %s in PowerDNS: %v\n", subdomainFQDN, err)
			} else {
				fmt.Printf("Zone for %s already exists in PowerDNS\n", subdomainFQDN)
			}
		} else {
			fmt.Printf("Created zone for %s with default nameservers: %v\n", subdomainFQDN, defaultNS)
		}
		return
	}

	// 使用自定义 NS
	// 1. 删除子域名的独立 zone（如果存在）
	if err := h.pdns.DeleteZone(subdomainFQDN); err != nil {
		// zone 不存在不是错误
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "Could not find") {
			fmt.Printf("Warning: Failed to delete zone for %s in PowerDNS: %v\n", subdomainFQDN, err)
		}
	} else {
		fmt.Printf("Deleted zone for %s (using custom NS)\n", subdomainFQDN)
	}

	// 2. 在 root domain zone 中添加子域名的 NS 记录
	entries := make([]powerdns.RecordEntry, 0, len(nameservers))
	for _, ns := range nameservers {
		entries = append(entries, powerdns.RecordEntry{
			Content: ns,
		})
	}

	if err := h.pdns.SetRecords(rootDomain, subdomainFQDN, "NS", entries, 3600); err != nil {
		fmt.Printf("Warning: Failed to set NS records for %s in PowerDNS: %v\n", subdomainFQDN, err)
	} else {
		fmt.Printf("Updated NS records for %s with custom nameservers: %v\n", subdomainFQDN, nameservers)
	}
}
