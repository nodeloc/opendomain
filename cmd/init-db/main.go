package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"opendomain/internal/config"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 创建 api_quotas 表
	sqlStmt := `
CREATE TABLE IF NOT EXISTS api_quotas (
    id SERIAL PRIMARY KEY,
    api_name VARCHAR(100) NOT NULL,
    date VARCHAR(10) NOT NULL,
    used_count INTEGER NOT NULL DEFAULT 0,
    daily_limit INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 为 api_name 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_api_quotas_api_name ON api_quotas(api_name);

-- 为日期查询创建索引
CREATE INDEX IF NOT EXISTS idx_api_quotas_date ON api_quotas(date);
`

	if err := db.Exec(sqlStmt).Error; err != nil {
		log.Fatalf("Failed to create api_quotas table: %v", err)
	}

	fmt.Println("✅ Successfully created api_quotas table and indexes")

	// 验证表是否创建成功
	var count int64
	if err := db.Raw("SELECT COUNT(*) FROM api_quotas").Scan(&count).Error; err != nil {
		log.Fatalf("Failed to verify table: %v", err)
	}

	fmt.Printf("✅ Table verified, current records: %d\n", count)
}
