# 域名恢复工具

本工具用于将被挂起（suspended）的域名批量恢复为活跃（active）状态。

## 工具列表

### 1. unsuspend-domains（基础版）
交互式的恢复工具，会列出所有挂起的域名并要求确认。

### 2. unsuspend-domains-advanced（高级版）
功能更强大的恢复工具，支持命令行参数。

### 3. Shell 脚本
提供了 `unsuspend-domains.sh` 便捷脚本。

## 使用方法

### 方法一：使用 Makefile（推荐）

```bash
# 编译所有管理工具
make build-tools

# 交互式恢复所有挂起的域名
make unsuspend-all

# 只列出挂起的域名，不恢复
make unsuspend-list

# 自动恢复所有挂起的域名（不询问）
make unsuspend-auto
```

### 方法二：使用 Shell 脚本

```bash
# 给脚本添加执行权限
chmod +x unsuspend-domains.sh

# 运行脚本
./unsuspend-domains.sh
```

### 方法三：直接运行 Go 程序

```bash
# 编译基础版
go build -o bin/unsuspend-domains cmd/unsuspend-domains/main.go

# 运行
./bin/unsuspend-domains
```

### 方法四：使用高级版工具

```bash
# 编译高级版
go build -o bin/unsuspend-domains-advanced cmd/unsuspend-domains-advanced/main.go

# 显示帮助
./bin/unsuspend-domains-advanced -h

# 只列出挂起的域名
./bin/unsuspend-domains-advanced -list

# 自动恢复所有域名（不询问）
./bin/unsuspend-domains-advanced -y

# 恢复指定域名（支持模糊匹配）
./bin/unsuspend-domains-advanced -domain example.com

# 组合使用：自动恢复匹配的域名
./bin/unsuspend-domains-advanced -domain test -y
```

## 高级版参数说明

| 参数 | 说明 | 示例 |
|------|------|------|
| `-list` | 只列出挂起的域名，不执行恢复 | `./bin/unsuspend-domains-advanced -list` |
| `-y` | 自动确认，不询问 | `./bin/unsuspend-domains-advanced -y` |
| `-domain <名称>` | 指定要恢复的域名（支持模糊匹配） | `./bin/unsuspend-domains-advanced -domain test.com` |
| `-h` | 显示帮助信息 | `./bin/unsuspend-domains-advanced -h` |

## 使用场景

### 场景 1：查看所有挂起的域名
```bash
make unsuspend-list
# 或
./bin/unsuspend-domains-advanced -list
```

### 场景 2：交互式恢复所有域名
```bash
make unsuspend-all
# 或
./bin/unsuspend-domains
```
工具会列出所有挂起的域名，并询问是否继续。

### 场景 3：批量自动恢复
```bash
make unsuspend-auto
# 或
./bin/unsuspend-domains-advanced -y
```
自动恢复所有挂起的域名，不需要确认。

### 场景 4：恢复特定域名
```bash
./bin/unsuspend-domains-advanced -domain example.com
```
只恢复包含 "example.com" 的域名。

### 场景 5：测试环境批量恢复
```bash
./bin/unsuspend-domains-advanced -domain test -y
```
自动恢复所有包含 "test" 的域名。

## 功能说明

恢复操作会执行以下操作：
1. 将域名的 `status` 从 `suspended` 更新为 `active`
2. 清除 `first_failed_at` 时间戳（如果存在）
3. 更新 `updated_at` 时间戳

恢复后的域名将在下次健康检查时重新扫描。

## 查看数据库中的挂起域名

```bash
# 使用 Docker 容器
docker exec opendomain_db psql -U opendomain -d opendomain \
  -c "SELECT id, full_domain, status, first_failed_at, updated_at FROM domains WHERE status = 'suspended' ORDER BY updated_at DESC;"

# 或者直接连接数据库
psql -h localhost -p 5433 -U opendomain -d opendomain \
  -c "SELECT id, full_domain, status, first_failed_at, updated_at FROM domains WHERE status = 'suspended' ORDER BY updated_at DESC;"
```

## 注意事项

1. **备份数据**：建议在批量操作前备份数据库
2. **确认环境**：确保 `.env` 文件配置正确
3. **权限要求**：需要数据库写权限
4. **影响范围**：恢复后的域名会在下次扫描时重新检查健康状态

## 故障排查

### 问题：无法连接数据库
```bash
# 检查数据库连接
docker ps | grep postgres

# 检查 .env 配置
cat .env | grep DB_
```

### 问题：没有找到挂起的域名
```bash
# 直接查询数据库确认
docker exec opendomain_db psql -U opendomain -d opendomain \
  -c "SELECT COUNT(*) FROM domains WHERE status = 'suspended';"
```

### 问题：编译失败
```bash
# 确保依赖已下载
go mod download
go mod tidy

# 重新编译
make build-tools
```

## 相关工具

- `cmd/recover-suspended/main.go` - 旧版恢复脚本
- `internal/scanner/scanner.go` - 域名健康扫描器（自动挂起域名）
- `migrations/000017_add_first_failed_at_to_domains.up.sql` - 相关数据库迁移

## 更新日志

- **2026-02-14**: 创建域名恢复工具，支持批量和指定恢复
