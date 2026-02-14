package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"opendomain/internal/models"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
}

func loadConfig() (*Config, error) {
	_ = godotenv.Load()

	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	config := &Config{}

	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetInt("DB_PORT")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")
	config.Database.DBName = viper.GetString("DB_NAME")
	config.Database.SSLMode = viper.GetString("DB_SSL_MODE")

	if config.Database.Host == "" {
		config.Database.Host = "localhost"
	}
	if config.Database.Port == 0 {
		config.Database.Port = 5432
	}
	if config.Database.SSLMode == "" {
		config.Database.SSLMode = "disable"
	}

	return config, nil
}

func connectDB(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.DBName,
		config.Database.Port,
		config.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func main() {
	// 命令行参数
	listOnly := flag.Bool("list", false, "只列出挂起的域名，不执行恢复")
	autoYes := flag.Bool("y", false, "自动确认，不询问")
	domain := flag.String("domain", "", "指定要恢复的域名（支持模糊匹配）")
	help := flag.Bool("h", false, "显示帮助信息")

	flag.Parse()

	if *help {
		fmt.Println("域名恢复工具 - 将挂起的域名恢复为活跃状态")
		fmt.Println()
		fmt.Println("用法:")
		fmt.Println("  unsuspend-domains [选项]")
		fmt.Println()
		fmt.Println("选项:")
		fmt.Println("  -list          只列出挂起的域名，不执行恢复")
		fmt.Println("  -y             自动确认，不询问")
		fmt.Println("  -domain <名称> 指定要恢复的域名（支持模糊匹配，如 'example.com' 或 '%.example.com'）")
		fmt.Println("  -h             显示此帮助信息")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  unsuspend-domains                    # 交互式恢复所有挂起的域名")
		fmt.Println("  unsuspend-domains -list              # 列出所有挂起的域名")
		fmt.Println("  unsuspend-domains -y                 # 自动恢复所有挂起的域名")
		fmt.Println("  unsuspend-domains -domain test.com   # 恢复包含 'test.com' 的域名")
		fmt.Println()
		return
	}

	fmt.Println("========================================")
	fmt.Println("域名恢复工具 - Unsuspend Domains")
	fmt.Println("========================================")
	fmt.Println()

	// 加载配置
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	db, err := connectDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	fmt.Println("✓ 数据库连接成功")
	fmt.Println()

	// 构建查询
	query := db.Where("status = ?", "suspended")
	if *domain != "" {
		query = query.Where("full_domain LIKE ?", "%"+*domain+"%")
	}

	// 查询挂起的域名
	var suspendedDomains []models.Domain
	if err := query.Order("updated_at DESC").Find(&suspendedDomains).Error; err != nil {
		log.Fatalf("Failed to query suspended domains: %v", err)
	}

	if len(suspendedDomains) == 0 {
		if *domain != "" {
			fmt.Printf("✓ 没有找到匹配 '%s' 的被挂起域名\n", *domain)
		} else {
			fmt.Println("✓ 没有找到被挂起的域名")
		}
		return
	}

	// 显示域名列表
	if *domain != "" {
		fmt.Printf("找到 %d 个匹配 '%s' 的被挂起域名:\n", len(suspendedDomains), *domain)
	} else {
		fmt.Printf("找到 %d 个被挂起的域名:\n", len(suspendedDomains))
	}
	fmt.Println("----------------------------------------")
	for i, d := range suspendedDomains {
		failedInfo := ""
		if d.FirstFailedAt != nil {
			duration := time.Since(*d.FirstFailedAt)
			days := int(duration.Hours() / 24)
			failedInfo = fmt.Sprintf(" (失败 %d 天)", days)
		}
		fmt.Printf("%d. %s (ID: %d)%s\n", i+1, d.FullDomain, d.ID, failedInfo)
	}
	fmt.Println("----------------------------------------")
	fmt.Println()

	// 如果只是列出，则退出
	if *listOnly {
		return
	}

	// 确认操作
	if !*autoYes {
		fmt.Print("是否要恢复这些域名？(y/n): ")
		var confirm string
		fmt.Scanln(&confirm)

		if confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "YES" {
			fmt.Println("操作已取消")
			return
		}
	}

	fmt.Println()
	fmt.Println("开始恢复域名...")
	fmt.Println()

	// 批量更新状态
	successCount := 0
	failedCount := 0
	now := time.Now()

	for _, d := range suspendedDomains {
		updates := map[string]interface{}{
			"status":     "active",
			"updated_at": now,
		}

		// 清除首次失败时间
		if d.FirstFailedAt != nil {
			updates["first_failed_at"] = nil
		}

		if err := db.Model(&models.Domain{}).Where("id = ?", d.ID).Updates(updates).Error; err != nil {
			fmt.Printf("✗ 失败: %s - %v\n", d.FullDomain, err)
			failedCount++
		} else {
			fmt.Printf("✓ 已恢复: %s\n", d.FullDomain)
			successCount++
		}
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("恢复完成！成功: %d, 失败: %d\n", successCount, failedCount)
	fmt.Println("========================================")

	if successCount > 0 {
		fmt.Println()
		fmt.Println("提示: 恢复的域名将在下次健康检查时重新扫描")
	}
}
