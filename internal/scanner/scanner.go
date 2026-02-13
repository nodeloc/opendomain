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
	vtMu            sync.Mutex
	vtLastCall      time.Time
	vtDailyCount    int
	vtDailyResetDate string // 格式: YYYY-MM-DD
}

func NewScanner(db *gorm.DB, cfg *config.Config) *Scanner {
	return &Scanner{db: db, cfg: cfg}
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
	return true
}

// strPtr 返回字符串指针
func strPtr(s string) *string {
	return &s
}

// ScanAllDomains 扫描所有活跃域名
func (s *Scanner) ScanAllDomains(ctx context.Context) error {
	var domains []models.Domain
	if err := s.db.Where("status = ?", "active").Find(&domains).Error; err != nil {
		return fmt.Errorf("failed to fetch domains: %w", err)
	}

	for _, domain := range domains {
		if err := s.ScanDomain(ctx, &domain); err != nil {
			fmt.Printf("Failed to scan domain %s: %v\n", domain.FullDomain, err)
		}
	}

	return nil
}

// ScanAllPendingDomains 扫描所有待激活域名
func (s *Scanner) ScanAllPendingDomains(ctx context.Context) error {
	var pendingDomains []models.PendingDomain
	if err := s.db.Where("deleted_at IS NULL").Find(&pendingDomains).Error; err != nil {
		return fmt.Errorf("failed to fetch pending domains: %w", err)
	}

	fmt.Printf("[INFO] Scanning %d pending domains\n", len(pendingDomains))

	for _, pending := range pendingDomains {
		if err := s.ScanPendingDomain(ctx, &pending); err != nil {
			fmt.Printf("Failed to scan pending domain %s: %v\n", pending.FullDomain, err)
		}
	}

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
		} else {
			safeBrowsingStatus = "unsafe"
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
		scan.Status = "failed"
		details, _ := json.Marshal(map[string]interface{}{
			"safe":    false,
			"threats": matches,
		})
		scan.ScanDetails = strPtr(string(details))
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
	if safeBrowsingStatus == "unsafe" || virusTotalStatus == "malicious" {
		if domain.Status != "suspended" {
			domain.Status = "suspended"
			s.db.Save(&domain)

			reason := fmt.Sprintf("Malicious content detected (Safe Browsing: %s, VirusTotal: %s)",
				safeBrowsingStatus, virusTotalStatus)
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
