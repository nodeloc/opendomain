package scanner

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"

	"opendomain/internal/config"
	"opendomain/internal/models"
	"opendomain/internal/services"
)

type Scanner struct {
	db  *gorm.DB
	cfg *config.Config

	// VirusTotal 速率限制: 4 lookups/min, 500/day
	vtMu             sync.Mutex
	vtLastCall       time.Time
	vtDailyCount     int
	vtDailyResetDate string // 格式: YYYY-MM-DD

	// Google Safe Browsing 速率限制: 1 request/sec, 10000/day
	gsbMu             sync.Mutex
	gsbLastCall       time.Time
	gsbDailyCount     int
	gsbDailyResetDate string // 格式: YYYY-MM-DD
}

func NewScanner(db *gorm.DB, cfg *config.Config) *Scanner {
	s := &Scanner{db: db, cfg: cfg}

	// 从数据库加载今天的配额使用情况
	s.loadQuotaFromDB()

	return s
}

// loadQuotaFromDB 从数据库加载今天的配额使用情况
func (s *Scanner) loadQuotaFromDB() {
	today := time.Now().Format("2006-01-02")

	// 加载 Google Safe Browsing 配额
	var gsbQuota models.APIQuota
	if err := s.db.Where("api_name = ?", "google_safe_browsing").First(&gsbQuota).Error; err == nil {
		if gsbQuota.Date == today {
			s.gsbDailyCount = gsbQuota.UsedCount
			s.gsbDailyResetDate = today
			fmt.Printf("[INFO] Loaded Google Safe Browsing quota: %d/10000 for %s\n", s.gsbDailyCount, today)
		}
	} else if err != gorm.ErrRecordNotFound {
		// 如果是表不存在等其他错误，记录警告但不中断启动
		fmt.Printf("[WARNING] Failed to load Google Safe Browsing quota: %v\n", err)
	}

	// 加载 VirusTotal 配额
	var vtQuota models.APIQuota
	if err := s.db.Where("api_name = ?", "virustotal").First(&vtQuota).Error; err == nil {
		if vtQuota.Date == today {
			s.vtDailyCount = vtQuota.UsedCount
			s.vtDailyResetDate = today
			fmt.Printf("[INFO] Loaded VirusTotal quota: %d/500 for %s\n", s.vtDailyCount, today)
		}
	} else if err != gorm.ErrRecordNotFound {
		// 如果是表不存在等其他错误，记录警告但不中断启动
		fmt.Printf("[WARNING] Failed to load VirusTotal quota: %v\n", err)
	}
}

// saveQuotaToDB 保存配额使用情况到数据库
func (s *Scanner) saveQuotaToDB(apiName string, count int, limit int) {
	today := time.Now().Format("2006-01-02")

	quota := models.APIQuota{
		APIName:    apiName,
		Date:       today,
		UsedCount:  count,
		DailyLimit: limit,
	}

	// 使用 upsert 更新或插入
	if err := s.db.Where("api_name = ?", apiName).Assign(quota).FirstOrCreate(&quota).Error; err != nil {
		// 如果表不存在或其他错误，记录警告但不中断运行
		fmt.Printf("[WARNING] Failed to save quota for %s: %v\n", apiName, err)
	}
}

// GetQuotaStatus 获取当前配额使用状态
func (s *Scanner) GetQuotaStatus() map[string]interface{} {
	return map[string]interface{}{
		"google_safe_browsing": map[string]interface{}{
			"used":  s.gsbDailyCount,
			"limit": 10000,
			"date":  s.gsbDailyResetDate,
		},
		"virustotal": map[string]interface{}{
			"used":  s.vtDailyCount,
			"limit": 500,
			"date":  s.vtDailyResetDate,
		},
	}
}

// vtRateLimit 等待至少 15 秒间隔以遵守 VirusTotal 免费 API 速率限制 (4次/分钟)
func (s *Scanner) vtRateLimit() {
	s.vtMu.Lock()
	defer s.vtMu.Unlock()

	if !s.vtLastCall.IsZero() {
		elapsed := time.Since(s.vtLastCall)
		if wait := 15*time.Second - elapsed; wait > 0 {
			time.Sleep(wait)
		}
	}
	s.vtLastCall = time.Now()
}

// vtCheckDailyQuota 检查并更新每日配额
// 返回 true 表示可以继续调用，false 表示已超出配额
func (s *Scanner) vtCheckDailyQuota() bool {
	s.vtMu.Lock()
	defer s.vtMu.Unlock()

	today := time.Now().Format("2006-01-02")

	// 新的一天，重置计数器
	if s.vtDailyResetDate != today {
		s.vtDailyResetDate = today
		s.vtDailyCount = 0
		fmt.Printf("[INFO] VirusTotal daily quota reset for %s\n", today)
	}

	// 检查是否超出每日限制
	if s.vtDailyCount >= 500 {
		return false
	}

	// 增加计数
	s.vtDailyCount++

	// 保存到数据库
	go s.saveQuotaToDB("virustotal", s.vtDailyCount, 500)

	return true
}

// gsbRateLimit 等待至少 1 秒间隔以遵守 Google Safe Browsing API 速率限制 (1次/秒)
func (s *Scanner) gsbRateLimit() {
	s.gsbMu.Lock()
	defer s.gsbMu.Unlock()

	if !s.gsbLastCall.IsZero() {
		elapsed := time.Since(s.gsbLastCall)
		if wait := time.Second - elapsed; wait > 0 {
			time.Sleep(wait)
		}
	}
	s.gsbLastCall = time.Now()
}

// gsbCheckDailyQuota 检查并更新 Google Safe Browsing 每日配额
// 返回 true 表示可以继续调用，false 表示已超出配额
func (s *Scanner) gsbCheckDailyQuota() bool {
	s.gsbMu.Lock()
	defer s.gsbMu.Unlock()

	today := time.Now().Format("2006-01-02")

	// 新的一天，重置计数器
	if s.gsbDailyResetDate != today {
		s.gsbDailyResetDate = today
		s.gsbDailyCount = 0
		fmt.Printf("[INFO] Google Safe Browsing daily quota reset for %s\n", today)
	}

	// 检查是否超出每日限制 (10000/day)
	if s.gsbDailyCount >= 10000 {
		return false
	}

	// 增加计数
	s.gsbDailyCount++

	// 保存到数据库
	go s.saveQuotaToDB("google_safe_browsing", s.gsbDailyCount, 10000)

	return true
}

// strPtr 返回字符串指针
func strPtr(s string) *string {
	return &s
}

// ScanAllDomains 扫描所有活跃域名（分批处理，遵守 API 速率限制）
func (s *Scanner) ScanAllDomains(ctx context.Context) error {
	var domains []models.Domain
	if err := s.db.Where("status = ?", "active").Find(&domains).Error; err != nil {
		return fmt.Errorf("failed to fetch domains: %w", err)
	}

	totalDomains := len(domains)
	if totalDomains == 0 {
		fmt.Println("[INFO] No active domains to scan")
		return nil
	}

	fmt.Printf("[INFO] Starting scan for %d domains\n", totalDomains)
	fmt.Printf("[INFO] Google Safe Browsing quota used: %d/10000\n", s.gsbDailyCount)
	fmt.Printf("[INFO] VirusTotal quota used: %d/500\n", s.vtDailyCount)

	// 分批处理，每批处理完成后等待一段时间
	batchSize := 50 // 每批50个域名
	for i := 0; i < totalDomains; i += batchSize {
		end := i + batchSize
		if end > totalDomains {
			end = totalDomains
		}

		batch := domains[i:end]
		fmt.Printf("[INFO] Processing batch %d-%d of %d domains\n", i+1, end, totalDomains)

		for _, domain := range batch {
			if err := s.ScanDomain(ctx, &domain); err != nil {
				fmt.Printf("[ERROR] Failed to scan domain %s: %v\n", domain.FullDomain, err)
			}
		}

		// 批次之间等待，避免过快
		if end < totalDomains {
			fmt.Printf("[INFO] Batch complete. Waiting 5 seconds before next batch...\n")
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Printf("[INFO] Scan complete. GSB used: %d/10000, VT used: %d/500\n",
		s.gsbDailyCount, s.vtDailyCount)
	return nil
}

// ScanAllPendingDomains 扫描所有待激活域名（分批处理，遵守 API 速率限制）
func (s *Scanner) ScanAllPendingDomains(ctx context.Context) error {
	var pendingDomains []models.PendingDomain
	if err := s.db.Where("deleted_at IS NULL").Find(&pendingDomains).Error; err != nil {
		return fmt.Errorf("failed to fetch pending domains: %w", err)
	}

	totalDomains := len(pendingDomains)
	if totalDomains == 0 {
		fmt.Println("[INFO] No pending domains to scan")
		return nil
	}

	fmt.Printf("[INFO] Starting scan for %d pending domains\n", totalDomains)

	// 分批处理
	batchSize := 50
	for i := 0; i < totalDomains; i += batchSize {
		end := i + batchSize
		if end > totalDomains {
			end = totalDomains
		}

		batch := pendingDomains[i:end]
		fmt.Printf("[INFO] Processing batch %d-%d of %d pending domains\n", i+1, end, totalDomains)

		for _, pending := range batch {
			if err := s.ScanPendingDomain(ctx, &pending); err != nil {
				fmt.Printf("[ERROR] Failed to scan pending domain %s: %v\n", pending.FullDomain, err)
			}
		}

		// 批次之间等待
		if end < totalDomains {
			fmt.Printf("[INFO] Batch complete. Waiting 5 seconds before next batch...\n")
			time.Sleep(5 * time.Second)
		}
	}

	fmt.Printf("[INFO] Pending domains scan complete\n")
	return nil
}

// ScanPendingDomain 扫描单个待激活域名
func (s *Scanner) ScanPendingDomain(ctx context.Context, pending *models.PendingDomain) error {
	// DNS 检查
	dnsStatus := s.checkDNS(pending.FullDomain)

	// HTTP 检查
	httpScan := s.checkHTTP(pending.FullDomain)
	httpStatus := httpScan.Status

	// 判断域名是否不可访问
	isDown := (dnsStatus == "failed" || httpStatus == "offline")

	now := time.Now()
	telegram := services.NewTelegramService(s.cfg)

	if isDown {
		if pending.FirstFailedAt == nil {
			// 记录首次失败时间
			pending.FirstFailedAt = &now
			pending.Status = "unhealthy"
			s.db.Save(pending)

			// 发送通知
			telegram.SendHealthAlert(pending.FullDomain,
				[]string{fmt.Sprintf("DNS: %s", dnsStatus), fmt.Sprintf("HTTP: %s", httpStatus)},
				"Pending domain unhealthy - will be removed in 30 days if not resolved")
			fmt.Printf("[INFO] Pending domain %s marked as unhealthy\n", pending.FullDomain)
		} else {
			daysSinceFailure := int(now.Sub(*pending.FirstFailedAt).Hours() / 24)

			if daysSinceFailure >= 30 {
				// 删除待激活域名
				s.db.Delete(pending)
				telegram.SendHealthAlert(pending.FullDomain,
					[]string{"Pending domain down for 30+ days"},
					"Pending domain REMOVED - now available for registration")
				fmt.Printf("[INFO] Pending domain %s removed after 30 days of failure\n", pending.FullDomain)
			} else if daysSinceFailure >= 7 {
				// 7天后发送警告
				pending.Status = "unhealthy"
				s.db.Save(pending)
				telegram.SendDeletionWarning(pending.FullDomain, 30-daysSinceFailure)
				fmt.Printf("[WARNING] Pending domain %s unhealthy for %d days, will be removed in %d days\n",
					pending.FullDomain, daysSinceFailure, 30-daysSinceFailure)
			}
		}
	} else {
		// 域名恢复健康
		if pending.FirstFailedAt != nil {
			pending.FirstFailedAt = nil
			pending.Status = "pending"
			s.db.Save(pending)
			fmt.Printf("[INFO] Pending domain %s recovered\n", pending.FullDomain)
		}
	}

	return nil
}

// ScanDomain 扫描单个域名
func (s *Scanner) ScanDomain(ctx context.Context, domain *models.Domain) error {
	// DNS 检查
	dnsStatus := s.checkDNS(domain.FullDomain)
	dnsScan := &models.DomainScan{
		DomainID:  domain.ID,
		ScanType:  "dns",
		Status:    dnsStatus,
		ScannedAt: time.Now(),
	}
	s.db.Create(dnsScan)

	// HTTP 检查
	httpScan := s.checkHTTP(domain.FullDomain)
	httpScan.DomainID = domain.ID
	httpScan.ScannedAt = time.Now()
	s.db.Create(httpScan)

	// SSL 检查
	sslScan := s.checkSSL(domain.FullDomain)
	sslScan.DomainID = domain.ID
	sslScan.ScannedAt = time.Now()
	s.db.Create(sslScan)

	// Google Safe Browsing 检查
	var safeBrowsingStatus string
	if sbScan := s.checkGoogleSafeBrowsing(domain.FullDomain); sbScan != nil {
		sbScan.DomainID = domain.ID
		sbScan.ScannedAt = time.Now()
		s.db.Create(sbScan)
		if sbScan.Status == "success" {
			safeBrowsingStatus = "safe"
		} else if sbScan.Status == "threat_detected" {
			safeBrowsingStatus = "unsafe"
		} else {
			// API 调用失败，不应标记为 unsafe
			safeBrowsingStatus = "unknown"
		}
	}

	// VirusTotal 检查
	var virusTotalStatus string
	if vtScan := s.checkVirusTotal(domain.FullDomain); vtScan != nil {
		vtScan.DomainID = domain.ID
		vtScan.ScannedAt = time.Now()
		s.db.Create(vtScan)
		if vtScan.Status == "success" {
			virusTotalStatus = "clean"
		} else if vtScan.Status == "not_found" {
			virusTotalStatus = "clean"
		} else if vtScan.Status == "failed" || vtScan.Status == "quota_exceeded" {
			// API 调用失败（如超出额度或每日配额），不应标记为恶意
			virusTotalStatus = "unknown"
		} else {
			// 其他状态，检查扫描详情以区分恶意和可疑
			virusTotalStatus = "malicious"
			if vtScan.ScanDetails != nil {
				var details map[string]interface{}
				if json.Unmarshal([]byte(*vtScan.ScanDetails), &details) == nil {
					if mal, ok := details["malicious"].(float64); ok && mal == 0 {
						virusTotalStatus = "suspicious"
					}
				}
			}
		}
	}

	// 更新摘要
	s.updateSummary(domain.ID, dnsStatus, httpScan.Status, sslScan.Status, safeBrowsingStatus, virusTotalStatus)

	return nil
}

// checkDNS 检查 DNS 解析
func (s *Scanner) checkDNS(domain string) string {
	_, err := net.LookupHost(domain)
	if err != nil {
		return "failed"
	}
	return "success"
}

// checkHTTP 检查 HTTP 状态
func (s *Scanner) checkHTTP(domain string) *models.DomainScan {
	scan := &models.DomainScan{
		ScanType: "http",
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	start := time.Now()
	resp, err := client.Get("http://" + domain)
	duration := int(time.Since(start).Milliseconds())
	scan.ResponseTime = &duration

	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = err.Error()
		return scan
	}
	defer resp.Body.Close()

	scan.Status = "success"
	scan.HTTPStatusCode = &resp.StatusCode

	return scan
}

// checkSSL 检查 SSL 证书
func (s *Scanner) checkSSL(domain string) *models.DomainScan {
	scan := &models.DomainScan{
		ScanType: "ssl",
	}

	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})

	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = err.Error()
		sslValid := false
		scan.SSLValid = &sslValid
		return scan
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		cert := certs[0]
		scan.Status = "success"
		scan.SSLExpiryDate = &cert.NotAfter

		// 检查证书是否有效
		now := time.Now()
		sslValid := now.After(cert.NotBefore) && now.Before(cert.NotAfter)
		scan.SSLValid = &sslValid
	} else {
		scan.Status = "failed"
		sslValid := false
		scan.SSLValid = &sslValid
	}

	return scan
}

// checkGoogleSafeBrowsing 通过 Google Safe Browsing API 检查域名安全性
func (s *Scanner) checkGoogleSafeBrowsing(domain string) *models.DomainScan {
	apiKey := s.cfg.Scanner.GoogleSafeBrowsingKey
	if apiKey == "" {
		return nil
	}

	scan := &models.DomainScan{
		ScanType: "safebrowsing",
	}

	// 检查每日配额
	if !s.gsbCheckDailyQuota() {
		scan.Status = "quota_exceeded"
		scan.ErrorMessage = fmt.Sprintf("Daily quota exceeded (10000/day). Count: %d", s.gsbDailyCount)
		fmt.Printf("[WARNING] Google Safe Browsing daily quota exceeded for domain %s (used: %d/10000)\n",
			domain, s.gsbDailyCount)
		return scan
	}

	// 速率限制：每秒最多1次请求
	s.gsbRateLimit()

	// 构建请求
	reqBody := map[string]interface{}{
		"client": map[string]string{
			"clientId":      "opendomain",
			"clientVersion": "1.0.0",
		},
		"threatInfo": map[string]interface{}{
			"threatTypes":      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"},
			"platformTypes":    []string{"ANY_PLATFORM"},
			"threatEntryTypes": []string{"URL"},
			"threatEntries": []map[string]string{
				{"url": "http://" + domain + "/"},
				{"url": "https://" + domain + "/"},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "failed to marshal request: " + err.Error()
		return scan
	}

	client := &http.Client{Timeout: 15 * time.Second}
	start := time.Now()
	resp, err := client.Post(
		"https://safebrowsing.googleapis.com/v4/threatMatches:find?key="+apiKey,
		"application/json",
		bytes.NewReader(jsonBody),
	)
	duration := int(time.Since(start).Milliseconds())
	scan.ResponseTime = &duration

	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "API request failed: " + err.Error()
		return scan
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "failed to read response: " + err.Error()
		return scan
	}

	if resp.StatusCode != http.StatusOK {
		scan.Status = "failed"
		scan.ErrorMessage = fmt.Sprintf("API returned status %d", resp.StatusCode)
		return scan
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "failed to parse response: " + err.Error()
		return scan
	}

	// 检查是否有威胁匹配
	if matches, ok := result["matches"]; ok && matches != nil {
		scan.Status = "threat_detected"
		details, _ := json.Marshal(map[string]interface{}{
			"safe":    false,
			"threats": matches,
		})
		scan.ScanDetails = strPtr(string(details))
		fmt.Printf("[SECURITY] Google Safe Browsing detected threat for domain %s: %v\n", domain, matches)
	} else {
		scan.Status = "success"
		details, _ := json.Marshal(map[string]interface{}{
			"safe": true,
		})
		scan.ScanDetails = strPtr(string(details))
	}

	return scan
}

// checkVirusTotal 通过 VirusTotal API 检查域名安全性
func (s *Scanner) checkVirusTotal(domain string) *models.DomainScan {
	apiKey := s.cfg.Scanner.VirusTotalKey
	if apiKey == "" {
		return nil
	}

	scan := &models.DomainScan{
		ScanType: "virustotal",
	}

	// 检查每日配额
	if !s.vtCheckDailyQuota() {
		scan.Status = "quota_exceeded"
		scan.ErrorMessage = fmt.Sprintf("Daily quota exceeded (500/day). Count: %d", s.vtDailyCount)
		fmt.Printf("[WARNING] VirusTotal daily quota exceeded for domain %s (used: %d/500)\n",
			domain, s.vtDailyCount)
		return scan
	}

	// 遵守 VirusTotal 速率限制
	s.vtRateLimit()

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", "https://www.virustotal.com/api/v3/domains/"+domain, nil)
	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "failed to create request: " + err.Error()
		return scan
	}
	req.Header.Set("x-apikey", apiKey)

	start := time.Now()
	resp, err := client.Do(req)
	duration := int(time.Since(start).Milliseconds())
	scan.ResponseTime = &duration

	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "API request failed: " + err.Error()
		return scan
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "failed to read response: " + err.Error()
		return scan
	}

	// 404 表示域名未被 VirusTotal 收录，视为安全
	if resp.StatusCode == http.StatusNotFound {
		scan.Status = "not_found"
		scan.ScanDetails = strPtr(`{"not_found":true}`)
		return scan
	}

	if resp.StatusCode != http.StatusOK {
		scan.Status = "failed"
		if resp.StatusCode == 429 {
			scan.ErrorMessage = "Rate limit exceeded (quota exhausted)"
		} else {
			scan.ErrorMessage = fmt.Sprintf("API returned status %d", resp.StatusCode)
		}
		return scan
	}

	var result struct {
		Data struct {
			Attributes struct {
				LastAnalysisStats struct {
					Malicious  int `json:"malicious"`
					Suspicious int `json:"suspicious"`
					Harmless   int `json:"harmless"`
					Undetected int `json:"undetected"`
				} `json:"last_analysis_stats"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		scan.Status = "failed"
		scan.ErrorMessage = "failed to parse response: " + err.Error()
		return scan
	}

	stats := result.Data.Attributes.LastAnalysisStats
	details, _ := json.Marshal(map[string]interface{}{
		"malicious":  stats.Malicious,
		"suspicious": stats.Suspicious,
		"harmless":   stats.Harmless,
		"undetected": stats.Undetected,
	})
	scan.ScanDetails = strPtr(string(details))

	if stats.Malicious > 0 {
		scan.Status = "failed"
		scan.ErrorMessage = fmt.Sprintf("%d engines detected as malicious", stats.Malicious)
	} else if stats.Suspicious > 0 {
		scan.Status = "failed"
		scan.ErrorMessage = fmt.Sprintf("%d engines detected as suspicious", stats.Suspicious)
	} else {
		scan.Status = "success"
	}

	return scan
}

// updateSummary 更新域名健康摘要
func (s *Scanner) updateSummary(domainID uint, dnsStatus, httpStatus, sslStatus, safeBrowsingStatus, virusTotalStatus string) {
	var summary models.DomainScanSummary
	result := s.db.Where("domain_id = ?", domainID).First(&summary)

	if result.Error == gorm.ErrRecordNotFound {
		summary = models.DomainScanSummary{
			DomainID: domainID,
		}
	}

	now := time.Now()
	summary.LastScannedAt = &now
	summary.TotalScans++

	// 更新状态
	if dnsStatus == "success" {
		summary.DNSStatus = "resolved"
	} else {
		summary.DNSStatus = "failed"
	}

	if httpStatus == "success" {
		summary.HTTPStatus = "online"
		summary.SuccessfulScans++
	} else {
		summary.HTTPStatus = "offline"
	}

	if sslStatus == "success" {
		summary.SSLStatus = "valid"
	} else if sslStatus == "failed" {
		summary.SSLStatus = "invalid"
	} else {
		summary.SSLStatus = "none"
	}

	// 更新安全扫描状态
	if safeBrowsingStatus != "" {
		summary.SafeBrowsingStatus = safeBrowsingStatus
	} else {
		summary.SafeBrowsingStatus = "unknown"
	}

	if virusTotalStatus != "" {
		summary.VirusTotalStatus = virusTotalStatus
	} else {
		summary.VirusTotalStatus = "unknown"
	}

	// 计算整体健康状态
	if safeBrowsingStatus == "unsafe" || virusTotalStatus == "malicious" {
		summary.OverallHealth = "degraded"
	} else if summary.DNSStatus == "resolved" && summary.HTTPStatus == "online" {
		summary.OverallHealth = "healthy"
	} else if summary.DNSStatus == "resolved" || summary.HTTPStatus == "online" {
		summary.OverallHealth = "degraded"
	} else {
		summary.OverallHealth = "down"
	}

	summary.UpdatedAt = now

	if result.Error == gorm.ErrRecordNotFound {
		s.db.Create(&summary)
	} else {
		s.db.Save(&summary)
	}

	// 处理自动 suspend/delete 逻辑
	s.handleAutoActions(domainID, safeBrowsingStatus, virusTotalStatus, summary)
}

// handleAutoActions 处理自动 suspend/delete 逻辑
func (s *Scanner) handleAutoActions(domainID uint, safeBrowsingStatus, virusTotalStatus string, summary models.DomainScanSummary) {
	var domain models.Domain
	if err := s.db.First(&domain, domainID).Error; err != nil {
		return
	}

	// 导入 Telegram 服务
	telegram := services.NewTelegramService(s.cfg)
	now := time.Now()

	// 1. 检测恶意内容 - 立即 suspend
	// 但需要确认是真正的威胁检测，而不是 API 错误
	if safeBrowsingStatus == "unsafe" || virusTotalStatus == "malicious" {
		// 获取最新的扫描记录以确认
		var latestScans []models.DomainScan
		s.db.Where("domain_id = ? AND scan_type IN (?, ?)", domainID, "safebrowsing", "virustotal").
			Order("scanned_at DESC").
			Limit(2).
			Find(&latestScans)

		// 检查是否真的是威胁检测，而不是 API 错误
		realThreat := false
		scanDetails := ""

		for _, scan := range latestScans {
			if scan.ScanType == "safebrowsing" && safeBrowsingStatus == "unsafe" {
				if scan.Status == "threat_detected" {
					realThreat = true
					scanDetails += fmt.Sprintf("SafeBrowsing: threat_detected")
					if scan.ScanDetails != nil {
						scanDetails += fmt.Sprintf(" %s", *scan.ScanDetails)
					}
				} else {
					scanDetails += fmt.Sprintf("SafeBrowsing: %s", scan.Status)
					if scan.ErrorMessage != "" {
						scanDetails += fmt.Sprintf(" (error: %s)", scan.ErrorMessage)
					}
				}
			}
			if scan.ScanType == "virustotal" && virusTotalStatus == "malicious" {
				if scan.Status == "threat_detected" {
					realThreat = true
					if scanDetails != "" {
						scanDetails += ", "
					}
					scanDetails += fmt.Sprintf("VirusTotal: threat_detected")
					if scan.ScanDetails != nil {
						scanDetails += fmt.Sprintf(" %s", *scan.ScanDetails)
					}
				} else {
					if scanDetails != "" {
						scanDetails += ", "
					}
					scanDetails += fmt.Sprintf("VirusTotal: %s", scan.Status)
					if scan.ErrorMessage != "" {
						scanDetails += fmt.Sprintf(" (error: %s)", scan.ErrorMessage)
					}
				}
			}
		}

		// 只有在确认是真正的威胁时才挂起
		if !realThreat {
			fmt.Printf("[WARNING] Domain %s marked as unsafe/malicious but latest scan shows API error, not suspending\n",
				domain.FullDomain)
			fmt.Printf("[DEBUG] Scan details: %s\n", scanDetails)
			return
		}

		if domain.Status != "suspended" {
			domain.Status = "suspended"
			s.db.Save(&domain)

			reason := fmt.Sprintf("Malicious content detected (Safe Browsing: %s, VirusTotal: %s)\nScan details: %s",
				safeBrowsingStatus, virusTotalStatus, scanDetails)
			telegram.SendAutoSuspendNotification(domain.FullDomain, reason)
		}
		return
	}

	// 2. 检测不可访问 - 7天后 suspend, 30天后删除
	isDown := (summary.DNSStatus == "failed" || summary.HTTPStatus == "offline")

	if isDown {
		if domain.FirstFailedAt == nil {
			// 记录首次失败时间
			domain.FirstFailedAt = &now
			s.db.Save(&domain)

			issues := []string{}
			if summary.DNSStatus == "failed" {
				issues = append(issues, fmt.Sprintf("DNS: %s", summary.DNSStatus))
			}
			if summary.HTTPStatus == "offline" {
				issues = append(issues, fmt.Sprintf("HTTP: %s", summary.HTTPStatus))
			}
			telegram.SendHealthAlert(domain.FullDomain, issues,
				"Monitoring started - will suspend in 7 days if not resolved")
		} else {
			daysSinceFailure := int(now.Sub(*domain.FirstFailedAt).Hours() / 24)

			if daysSinceFailure >= 30 {
				// 删除域名
				s.db.Delete(&domain)
				telegram.SendHealthAlert(domain.FullDomain,
					[]string{"Domain down for 30+ days"},
					"Domain DELETED")
			} else if daysSinceFailure >= 7 && domain.Status != "suspended" {
				// Suspend 域名
				domain.Status = "suspended"
				s.db.Save(&domain)
				telegram.SendHealthAlert(domain.FullDomain,
					[]string{fmt.Sprintf("Domain down for %d days", daysSinceFailure)},
					"Domain SUSPENDED")
			} else if daysSinceFailure == 25 || daysSinceFailure == 28 {
				// 发送删除警告
				telegram.SendDeletionWarning(domain.FullDomain, 30-daysSinceFailure)
			}
		}
	} else {
		// 域名恢复正常，清除失败记录
		if domain.FirstFailedAt != nil {
			domain.FirstFailedAt = nil
			if domain.Status == "suspended" {
				domain.Status = "active"
			}
			s.db.Save(&domain)
			telegram.SendHealthAlert(domain.FullDomain,
				[]string{"Domain is back online"},
				"Recovery detected")
		}
	}
}

// StartPeriodicScanning 启动定期扫描 (可选功能)
func (s *Scanner) StartPeriodicScanning(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 立即执行一次
	s.ScanAllDomains(ctx)
	s.ScanAllPendingDomains(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.ScanAllDomains(ctx)
			s.ScanAllPendingDomains(ctx)
		}
	}
}
