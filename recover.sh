#!/bin/bash

# 恢复被错误暂停的域名 - 快速执行脚本
# 使用方法: ./recover.sh

echo "=== 开始恢复被错误暂停的域名 ==="
echo ""

# 检查是否已编译
if [ ! -f "bin/recover-suspended" ]; then
    echo "工具未编译，正在编译..."
    make build-tools
    echo ""
fi

# 加载环境变量（如果存在 .env 文件）
if [ -f ".env" ]; then
    echo "加载 .env 配置..."
    export $(cat .env | grep -v '^#' | xargs)
fi

# 执行恢复工具
./bin/recover-suspended

echo ""
echo "完成！"
