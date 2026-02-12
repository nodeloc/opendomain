#!/bin/sh
set -e

# 复制前端文件到共享 volume（如果目标目录为空）
if [ -d "/app/web/dist" ] && [ ! -f "/shared/dist/.copied" ]; then
    echo "Copying frontend files to shared volume..."
    mkdir -p /shared/dist
    cp -r /app/web/dist/* /shared/dist/
    touch /shared/dist/.copied
    echo "Frontend files copied successfully"
fi

# 执行原始命令
exec "$@"
