# 恢复被错误暂停的域名 - 服务器部署指南

## 快速上传并执行

### 方法 1: 一键执行（推荐）

```bash
# 在本地执行，上传并在服务器上运行
scp bin/recover-suspended root@your-server:/root/
scp .env root@your-server:/root/
ssh root@your-server 'cd /root && DB_HOST=localhost DB_PORT=5432 DB_USER=opendomain DB_PASSWORD=your_password DB_NAME=opendomain ./recover-suspended'
```

### 方法 2: 使用部署脚本上传

```bash
# 1. 先编译工具
make build-tools

# 2. 上传到服务器
scp bin/recover-suspended your-server:/path/to/opendomain/
scp .env your-server:/path/to/opendomain/

# 3. SSH 到服务器执行
ssh your-server
cd /path/to/opendomain
./recover-suspended
```

### 方法 3: Git 同步后编译

```bash
# 在服务器上执行
git pull
make build-tools
./recover.sh
```

## 直接在服务器上使用（无需上传）

如果服务器上已经有源代码：

```bash
# SSH 到服务器
ssh your-server

# 进入项目目录
cd /path/to/opendomain

# 拉取最新代码
git pull

# 编译工具
make build-tools

# 执行恢复
./recover.sh
```

## 手动指定数据库

如果不使用 .env 文件：

```bash
DB_HOST=localhost \
DB_PORT=5432 \
DB_USER=opendomain \
DB_PASSWORD=your_password \
DB_NAME=opendomain \
./bin/recover-suspended
```

## 使用 Docker 环境

如果使用 Docker 部署：

```bash
# 进入 API 容器
docker exec -it opendomain-api-1 sh

# 执行恢复
cd /app
./bin/recover-suspended
```

## 注意事项

1. ✅ 确保有数据库访问权限
2. ✅ 建议先备份数据库
3. ✅ 工具只会恢复非恶意的域名
4. ✅ 真正的恶意域名（unsafe + malicious）不会被恢复

## 验证结果

恢复后，可以通过以下方式验证：

```bash
# 查看域名状态
psql -U opendomain -d opendomain -c "SELECT full_domain, status FROM domains WHERE status = 'suspended';"

# 或访问管理页面
https://your-domain.com/admin/root-domains
```
