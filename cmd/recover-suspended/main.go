package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"opendomain/internal/models"
)

func main() {
	// 从环境变量获取数据库连接信息
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "opendomain")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("=== 恢复被错误暂停的域名 ===")
	fmt.Println("正在查找被暂停的域名...")

	// 查找所有被暂停的域名
	var suspendedDomains []models.Domain
	if err := db.Where("status = ?", "suspended").Find(&suspendedDomains).Error; err != nil {
		log.Fatalf("Failed to query domains: %v", err)
	}

	if len(suspendedDomains) == 0 {
		fmt.Println("没有找到被暂停的域名")
		return
	}

	fmt.Printf("找到 %d 个被暂停的域名\n\n", len(suspendedDomains))

	// 统计信息
	var recovered = 0
	var skipped = 0

	for _, domain := range suspendedDomains {
		// 检查该域名的扫描摘要
		var summary models.DomainScanSummary
		if err := db.Where("domain_id = ?", domain.ID).First(&summary).Error; err != nil {
			fmt.Printf("❌ %s - 未找到扫描记录，跳过\n", domain.FullDomain)
			skipped++
			continue
		}

		// 判断是否应该恢复
		shouldRecover := false
		reason := ""

		// 如果 Safe Browsing 状态不是真正的 unsafe，应该恢复
		if summary.SafeBrowsingStatus == "safe" || summary.SafeBrowsingStatus == "unknown" || summary.SafeBrowsingStatus == "" {
			if summary.VirusTotalStatus != "malicious" {
				shouldRecover = true
				reason = fmt.Sprintf("Safe Browsing: %s, VirusTotal: %s (非恶意)",
					getStatusString(summary.SafeBrowsingStatus),
					getStatusString(summary.VirusTotalStatus))
			}
		}

		// 如果域名正常但被错误暂停
		if summary.OverallHealth == "healthy" {
			shouldRecover = true
			reason = "域名健康状态正常"
		}

		if shouldRecover {
			// 恢复域名
			domain.Status = "active"
			domain.FirstFailedAt = nil
			if err := db.Save(&domain).Error; err != nil {
				fmt.Printf("❌ %s - 恢复失败: %v\n", domain.FullDomain, err)
				skipped++
			} else {
				fmt.Printf("✅ %s - 已恢复 (%s)\n", domain.FullDomain, reason)
				recovered++
			}
		} else {
			fmt.Printf("⚠️  %s - 保持暂停状态 (Safe Browsing: %s, VirusTotal: %s)\n",
				domain.FullDomain,
				summary.SafeBrowsingStatus,
				summary.VirusTotalStatus)
			skipped++
		}
	}

	fmt.Printf("\n=== 完成 ===\n")
	fmt.Printf("恢复: %d 个\n", recovered)
	fmt.Printf("跳过: %d 个\n", skipped)
	fmt.Printf("总计: %d 个\n", len(suspendedDomains))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getStatusString(status string) string {
	if status == "" {
		return "未知"
	}
	return status
}
