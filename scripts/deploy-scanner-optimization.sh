#!/bin/bash

# 扫描器优化部署脚本

echo "=== OpenDomain 扫描器速率限制优化 ==="
echo ""

# 1. 运行数据库迁移
echo "步骤 1: 运行数据库迁移..."
make migrate-up

if [ $? -ne 0 ]; then
    echo "❌ 数据库迁移失败"
    exit 1
fi

echo "✅ 数据库迁移完成"
echo ""

# 2. 重新编译
echo "步骤 2: 重新编译项目..."
make build

if [ $? -ne 0 ]; then
    echo "❌ 编译失败"
    exit 1
fi

echo "✅ 编译完成"
echo ""

# 3. 重启服务
echo "步骤 3: 重启服务..."
if [ -f "docker-compose.yml" ]; then
    docker-compose restart api
    echo "✅ Docker 服务已重启"
elif command -v systemctl &> /dev/null; then
    sudo systemctl restart opendomain
    echo "✅ Systemd 服务已重启"
else
    echo "⚠️  请手动重启服务"
fi

echo ""
echo "=== 部署完成 ==="
echo ""
echo "📊 查看配额状态: curl http://localhost:8080/api/admin/api-quota"
echo "📖 查看文档: docs/SCANNER_RATE_LIMIT.md"
