package main

import (
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
	// 加载 .env 文件
	_ = godotenv.Load()

	viper.SetEnvPrefix("") // 不使用前缀
	viper.AutomaticEnv()

	config := &Config{}

	// 数据库配置
	config.Database.Host = viper.GetString("DB_HOST")
	config.Database.Port = viper.GetInt("DB_PORT")
	config.Database.User = viper.GetString("DB_USER")
	config.Database.Password = viper.GetString("DB_PASSWORD")
	config.Database.DBName = viper.GetString("DB_NAME")
	config.Database.SSLMode = viper.GetString("DB_SSL_MODE")

	// 设置默认值
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
	fmt.Println("========================================")
	fmt.Println("域名恢复工具 - Unsuspend All Domains")
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

	// 查询所有 suspended 的域名
	var suspendedDomains []models.Domain
	if err := db.Where("status = ?", "suspended").Find(&suspendedDomains).Error; err != nil {
		log.Fatalf("Failed to query suspended domains: %v", err)
	}

	if len(suspendedDomains) == 0 {
		fmt.Println("✓ 没有找到被挂起的域名")
		return
	}

	fmt.Printf("找到 %d 个被挂起的域名:\n", len(suspendedDomains))
	fmt.Println("----------------------------------------")
	for i, domain := range suspendedDomains {
		fmt.Printf("%d. %s (ID: %d)\n", i+1, domain.FullDomain, domain.ID)
	}
	fmt.Println("----------------------------------------")
	fmt.Println()

	// 确认操作
	fmt.Print("是否要恢复这些域名？(y/n): ")
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" && confirm != "yes" && confirm != "YES" {
		fmt.Println("操作已取消")
		return
	}

	fmt.Println()
	fmt.Println("开始恢复域名...")
	fmt.Println()

	// 批量更新状态
	successCount := 0
	failedCount := 0
	now := time.Now()

	for _, domain := range suspendedDomains {
		// 更新状态
		updates := map[string]interface{}{
			"status":     "active",
			"updated_at": now,
		}

		// 清除首次失败时间（如果有的话）
		if domain.FirstFailedAt != nil {
			updates["first_failed_at"] = nil
		}

		if err := db.Model(&models.Domain{}).Where("id = ?", domain.ID).Updates(updates).Error; err != nil {
			fmt.Printf("✗ 失败: %s - %v\n", domain.FullDomain, err)
			failedCount++
		} else {
			fmt.Printf("✓ 已恢复: %s\n", domain.FullDomain)
			successCount++
		}
	}

	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("恢复完成！成功: %d, 失败: %d\n", successCount, failedCount)
	fmt.Println("========================================")
}
