#!/bin/bash
# 重置管理员密码脚本

set -e

echo "OpenDomain 管理员密码重置工具"
echo "================================"
echo ""

# 读取新密码
read -p "请输入新密码: " NEW_PASSWORD
echo ""

if [ -z "$NEW_PASSWORD" ]; then
    echo "错误: 密码不能为空"
    exit 1
fi

# 生成 bcrypt 哈希（使用 Go 程序）
echo "正在生成密码哈希..."

# 创建临时 Go 程序来生成 bcrypt 哈希
cat > /tmp/hash_password.go <<'EOF'
package main

import (
    "fmt"
    "os"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: hash_password <password>")
        os.Exit(1)
    }

    password := os.Args[1]
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Print(string(hash))
}
EOF

# 使用 Docker 容器编译并运行
PASSWORD_HASH=$(docker run --rm -v /tmp:/tmp golang:1.24-alpine sh -c "
cd /tmp && \
go mod init hash_password 2>/dev/null || true && \
go get golang.org/x/crypto/bcrypt && \
go run hash_password.go '$NEW_PASSWORD'
")

if [ -z "$PASSWORD_HASH" ]; then
    echo "错误: 生成密码哈希失败"
    rm -f /tmp/hash_password.go
    exit 1
fi

echo "密码哈希生成成功"

# 更新数据库
echo "正在更新数据库..."
docker exec opendomain_db psql -U opendomain -d opendomain -c \
    "UPDATE users SET password_hash = '$PASSWORD_HASH' WHERE email = 'admin@opendomain.local';"

# 清理
rm -f /tmp/hash_password.go

echo ""
echo "✓ 管理员密码重置成功！"
echo ""
echo "邮箱: admin@opendomain.local"
echo "新密码: $NEW_PASSWORD"
echo ""
