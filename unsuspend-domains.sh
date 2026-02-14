#!/bin/bash
# 域名恢复脚本 - 将所有挂起(suspended)的域名恢复为活跃(active)状态

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$SCRIPT_DIR"

echo "========================================"
echo "OpenDomain - 域名批量恢复工具"
echo "========================================"
echo ""

# 检查是否在项目根目录
if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
    echo "错误: 请在项目根目录运行此脚本"
    exit 1
fi

# 检查 .env 文件
if [ ! -f "$PROJECT_ROOT/.env" ]; then
    echo "警告: 未找到 .env 文件，将使用环境变量"
fi

cd "$PROJECT_ROOT"

# 编译程序
echo "正在编译恢复工具..."
if ! go build -o bin/unsuspend-domains cmd/unsuspend-domains/main.go; then
    echo "错误: 编译失败"
    exit 1
fi
echo "✓ 编译成功"
echo ""

# 运行程序
echo "运行恢复工具..."
echo ""
./bin/unsuspend-domains

echo ""
echo "如需查看当前挂起的域名列表，可以运行:"
echo "  docker exec opendomain_db psql -U opendomain -d opendomain -c \"SELECT id, full_domain, status, updated_at FROM domains WHERE status = 'suspended' ORDER BY updated_at DESC;\""
