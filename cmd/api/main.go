package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"opendomain/internal/config"
	"opendomain/internal/handler"
	"opendomain/internal/i18n"
	"opendomain/internal/router"
	"opendomain/internal/scanner"
	"opendomain/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger.Init(cfg.LogLevel)
	defer logger.Sync()

	logger.Info("Starting OpenDomain API Server...")
	logger.Infof("Environment: %s", cfg.Env)

	// 初始化数据库
	db, err := config.InitDatabase(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	logger.Info("Database connected successfully")

	// 初始化 Redis
	rdb, err := config.InitRedis(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize Redis: %v", err)
	}
	logger.Info("Redis connected successfully")

	// 初始化 i18n
	localesPath := filepath.Join("internal", "i18n", "locales")
	if err := i18n.Init(localesPath); err != nil {
		logger.Warnf("Failed to initialize i18n: %v (continuing without i18n)", err)
	} else {
		logger.Info("I18n initialized successfully")
	}

	// 设置 Gin 模式
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建路由
	r := router.Setup(db, rdb, cfg)

	// 启动定时健康扫描任务
	scannerCtx, scannerCancel := context.WithCancel(context.Background())
	defer scannerCancel()

	healthScanner := scanner.NewScanner(db, cfg)
	go func() {
		logger.Info("Starting periodic health scanning (every 1 hour)...")
		healthScanner.StartPeriodicScanning(scannerCtx, 1*time.Hour)
	}()

	// 启动过期域名自动清理任务
	domainHandler := handler.NewDomainHandler(db, cfg)
	go func() {
		logger.Info("Starting periodic expired domain cleanup (every 24 hours)...")
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		// 立即执行一次清理
		domainHandler.CleanupExpiredDomains(30) // 删除过期超过 30 天的域名

		for {
			select {
			case <-ticker.C:
				logger.Info("Running scheduled expired domain cleanup...")
				domainHandler.CleanupExpiredDomains(30)
			case <-scannerCtx.Done():
				logger.Info("Stopping expired domain cleanup task...")
				return
			}
		}
	}()

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	go func() {
		logger.Infof("Server listening on port %s", cfg.Port)
		logger.Infof("API Documentation: http://localhost:%s/swagger/index.html", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}
