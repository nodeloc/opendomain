#!/bin/bash

# 恢复被错误暂停的域名脚本
# 使用方法：./recover_domains.sh

echo "=== OpenDomain 域名恢复工具 ==="
echo ""

# 方法 1: 使用 Go 脚本（推荐）
echo "方法 1: 使用 Go 脚本批量恢复"
echo "----------------------------------------"
echo "1. 设置数据库环境变量："
echo "   export DB_HOST=localhost"
echo "   export DB_PORT=5432"
echo "   export DB_USER=postgres"
echo "   export DB_PASSWORD=your_password"
echo "   export DB_NAME=opendomain"
echo ""
echo "2. 运行恢复脚本："
echo "   go run scripts/recover_domains.go"
echo ""

# 方法 2: 使用 SQL 脚本
echo "方法 2: 使用 SQL 脚本手动恢复"
echo "----------------------------------------"
echo "   psql -U postgres -d opendomain -f scripts/recover_suspended_domains.sql"
echo ""

# 方法 3: 快速恢复单个域名
echo "方法 3: 快速恢复单个域名（dsm.loc.cc）"
echo "----------------------------------------"

read -p "是否立即恢复域名 dsm.loc.cc？(yes/no): " confirm

if [ "$confirm" = "yes" ] || [ "$confirm" = "y" ]; then
    echo "正在恢复域名..."
    
    # 从配置文件或环境变量读取数据库信息
    DB_HOST=${DB_HOST:-localhost}
    DB_PORT=${DB_PORT:-5432}
    DB_USER=${DB_USER:-postgres}
    DB_NAME=${DB_NAME:-opendomain}
    
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "
        UPDATE domains 
        SET status = 'active', first_failed_at = NULL 
        WHERE full_domain = 'dsm.loc.cc';
        
        SELECT id, full_domain, status 
        FROM domains 
        WHERE full_domain = 'dsm.loc.cc';
    "
    
    if [ $? -eq 0 ]; then
        echo "✅ 域名恢复成功！"
        echo ""
        echo "注意：如果需要同步 PowerDNS 记录，请通过管理面板再次切换一次状态。"
    else
        echo "❌ 域名恢复失败，请检查数据库连接。"
    fi
else
    echo "操作已取消。"
fi

echo ""
echo "=== 其他选项 ==="
echo "- 通过管理面板恢复：登录管理后台 -> Root Domains -> 找到对应域名 -> 点击 Activate"
echo "- 查看所有被暂停的域名："
echo "  psql -U $DB_USER -d $DB_NAME -c \"SELECT id, full_domain, status, first_failed_at FROM domains WHERE status = 'suspended';\""
