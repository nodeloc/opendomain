#!/bin/bash

# OpenDomain 部署脚本
# 用法: ./deploy.sh

set -e  # 遇到错误立即退出

echo "========================================="
echo "OpenDomain 部署脚本"
echo "========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 项目目录（脚本所在目录）
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

echo -e "${GREEN}[1/8] 拉取最新代码...${NC}"
git pull origin main || git pull origin master

echo ""
echo -e "${GREEN}[2/8] 检查环境配置...${NC}"
if [ ! -f ".env" ]; then
    echo -e "${RED}错误: .env 文件不存在${NC}"
    echo "请复制 .env.example 并配置环境变量："
    echo "  cp .env.example .env"
    echo "  vim .env"
    exit 1
fi

echo ""
echo -e "${GREEN}[3/8] 安装 Go 依赖...${NC}"
go mod download
go mod tidy

echo ""
echo -e "${GREEN}[4/8] 编译后端程序...${NC}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o opendomain ./cmd/api
chmod +x opendomain
echo "后端编译完成: ./opendomain"

echo ""
echo -e "${GREEN}[5/8] 运行数据库迁移...${NC}"
if command -v migrate &> /dev/null; then
    # 从 .env 读取数据库配置
    export $(cat .env | grep -v '^#' | xargs)
    DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}"
    migrate -path ./migrations -database "$DB_URL" up
    echo "数据库迁移完成"
else
    echo -e "${YELLOW}警告: migrate 命令未找到，跳过数据库迁移${NC}"
    echo "如需安装 migrate: https://github.com/golang-migrate/migrate"
fi

echo ""
echo -e "${GREEN}[6/8] 安装前端依赖...${NC}"
cd web
if [ -f "package-lock.json" ]; then
    npm ci
else
    npm install
fi

echo ""
echo -e "${GREEN}[7/8] 构建前端应用...${NC}"
npm run build
echo "前端构建完成: web/dist/"
cd ..

echo ""
echo -e "${GREEN}[8/8] 重启服务...${NC}"

# 检查是否使用 systemd
if systemctl is-active --quiet opendomain; then
    echo "使用 systemd 重启服务..."
    sudo systemctl restart opendomain
    sudo systemctl status opendomain --no-pager
    echo -e "${GREEN}服务已通过 systemd 重启${NC}"
# 检查是否有运行中的进程
elif pgrep -f "./opendomain" > /dev/null; then
    echo "检测到运行中的进程，正在重启..."
    pkill -f "./opendomain"
    sleep 2
    nohup ./opendomain > logs/app.log 2>&1 &
    echo -e "${GREEN}进程已重启 (PID: $!)${NC}"
else
    echo -e "${YELLOW}未检测到运行中的服务${NC}"
    echo "要启动服务，请运行："
    echo "  nohup ./opendomain > logs/app.log 2>&1 &"
    echo "或创建 systemd 服务"
fi

echo ""
echo "========================================="
echo -e "${GREEN}部署完成！${NC}"
echo "========================================="
echo ""
echo "服务信息："
echo "  - 后端程序: ./opendomain"
echo "  - 前端文件: web/dist/"
echo "  - 日志文件: logs/app.log (如果使用 nohup)"
echo ""
echo "有用的命令："
echo "  查看日志: tail -f logs/app.log"
echo "  查看进程: ps aux | grep opendomain"
echo "  停止服务: pkill -f opendomain"
echo ""
