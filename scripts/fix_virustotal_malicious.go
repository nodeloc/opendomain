package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"opendomain/internal/config"
	"opendomain/internal/models"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(cfg.Database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("========================================")
	fmt.Println("VirusTotal Malicious Status Fix Script")
	fmt.Println("========================================")
	fmt.Println()

	// 1. 查找所有标记为 malicious 的域名摘要
	var summaries []models.DomainScanSummary
	err = db.Where("virus_total_status = ?", "malicious").Find(&summaries).Error
	if err != nil {
		log.Fatalf("Failed to query summaries: %v", err)
	}

	fmt.Printf("Found %d domains marked as 'malicious'\n\n", len(summaries))

	if len(summaries) == 0 {
		fmt.Println("No domains to fix. Exiting.")
		return
	}

	fixed := 0
	alreadyCorrect := 0
	actuallyMalicious := 0

	for _, summary := range summaries {
		// 获取域名信息
		var domain models.Domain
		if err := db.First(&domain, summary.DomainID).Error; err != nil {
			fmt.Printf("⚠️  Skip: Cannot find domain ID %d\n", summary.DomainID)
			continue
		}

		// 查找最近的 VirusTotal 扫描记录
		var vtScan models.DomainScan
		err := db.Where("domain_id = ? AND scan_type = ?", summary.DomainID, "virustotal").
			Order("scanned_at DESC").
			First(&vtScan).Error

		if err != nil {
			fmt.Printf("⚠️  %s: No VirusTotal scan record found\n", domain.FullDomain)
			continue
		}

		// 检查扫描状态
		if vtScan.Status == "failed" || vtScan.Status == "quota_exceeded" {
			// 这是 API 失败导致的误判，需要修复
			summary.VirusTotalStatus = "unknown"

			// 重新计算整体健康状态
			recalculateOverallHealth(&summary)

			// 保存更新
			if err := db.Save(&summary).Error; err != nil {
				fmt.Printf("❌ %s: Failed to update - %v\n", domain.FullDomain, err)
				continue
			}

			fmt.Printf("✅ %s: Fixed (scan status: %s, error: %s)\n",
				domain.FullDomain, vtScan.Status, vtScan.ErrorMessage)
			fixed++
		} else if vtScan.Status == "success" {
			// 检查扫描详情，确认是否真的是恶意
			if vtScan.ScanDetails != nil {
				details := *vtScan.ScanDetails
				// 如果包含 "malicious":0，说明实际上是安全的
				if contains(details, `"malicious":0`) {
					summary.VirusTotalStatus = "clean"
					recalculateOverallHealth(&summary)
					db.Save(&summary)
					fmt.Printf("✅ %s: Fixed (actually clean, malicious count: 0)\n", domain.FullDomain)
					fixed++
				} else {
					fmt.Printf("⚪ %s: Confirmed malicious (scan details: %s)\n",
						domain.FullDomain, truncate(*vtScan.ScanDetails, 100))
					actuallyMalicious++
				}
			} else {
				fmt.Printf("⚪ %s: Confirmed malicious (no details available)\n", domain.FullDomain)
				actuallyMalicious++
			}
		} else {
			fmt.Printf("ℹ️  %s: Unknown scan status '%s'\n", domain.FullDomain, vtScan.Status)
			alreadyCorrect++
		}
	}

	// 打印统计
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Summary:")
	fmt.Println("========================================")
	fmt.Printf("Total processed:       %d\n", len(summaries))
	fmt.Printf("✅ Fixed:              %d\n", fixed)
	fmt.Printf("⚪ Actually malicious: %d\n", actuallyMalicious)
	fmt.Printf("ℹ️  Other:              %d\n", alreadyCorrect)
	fmt.Println()
}

// recalculateOverallHealth 重新计算整体健康状态
func recalculateOverallHealth(summary *models.DomainScanSummary) {
	if summary.SafeBrowsingStatus == "unsafe" || summary.VirusTotalStatus == "malicious" {
		summary.OverallHealth = "degraded"
	} else if summary.DNSStatus == "resolved" && summary.HTTPStatus == "online" {
		summary.OverallHealth = "healthy"
	} else if summary.DNSStatus == "resolved" || summary.HTTPStatus == "online" {
		summary.OverallHealth = "degraded"
	} else {
		summary.OverallHealth = "down"
	}
	summary.UpdatedAt = time.Now()
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsRec(s, substr))
}

func containsRec(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
