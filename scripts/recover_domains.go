package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Domain struct {
	ID             uint   `gorm:"primaryKey"`
	FullDomain     string `gorm:"column:full_domain"`
	Status         string
	FirstFailedAt  *string `gorm:"column:first_failed_at"`
}

type DomainScan struct {
	ID           uint
	DomainID     uint   `gorm:"column:domain_id"`
	ScanType     string `gorm:"column:scan_type"`
	Status       string
	ErrorMessage string `gorm:"column:error_message"`
}

func (Domain) TableName() string {
	return "domains"
}

func (DomainScan) TableName() string {
	return "domain_scans"
}

func main() {
	// 从环境变量读取数据库连接信息
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "opendomain")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("=== 查找被错误暂停的域名 ===\n")

	// 查找所有被暂停的域名
	var suspendedDomains []Domain
	if err := db.Where("status = ?", "suspended").Find(&suspendedDomains).Error; err != nil {
		log.Fatalf("Failed to query suspended domains: %v", err)
	}

	if len(suspendedDomains) == 0 {
		fmt.Println("没有找到被暂停的域名。")
		return
	}

	fmt.Printf("找到 %d 个被暂停的域名：\n\n", len(suspendedDomains))

	// 检查每个域名的扫描记录
	var toRecover []uint
	for _, domain := range suspendedDomains {
		fmt.Printf("域名: %s (ID: %d)\n", domain.FullDomain, domain.ID)

		// 查找最近的 Safe Browsing 扫描记录
		var recentScans []DomainScan
		db.Where("domain_id = ? AND scan_type = ?", domain.ID, "safebrowsing").
			Order("id DESC").
			Limit(3).
			Find(&recentScans)

		if len(recentScans) == 0 {
			fmt.Println("  ❌ 没有找到 Safe Browsing 扫描记录")
			continue
		}

		// 检查是否是因为 API 失败
		hasAPIFailure := false
		hasThreatDetection := false

		for _, scan := range recentScans {
			if scan.Status == "failed" && scan.ErrorMessage != "" {
				hasAPIFailure = true
				fmt.Printf("  ⚠️  API 调用失败: %s\n", scan.ErrorMessage)
			} else if scan.Status == "threat_detected" {
				hasThreatDetection = true
				fmt.Println("  ⚠️  检测到真实威胁")
			}
		}

		if hasAPIFailure && !hasThreatDetection {
			fmt.Println("  ✅ 建议恢复（疑似 API 失败误判）")
			toRecover = append(toRecover, domain.ID)
		} else if hasThreatDetection {
			fmt.Println("  ❌ 不建议恢复（检测到真实威胁）")
		}

		fmt.Println()
	}

	if len(toRecover) == 0 {
		fmt.Println("没有需要恢复的域名。")
		return
	}

	fmt.Printf("\n=== 准备恢复 %d 个域名 ===\n", len(toRecover))
	fmt.Println("这些域名将被恢复为 active 状态。")
	fmt.Print("是否继续？(yes/no): ")

	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "yes" && confirm != "y" {
		fmt.Println("操作已取消。")
		return
	}

	// 批量恢复域名
	result := db.Model(&Domain{}).
		Where("id IN ?", toRecover).
		Updates(map[string]interface{}{
			"status":           "active",
			"first_failed_at":  nil,
		})

	if result.Error != nil {
		log.Fatalf("恢复失败: %v", result.Error)
	}

	fmt.Printf("\n✅ 成功恢复 %d 个域名！\n", result.RowsAffected)
	
	// 显示恢复的域名
	var recovered []Domain
	db.Where("id IN ?", toRecover).Find(&recovered)
	fmt.Println("\n恢复的域名列表：")
	for _, d := range recovered {
		fmt.Printf("  - %s (ID: %d)\n", d.FullDomain, d.ID)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
