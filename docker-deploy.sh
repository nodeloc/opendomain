#!/bin/bash

# OpenDomain Docker 部署脚本
# 用法: ./docker-deploy.sh [选项]
#   --build    强制重新构建镜像
#   --down     停止并删除容器
#   --logs     查看日志
#   --restart  重启服务

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目目录
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}OpenDomain Docker 部署脚本${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

# 检查 Docker 和 Docker Compose
check_dependencies() {
    echo -e "${GREEN}[1/7] 检查依赖...${NC}"
    
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}错误: Docker 未安装${NC}"
        echo "请安装 Docker: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        echo -e "${RED}错误: Docker Compose 未安装${NC}"
        echo "请安装 Docker Compose: https://docs.docker.com/compose/install/"
        exit 1
    fi
    
    # 使用新版或旧版 docker-compose
    if docker compose version &> /dev/null; then
        DOCKER_COMPOSE="docker compose"
    else
        DOCKER_COMPOSE="docker-compose"
    fi
    
    echo "✓ Docker 版本: $(docker --version)"
    echo "✓ Docker Compose 版本: $($DOCKER_COMPOSE --version)"
}

# 检查环境配置
check_env() {
    echo ""
    echo -e "${GREEN}[2/7] 检查环境配置...${NC}"
    
    if [ ! -f ".env" ]; then
        echo -e "${YELLOW}警告: .env 文件不存在${NC}"
        echo "正在从 .env.example 创建 .env 文件..."
        cp .env.example .env
        echo -e "${YELLOW}请编辑 .env 文件并配置必要的环境变量${NC}"
        echo ""
        echo "必须配置的变量："
        echo "  - DB_PASSWORD (数据库密码)"
        echo "  - JWT_SECRET (JWT 密钥)"
        echo "  - POWERDNS_API_KEY (PowerDNS API 密钥)"
        echo "  - FRONTEND_URL (前端 URL)"
        echo ""
        read -p "按 Enter 继续编辑 .env 文件..." 
        ${EDITOR:-vim} .env
    fi
    
    echo "✓ 环境配置文件存在"
}

# 拉取最新代码
pull_code() {
    echo ""
    echo -e "${GREEN}[3/7] 拉取最新代码...${NC}"
    
    if [ -d ".git" ]; then
        git pull origin main || git pull origin master
        echo "✓ 代码已更新"
    else
        echo -e "${YELLOW}! 不是 Git 仓库，跳过代码拉取${NC}"
    fi
}

# 构建镜像
build_images() {
    echo ""
    echo -e "${GREEN}[4/7] 构建 Docker 镜像...${NC}"
    
    if [ "$FORCE_BUILD" = true ]; then
        echo "强制重新构建镜像..."
        $DOCKER_COMPOSE build --no-cache
    else
        $DOCKER_COMPOSE build
    fi
    
    echo "✓ 镜像构建完成"
}

# 停止旧容器
stop_containers() {
    echo ""
    echo -e "${GREEN}[5/7] 停止旧容器...${NC}"
    
    $DOCKER_COMPOSE down
    echo "✓ 旧容器已停止"
}

# 启动服务
start_services() {
    echo ""
    echo -e "${GREEN}[6/7] 启动服务...${NC}"

    # 启动服务（后台运行）
    $DOCKER_COMPOSE up -d

    echo "✓ 服务已启动"
}

# 复制前端文件
copy_frontend() {
    echo ""
    echo -e "${GREEN}[6.5/7] 复制前端文件到宿主机...${NC}"

    # 等待应用容器启动
    sleep 3

    # 创建目标目录
    FRONTEND_DIR="$PROJECT_DIR/public"
    mkdir -p "$FRONTEND_DIR"

    # 从容器复制前端文件
    if docker cp opendomain-app:/app/web/dist/. "$FRONTEND_DIR/"; then
        echo "✓ 前端文件已复制到: $FRONTEND_DIR"
        echo ""
        echo "请在您的 Nginx 配置中将网站根目录指向: $FRONTEND_DIR"
        echo "并配置反向代理："
        echo "  location /api {"
        echo "    proxy_pass http://localhost:8000;"
        echo "  }"
    else
        echo -e "${RED}✗ 复制前端文件失败${NC}"
        echo "您可以手动复制: docker cp opendomain-app:/app/web/dist/. ./public/"
    fi
}

# 显示状态
show_status() {
    echo ""
    echo -e "${GREEN}[7/7] 检查服务状态...${NC}"
    echo ""
    
    $DOCKER_COMPOSE ps
    
    echo ""
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}部署完成！${NC}"
    echo -e "${GREEN}=========================================${NC}"
    echo ""
    echo "服务访问地址："
    echo "  - API 后端: http://localhost:8000"
    echo "  - 健康检查: http://localhost:8000/health"
    echo "  - 前端文件: $PROJECT_DIR/public/"
    echo ""
    echo "Nginx 配置示例："
    echo "  root $PROJECT_DIR/public;"
    echo "  location /api {"
    echo "    proxy_pass http://localhost:8000;"
    echo "  }"
    echo ""
    echo "常用命令："
    echo "  查看日志: $DOCKER_COMPOSE logs -f"
    echo "  查看应用日志: $DOCKER_COMPOSE logs -f app"
    echo "  重启服务: $DOCKER_COMPOSE restart"
    echo "  停止服务: $DOCKER_COMPOSE down"
    echo "  进入容器: $DOCKER_COMPOSE exec app sh"
    echo "  复制前端: docker cp opendomain-app:/app/web/dist/. ./public/"
    echo ""
}

# 查看日志
view_logs() {
    echo -e "${BLUE}查看服务日志...${NC}"
    $DOCKER_COMPOSE logs -f
}

# 处理命令行参数
FORCE_BUILD=false
ACTION="deploy"

while [[ $# -gt 0 ]]; do
    case $1 in
        --build)
            FORCE_BUILD=true
            shift
            ;;
        --down)
            ACTION="down"
            shift
            ;;
        --logs)
            ACTION="logs"
            shift
            ;;
        --restart)
            ACTION="restart"
            shift
            ;;
        *)
            echo -e "${RED}未知参数: $1${NC}"
            echo "用法: $0 [--build] [--down] [--logs] [--restart]"
            exit 1
            ;;
    esac
done

# 执行操作
case $ACTION in
    deploy)
        check_dependencies
        check_env
        pull_code
        build_images
        stop_containers
        start_services
        copy_frontend
        show_status
        ;;
    down)
        echo -e "${YELLOW}停止并删除所有容器...${NC}"
        $DOCKER_COMPOSE down -v
        echo -e "${GREEN}✓ 容器已停止并删除${NC}"
        ;;
    logs)
        view_logs
        ;;
    restart)
        echo -e "${YELLOW}重启服务...${NC}"
        $DOCKER_COMPOSE restart
        echo -e "${GREEN}✓ 服务已重启${NC}"
        $DOCKER_COMPOSE ps
        ;;
esac
