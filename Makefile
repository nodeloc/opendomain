.PHONY: help build run test clean migrate-up migrate-down migrate-create docker-build docker-up docker-down

help: ## 显示帮助信息
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## 编译项目
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go

build-tools: ## 编译管理工具
	@echo "Building tools..."
	@go build -o bin/unsuspend-domains cmd/unsuspend-domains/main.go
	@go build -o bin/unsuspend-domains-advanced cmd/unsuspend-domains-advanced/main.go
	@echo "Built tools:"
	@echo "  - bin/unsuspend-domains (批量恢复域名工具)"
	@echo "  - bin/unsuspend-domains-advanced (高级恢复工具，支持参数)"

unsuspend-all: ## 恢复所有挂起的域名
	@echo "Building unsuspend tool..."
	@go build -o bin/unsuspend-domains cmd/unsuspend-domains/main.go
	@echo "Running unsuspend tool..."
	@./bin/unsuspend-domains

unsuspend-list: ## 列出所有挂起的域名
	@echo "Building advanced unsuspend tool..."
	@go build -o bin/unsuspend-domains-advanced cmd/unsuspend-domains-advanced/main.go
	@echo "Listing suspended domains..."
	@./bin/unsuspend-domains-advanced -list

unsuspend-auto: ## 自动恢复所有挂起的域名（不询问）
	@echo "Building advanced unsuspend tool..."
	@go build -o bin/unsuspend-domains-advanced cmd/unsuspend-domains-advanced/main.go
	@echo "Auto-unsuspending domains..."
	@./bin/unsuspend-domains-advanced -y

run: ## 运行 API 服务
	@echo "Starting API server..."
	@go run cmd/api/main.go

test: ## 运行测试
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## 查看测试覆盖率
	@go tool cover -html=coverage.out

clean: ## 清理构建文件
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out

deps: ## 下载依赖
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

migrate-up: ## 执行数据库迁移
	@echo "Running migrations..."
	@migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" up

migrate-down: ## 回滚数据库迁移
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" down

migrate-create: ## 创建新的迁移文件 (usage: make migrate-create name=create_users_table)
	@migrate create -ext sql -dir migrations -seq $(name)

docker-build: ## 构建 Docker 镜像
	@echo "Building Docker images..."
	@docker-compose build

docker-up: ## 启动 Docker 容器
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down: ## 停止 Docker 容器
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs: ## 查看 Docker 日志
	@docker-compose logs -f

dev: ## 开发模式（使用 air 热重载）
	@echo "Starting development mode..."
	@air

.DEFAULT_GOAL := help
