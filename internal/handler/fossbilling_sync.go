package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/middleware"
	"opendomain/internal/models"
)

type FOSSBillingSyncHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewFOSSBillingSyncHandler(db *gorm.DB, cfg *config.Config) *FOSSBillingSyncHandler {
	return &FOSSBillingSyncHandler{
		db:  db,
		cfg: cfg,
	}
}

// FOSSBilling API 响应结构
type FOSSBillingResponse struct {
	Result interface{} `json:"result"`
	Error  *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error"`
}

type FOSSBillingDomain struct {
	ID           int       `json:"id"`
	Domain       string    `json:"domain"`
	RegisteredAt time.Time `json:"registered_at"`
	ExpiresAt    time.Time `json:"expires_at"`
	Status       string    `json:"status"`
}

// SyncFromFOSSBilling 从 FOSSBilling 同步域名
// @Summary 从 FOSSBilling 同步域名
// @Tags User
// @Accept json
// @Produce json
// @Param request body object true "同步请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/user/sync-from-fossbilling [post]
// @Security Bearer
func (h *FOSSBillingSyncHandler) SyncFromFOSSBilling(c *gin.Context) {
	// 检查 FOSSBilling 功能是否启用
	if !h.cfg.FOSSBilling.Enabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "FOSSBilling integration is not enabled"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 获取用户信息
	var user models.User
	if err := h.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 获取请求参数
	var req struct {
		FOSSBillingURL    string `json:"fossbilling_url"`    // FOSSBilling 服务器地址（可选，使用配置默认值）
		FOSSBillingAPIKey string `json:"fossbilling_api_key" binding:"required"` // FOSSBilling API Key（必填）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// 使用配置的默认 URL（如果未提供）
	fossBillingURL := req.FOSSBillingURL
	if fossBillingURL == "" {
		fossBillingURL = h.cfg.FOSSBilling.URL
	}
	if fossBillingURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FOSSBilling URL is required"})
		return
	}

	// 使用 API Key 作为 token
	token := req.FOSSBillingAPIKey

	// 2. 获取域名列表
	domains, err := h.fetchFOSSBillingDomains(fossBillingURL, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains: " + err.Error()})
		return
	}

	if len(domains) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":        "No domains found in FOSSBilling",
			"synced_count":   0,
			"skipped_count":  0,
			"existing_count": 0,
		})
		return
	}

	// 3. 同步域名到 OpenDomain
	syncResult := h.syncDomains(user.ID, domains)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Sync completed",
		"synced_count":   syncResult.Synced,
		"skipped_count":  syncResult.Skipped,
		"existing_count": syncResult.Existing,
		"error_count":    syncResult.Errors,
		"details":        syncResult.Details,
	})
}

// loginFOSSBilling 登录 FOSSBilling 获取 API Token
func (h *FOSSBillingSyncHandler) loginFOSSBilling(baseURL, email, password string) (string, error) {
	url := strings.TrimRight(baseURL, "/") + "/api/guest/client/login"

	payload := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result FOSSBillingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	// 提取 token
	resultMap, ok := result.Result.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}

	token, ok := resultMap["api_token"].(string)
	if !ok {
		return "", fmt.Errorf("no api_token in response")
	}

	return token, nil
}

// fetchFOSSBillingDomains 获取 FOSSBilling 域名列表
func (h *FOSSBillingSyncHandler) fetchFOSSBillingDomains(baseURL, apiKey string) ([]FOSSBillingDomain, error) {
	// FOSSBilling 使用 order/get_list 获取所有订单，然后过滤域名类型
	url := strings.TrimRight(baseURL, "/") + "/api/client/order/get_list"

	fmt.Printf("[DEBUG] Fetching orders from: %s\n", url)
	fmt.Printf("[DEBUG] Using HTTP Basic Auth with username: client\n")

	// 创建空的 JSON payload
	payload := map[string]interface{}{}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 HTTP Basic Authentication
	req.SetBasicAuth("client", apiKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	fmt.Printf("[DEBUG] Response status: %d\n", resp.StatusCode)
	fmt.Printf("[DEBUG] Response body: %s\n", string(body))

	var result FOSSBillingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("API error: %s", result.Error.Message)
	}

	// 解析订单列表并过滤域名类型
	var domains []FOSSBillingDomain
	resultMap, ok := result.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	list, ok := resultMap["list"].([]interface{})
	if !ok {
		return domains, nil // 空列表
	}

	for _, item := range list {
		orderMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// 只处理域名类型的订单
		serviceType, _ := orderMap["service_type"].(string)
		if serviceType != "domain" {
			continue
		}

		domain := FOSSBillingDomain{}

		// 从 config 中提取域名信息
		config, ok := orderMap["config"].(map[string]interface{})
		if !ok {
			continue
		}

		// 获取域名: register_sld + register_tld
		sld, _ := config["register_sld"].(string)
		tld, _ := config["register_tld"].(string)
		if sld != "" && tld != "" {
			domain.Domain = sld + tld
		}

		if domain.Domain == "" {
			continue
		}

		// 获取订单 ID
		if id, ok := orderMap["id"].(float64); ok {
			domain.ID = int(id)
		}

		// 获取激活时间作为注册时间
		if activatedAt, ok := orderMap["activated_at"].(string); ok && activatedAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", activatedAt); err == nil {
				domain.RegisteredAt = t
			}
		}

		// 获取过期时间
		if expiresAt, ok := orderMap["expires_at"].(string); ok && expiresAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", expiresAt); err == nil {
				domain.ExpiresAt = t
			}
		}

		// 获取状态
		if status, ok := orderMap["status"].(string); ok {
			domain.Status = status
		}

		domains = append(domains, domain)
	}

	return domains, nil
}

type SyncResult struct {
	Synced   int      `json:"synced"`
	Skipped  int      `json:"skipped"`
	Existing int      `json:"existing"`
	Errors   int      `json:"errors"`
	Details  []string `json:"details"`
}

// syncDomains 同步域名到 OpenDomain
func (h *FOSSBillingSyncHandler) syncDomains(userID uint, fossDomains []FOSSBillingDomain) SyncResult {
	result := SyncResult{
		Details: make([]string, 0),
	}

	for _, fossDomain := range fossDomains {
		// 解析域名
		subdomain, rootDomainStr := parseDomainString(fossDomain.Domain)
		if subdomain == "" || rootDomainStr == "" {
			result.Skipped++
			result.Details = append(result.Details, fmt.Sprintf("Skipped invalid domain: %s", fossDomain.Domain))
			continue
		}

		// 检查域名是否已存在
		var existingDomain models.Domain
		err := h.db.Where("full_domain = ? AND deleted_at IS NULL", fossDomain.Domain).First(&existingDomain).Error
		if err == nil {
			// 域名已存在，检查是否属于当前用户
			if existingDomain.UserID == userID {
				result.Existing++
				result.Details = append(result.Details, fmt.Sprintf("Domain already synced: %s", fossDomain.Domain))
			} else {
				// 域名已被其他用户注册
				result.Skipped++
				result.Details = append(result.Details, fmt.Sprintf("Domain already registered by another user: %s", fossDomain.Domain))
			}
			continue
		} else if err != gorm.ErrRecordNotFound {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("Error checking domain %s: %v", fossDomain.Domain, err))
			continue
		}

		// 检查是否在待激活列表中
		var pendingDomain models.PendingDomain
		err = h.db.Where("full_domain = ? AND deleted_at IS NULL", fossDomain.Domain).First(&pendingDomain).Error
		isPending := (err == nil)
		if err != nil && err != gorm.ErrRecordNotFound {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("Error checking pending domain %s: %v", fossDomain.Domain, err))
			continue
		}

		// 获取或创建根域名
		var rootDomain models.RootDomain
		err = h.db.Where("domain = ?", rootDomainStr).First(&rootDomain).Error
		if err == gorm.ErrRecordNotFound {
			// 创建根域名
			rootDomain = models.RootDomain{
				Domain:                rootDomainStr,
				IsActive:              true,
				IsFree:                true,
				Priority:              0,
				UseDefaultNameservers: true,
				Nameservers:           `["ns1.example.com","ns2.example.com"]`,
			}
			if err := h.db.Create(&rootDomain).Error; err != nil {
				result.Errors++
				result.Details = append(result.Details, fmt.Sprintf("Failed to create root domain %s: %v", rootDomainStr, err))
				continue
			}
			result.Details = append(result.Details, fmt.Sprintf("Created root domain: %s", rootDomainStr))
		} else if err != nil {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("Error querying root domain %s: %v", rootDomainStr, err))
			continue
		}

		// 设置日期
		registeredAt := time.Now()
		if !fossDomain.RegisteredAt.IsZero() {
			registeredAt = fossDomain.RegisteredAt
		}

		expiresAt := registeredAt.AddDate(1, 0, 0)
		if !fossDomain.ExpiresAt.IsZero() {
			expiresAt = fossDomain.ExpiresAt
		}

		// 映射状态
		status := "active"
		if fossDomain.Status == "suspended" {
			status = "suspended"
		}

		// 创建域名
		domain := models.Domain{
			UserID:                userID,
			RootDomainID:          rootDomain.ID,
			Subdomain:             subdomain,
			FullDomain:            fossDomain.Domain,
			Status:                status,
			RegisteredAt:          registeredAt,
			ExpiresAt:             expiresAt,
			AutoRenew:             false,
			UseDefaultNameservers: true,
			Nameservers:           `["ns1.example.com","ns2.example.com"]`,
			DNSSynced:             false,
		}

		if err := h.db.Create(&domain).Error; err != nil {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("Failed to create domain %s: %v", fossDomain.Domain, err))
			continue
		}

		// 如果是从待激活列表激活的，删除pending_domain记录
		if isPending {
			if err := h.db.Where("id = ?", pendingDomain.ID).Delete(&models.PendingDomain{}).Error; err != nil {
				// 只记录日志，不影响同步结果
				fmt.Printf("[WARNING] Failed to delete pending domain %s: %v\n", fossDomain.Domain, err)
			} else {
				result.Details = append(result.Details, fmt.Sprintf("Activated domain from pending list: %s", fossDomain.Domain))
			}
		}

		// 更新根域名注册计数
		h.db.Model(&models.RootDomain{}).Where("id = ?", rootDomain.ID).
			UpdateColumn("registration_count", gorm.Expr("registration_count + ?", 1))

		result.Synced++
		if !isPending {
			result.Details = append(result.Details, fmt.Sprintf("Synced domain: %s", fossDomain.Domain))
		}
	}

	return result
}

// parseDomainString 解析完整域名
func parseDomainString(fullDomain string) (subdomain, rootDomain string) {
	// 清理域名，移除空格和特殊字符
	fullDomain = strings.TrimSpace(fullDomain)
	fullDomain = strings.ToLower(fullDomain)

	parts := strings.Split(fullDomain, ".")
	if len(parts) < 2 {
		return "", ""
	}

	// 对于 test.com -> subdomain="test", rootDomain="com"
	// 对于 test.loc.cc -> subdomain="test", rootDomain="loc.cc"
	subdomain = parts[0]
	rootDomain = strings.Join(parts[1:], ".")
	return subdomain, rootDomain
}

// GetSyncStatus 获取同步状态
// @Summary 获取同步状态
// @Tags User
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/user/sync-status [get]
// @Security Bearer
func (h *FOSSBillingSyncHandler) GetSyncStatus(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 获取用户的域名数量
	var domainCount int64
	h.db.Model(&models.Domain{}).Where("user_id = ? AND deleted_at IS NULL", userID).Count(&domainCount)

	// 获取最近同步时间（最后一个域名的创建时间）
	var lastDomain models.Domain
	h.db.Where("user_id = ?", userID).Order("created_at DESC").First(&lastDomain)

	c.JSON(http.StatusOK, gin.H{
		"current_domain_count": domainCount,
		"last_sync_at":         lastDomain.CreatedAt,
		"can_sync":             true,
	})
}

// AdminSyncAllDomains 管理员同步所有FOSSBilling域名到待激活列表
// @Summary 管理员同步所有域名
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/sync-fossbilling-domains [post]
// @Security Bearer
func (h *FOSSBillingSyncHandler) AdminSyncAllDomains(c *gin.Context) {
	// 检查 FOSSBilling 功能是否启用
	if !h.cfg.FOSSBilling.Enabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "FOSSBilling integration is not enabled"})
		return
	}

	// 检查是否配置了Admin API Key
	if h.cfg.FOSSBilling.AdminAPIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin API Key is not configured"})
		return
	}

	// 使用Admin API获取所有域名
	domains, err := h.fetchFOSSBillingDomainsAsAdmin(h.cfg.FOSSBilling.URL, h.cfg.FOSSBilling.AdminAPIKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains: " + err.Error()})
		return
	}

	if len(domains) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message":      "No domains found in FOSSBilling",
			"synced_count": 0,
		})
		return
	}

	// 检查并去重（优先保留active状态的订单）
	uniqueDomains := make(map[string]FOSSBillingDomain)
	duplicateCount := 0
	statusCounts := make(map[string]int)
	replacedCount := 0

	for _, domain := range domains {
		statusCounts[domain.Status]++

		if existing, exists := uniqueDomains[domain.Domain]; exists {
			duplicateCount++
			// 优先保留 active 状态的订单
			if domain.Status == "active" && existing.Status != "active" {
				uniqueDomains[domain.Domain] = domain
				replacedCount++
			}
		} else {
			uniqueDomains[domain.Domain] = domain
		}
	}

	// 转换回slice
	deduplicatedDomains := make([]FOSSBillingDomain, 0, len(uniqueDomains))
	for _, domain := range uniqueDomains {
		deduplicatedDomains = append(deduplicatedDomains, domain)
	}

	fmt.Printf("[INFO] Fetched %d orders from FOSSBilling API\n", len(domains))
	fmt.Printf("[INFO] Order status breakdown: ")
	for status, count := range statusCounts {
		fmt.Printf("%s=%d ", status, count)
	}
	fmt.Printf("\n")
	fmt.Printf("[INFO] Deduplication: %d duplicates found, %d replaced with active status\n", duplicateCount, replacedCount)
	fmt.Printf("[INFO] Result: %d unique domains to sync\n", len(deduplicatedDomains))

	// 同步到pending_domains表
	syncResult := h.syncToPendingDomains(deduplicatedDomains)

	c.JSON(http.StatusOK, gin.H{
		"message":        "Admin sync completed",
		"synced_count":   syncResult.Synced,
		"skipped_count":  syncResult.Skipped,
		"existing_count": syncResult.Existing,
		"error_count":    syncResult.Errors,
		"details":        syncResult.Details,
	})
}

// fetchFOSSBillingDomainsAsAdmin 使用Admin API获取所有域名
// 注意: FOSSBilling的API分页有bug,page参数被忽略,所以我们一次性获取所有数据
func (h *FOSSBillingSyncHandler) fetchFOSSBillingDomainsAsAdmin(baseURL, adminAPIKey string) ([]FOSSBillingDomain, error) {
	fmt.Printf("[INFO] Fetching all domains from FOSSBilling...\n")
	fmt.Printf("[WARNING] FOSSBilling API pagination is broken, using per_page=10000 to fetch all data at once\n")

	// FOSSBilling API的分页不工作,所以设置一个很大的per_page值一次性获取所有数据
	domains, _, err := h.fetchFOSSBillingDomainsPage(baseURL, adminAPIKey, 1, 10000)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[INFO] Fetched %d total orders from FOSSBilling\n", len(domains))
	return domains, nil
}

// fetchFOSSBillingDomainsPage 获取单页域名数据
func (h *FOSSBillingSyncHandler) fetchFOSSBillingDomainsPage(baseURL, adminAPIKey string, page, perPage int) ([]FOSSBillingDomain, bool, error) {
	url := strings.TrimRight(baseURL, "/") + "/api/admin/order/get_list"

	// 创建带分页参数的 payload
	payload := map[string]interface{}{
		"page":     page,
		"per_page": perPage,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, false, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, false, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 HTTP Basic Authentication
	req.SetBasicAuth("admin", adminAPIKey)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, fmt.Errorf("failed to read response: %w", err)
	}

	var result FOSSBillingResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, false, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return nil, false, fmt.Errorf("API error: %s", result.Error.Message)
	}

	// 解析订单列表
	var domains []FOSSBillingDomain
	resultMap, ok := result.Result.(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("invalid response format")
	}

	list, ok := resultMap["list"].([]interface{})
	if !ok {
		return domains, false, nil
	}

	// 检查分页信息
	totalPages := 1
	if totalPagesFloat, ok := resultMap["pages"].(float64); ok {
		totalPages = int(totalPagesFloat)
	}
	hasMore := page < totalPages

	fmt.Printf("[DEBUG] Page %d response: total_pages=%d, items_in_list=%d, hasMore=%v\n", page, totalPages, len(list), hasMore)

	// 收集订单ID用于验证是否重复
	orderIDs := make([]int, 0)

	// 解析域名订单
	for _, item := range list {
		orderMap, ok := item.(map[string]interface{})
		if !ok {
			fmt.Printf("[WARNING] Skipping invalid order item\n")
			continue
		}

		// 只处理域名类型的订单
		serviceType, _ := orderMap["service_type"].(string)
		if serviceType != "domain" {
			continue
		}

		domain := FOSSBillingDomain{}

		// 从 config 中提取域名信息
		config, ok := orderMap["config"].(map[string]interface{})
		if !ok {
			orderID := orderMap["id"]
			fmt.Printf("[WARNING] Order %v: missing config\n", orderID)
			continue
		}

		// 获取域名: register_sld + register_tld
		sld, _ := config["register_sld"].(string)
		tld, _ := config["register_tld"].(string)
		if sld != "" && tld != "" {
			domain.Domain = sld + tld
		}

		if domain.Domain == "" {
			orderID := orderMap["id"]
			fmt.Printf("[WARNING] Order %v: empty domain (sld=%s, tld=%s)\n", orderID, sld, tld)
			continue
		}

		// 获取订单 ID
		if id, ok := orderMap["id"].(float64); ok {
			domain.ID = int(id)
			orderIDs = append(orderIDs, domain.ID)
		}

		// 获取激活时间作为注册时间
		if activatedAt, ok := orderMap["activated_at"].(string); ok && activatedAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", activatedAt); err == nil {
				domain.RegisteredAt = t
			}
		}

		// 获取过期时间
		if expiresAt, ok := orderMap["expires_at"].(string); ok && expiresAt != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", expiresAt); err == nil {
				domain.ExpiresAt = t
			}
		}

		// 获取状态
		if status, ok := orderMap["status"].(string); ok {
			domain.Status = status
		}

		domains = append(domains, domain)
	}

	// 打印订单ID范围用于验证分页
	if len(orderIDs) > 0 {
		fmt.Printf("[DEBUG] Page %d order IDs: first=%d, last=%d, sample=%v\n",
			page, orderIDs[0], orderIDs[len(orderIDs)-1], orderIDs[:min(5, len(orderIDs))])
	}

	return domains, hasMore, nil
}

// min helper function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// syncToPendingDomains 同步域名到待激活表
func (h *FOSSBillingSyncHandler) syncToPendingDomains(fossDomains []FOSSBillingDomain) SyncResult {
	result := SyncResult{
		Details: make([]string, 0),
	}

	fmt.Printf("[INFO] Starting to sync %d domains to pending_domains table\n", len(fossDomains))

	for i, fossDomain := range fossDomains {
		// 每处理100个域名输出一次进度
		if (i+1)%100 == 0 || i == len(fossDomains)-1 {
			fmt.Printf("[PROGRESS] Processed %d/%d domains (Synced: %d, Skipped: %d, Existing: %d, Errors: %d)\n",
				i+1, len(fossDomains), result.Synced, result.Skipped, result.Existing, result.Errors)
		}

		// 解析域名
		subdomain, rootDomainStr := parseDomainString(fossDomain.Domain)
		if subdomain == "" || rootDomainStr == "" {
			result.Skipped++
			errMsg := fmt.Sprintf("Invalid domain format: %s", fossDomain.Domain)
			result.Details = append(result.Details, errMsg)
			continue
		}

		// 检查是否已经在pending_domains中
		var existingPending models.PendingDomain
		err := h.db.Where("full_domain = ? AND deleted_at IS NULL", fossDomain.Domain).First(&existingPending).Error
		if err == nil {
			result.Existing++
			continue // 跳过，不记录Details
		} else if err != gorm.ErrRecordNotFound {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("DB error checking pending %s: %v", fossDomain.Domain, err))
			continue
		}

		// 检查是否已经被激活（在domains表中）
		var existingDomain models.Domain
		err = h.db.Where("full_domain = ? AND deleted_at IS NULL", fossDomain.Domain).First(&existingDomain).Error
		if err == nil {
			result.Skipped++
			continue // 跳过，不记录Details
		} else if err != gorm.ErrRecordNotFound {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("DB error checking domain %s: %v", fossDomain.Domain, err))
			continue
		}

		// 获取或创建根域名
		var rootDomain models.RootDomain
		err = h.db.Where("domain = ?", rootDomainStr).First(&rootDomain).Error
		if err == gorm.ErrRecordNotFound {
			// 创建根域名
			rootDomain = models.RootDomain{
				Domain:                rootDomainStr,
				IsActive:              true,
				IsFree:                true,
				Priority:              0,
				UseDefaultNameservers: true,
				Nameservers:           `["ns1.example.com","ns2.example.com"]`,
			}
			if err := h.db.Create(&rootDomain).Error; err != nil {
				result.Errors++
				result.Details = append(result.Details, fmt.Sprintf("Failed to create root domain %s: %v", rootDomainStr, err))
				continue
			}
			// 根域名创建成功
		} else if err != nil {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("DB error querying root domain %s: %v", rootDomainStr, err))
			continue
		}

		// 设置日期
		registeredAt := time.Now()
		if !fossDomain.RegisteredAt.IsZero() {
			registeredAt = fossDomain.RegisteredAt
		}

		expiresAt := registeredAt.AddDate(1, 0, 0)
		if !fossDomain.ExpiresAt.IsZero() {
			expiresAt = fossDomain.ExpiresAt
		}

		// 创建待激活域名
		pendingDomain := models.PendingDomain{
			RootDomainID:       rootDomain.ID,
			Subdomain:          subdomain,
			FullDomain:         fossDomain.Domain,
			FOSSBillingOrderID: fossDomain.ID,
			Status:             "pending",
			RegisteredAt:       registeredAt,
			ExpiresAt:          expiresAt,
		}

		if err := h.db.Create(&pendingDomain).Error; err != nil {
			result.Errors++
			result.Details = append(result.Details, fmt.Sprintf("Failed to create pending domain %s: %v", fossDomain.Domain, err))
			continue
		}

		result.Synced++
		// 成功创建，不输出详细日志
	}

	// 打印汇总
	fmt.Printf("\n[SUMMARY] Sync completed:\n")
	fmt.Printf("  - Total processed: %d\n", len(fossDomains))
	fmt.Printf("  - Synced: %d\n", result.Synced)
	fmt.Printf("  - Skipped (already activated): %d\n", result.Skipped)
	fmt.Printf("  - Existing (already in pending): %d\n", result.Existing)
	fmt.Printf("  - Errors: %d\n", result.Errors)

	// 打印前20个错误详情
	if len(result.Details) > 0 {
		fmt.Printf("\n[ERROR/WARNING DETAILS] (showing first 20):\n")
		count := 0
		for _, detail := range result.Details {
			fmt.Printf("  - %s\n", detail)
			count++
			if count >= 20 {
				if len(result.Details) > 20 {
					fmt.Printf("  ... and %d more issues\n", len(result.Details)-20)
				}
				break
			}
		}
	}

	return result
}

// ListPendingDomains 获取待激活域名列表
// @Summary 获取待激活域名列表
// @Tags Admin
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/pending-domains [get]
// @Security Bearer
func (h *FOSSBillingSyncHandler) ListPendingDomains(c *gin.Context) {
	// 分页参数
	page := 1
	perPage := 50

	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if pp := c.Query("per_page"); pp != "" {
		if perPageNum, err := strconv.Atoi(pp); err == nil && perPageNum > 0 && perPageNum <= 100 {
			perPage = perPageNum
		}
	}

	// 计算偏移量
	offset := (page - 1) * perPage

	// 获取总数
	var total int64
	if err := h.db.Model(&models.PendingDomain{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count pending domains"})
		return
	}

	// 获取分页数据
	var pendingDomains []models.PendingDomain
	if err := h.db.Preload("RootDomain").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&pendingDomains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending domains"})
		return
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	c.JSON(http.StatusOK, gin.H{
		"pending_domains": pendingDomains,
		"pagination": gin.H{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// DeletePendingDomain 删除待激活域名
// @Summary 删除待激活域名
// @Tags Admin
// @Param id path int true "Pending Domain ID"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/pending-domains/{id} [delete]
// @Security Bearer
func (h *FOSSBillingSyncHandler) DeletePendingDomain(c *gin.Context) {
	id := c.Param("id")

	var pendingDomain models.PendingDomain
	if err := h.db.First(&pendingDomain, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pending domain not found"})
		return
	}

	if err := h.db.Delete(&pendingDomain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete pending domain"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pending domain deleted successfully",
	})
}

// GetPublicPendingDomains 获取公开的待激活域名列表(分页)
// @Summary 获取公开的待激活域名列表
// @Tags Public
// @Produce json
// @Param page query int false "页码" default(1)
// @Param per_page query int false "每页数量" default(50)
// @Success 200 {object} map[string]interface{}
// @Router /api/public/pending-domains [get]
func (h *FOSSBillingSyncHandler) GetPublicPendingDomains(c *gin.Context) {
	// 获取分页参数
	page := 1
	perPage := 50

	if p := c.Query("page"); p != "" {
		if pageNum, err := strconv.Atoi(p); err == nil && pageNum > 0 {
			page = pageNum
		}
	}

	if pp := c.Query("per_page"); pp != "" {
		if perPageNum, err := strconv.Atoi(pp); err == nil && perPageNum > 0 && perPageNum <= 100 {
			perPage = perPageNum
		}
	}

	// 计算偏移量
	offset := (page - 1) * perPage

	// 获取总数
	var total int64
	if err := h.db.Model(&models.PendingDomain{}).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count pending domains"})
		return
	}

	// 获取分页数据
	var pendingDomains []models.PendingDomain
	if err := h.db.Preload("RootDomain").
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Limit(perPage).
		Offset(offset).
		Find(&pendingDomains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pending domains"})
		return
	}

	// 计算总页数
	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"pending_domains": pendingDomains,
		"pagination": gin.H{
			"page":        page,
			"per_page":    perPage,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}
