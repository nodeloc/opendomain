# Docker 部署指南

使用 Docker 和 Docker Compose 快速部署 OpenDomain。

## 🚀 快速开始（5分钟）

### 前置要求

- **Docker**: 20.10+ ([安装指南](https://docs.docker.com/get-docker/))
- **Docker Compose**: 2.0+ ([安装指南](https://docs.docker.com/compose/install/))
- **Git**: 用于克隆代码

### 一键部署

```bash
# 1. 克隆代码
git clone https://github.com/your-username/opendomain.git
cd opendomain

# 2. 配置环境变量
cp .env.example .env
vim .env  # 修改必要的配置

# 3. 一键部署
chmod +x docker-deploy.sh
./docker-deploy.sh
```

部署完成后访问：
- 前端: http://localhost:8000
- 健康检查: http://localhost:8000/health

---

## 📋 详细步骤

### 1. 安装 Docker

**Ubuntu/Debian:**
```bash
# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 添加当前用户到 docker 组
sudo usermod -aG docker $USER
newgrp docker

# 验证安装
docker --version
docker compose version
```

**CentOS/RHEL:**
```bash
sudo yum install -y yum-utils
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
sudo yum install docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo systemctl start docker
sudo systemctl enable docker
```

### 2. 配置环境变量

```bash
# 复制模板
cp .env.example .env

# 编辑配置
vim .env
```

**必须配置的变量**:

```bash
# 数据库密码（请修改为强密码）
DB_PASSWORD=your-secure-password

# JWT 密钥（随机字符串，至少32位）
JWT_SECRET=your-super-secret-jwt-key-at-least-32-chars

# PowerDNS API 配置
POWERDNS_API_URL=http://host.docker.internal:8081
POWERDNS_API_KEY=your-powerdns-api-key

# 前端 URL（生产环境域名）
FRONTEND_URL=https://your-domain.com

# 默认 NS 服务器
DEFAULT_NS1=ns1.your-domain.com
DEFAULT_NS2=ns2.your-domain.com
```

### 3. 构建和启动

#### 方式 1: 使用部署脚本（推荐）

```bash
chmod +x docker-deploy.sh
./docker-deploy.sh
```

#### 方式 2: 使用 Docker Compose

```bash
# 构建镜像
docker compose build

# 启动所有服务
docker compose up -d

# 查看状态
docker compose ps

# 查看日志
docker compose logs -f
```

---

## 🔧 服务组件

Docker Compose 会启动以下服务：

| 服务 | 说明 | 端口 |
|------|------|------|
| **app** | OpenDomain 应用 | 8000 |
| **postgres** | PostgreSQL 数据库 | 5432 |
| **redis** | Redis 缓存 | 6379 |
| **migrate** | 数据库迁移（一次性） | - |
| **nginx** | Nginx 反向代理（可选） | 80, 443 |

---

## 📊 常用命令

### 查看服务状态

```bash
# 查看所有服务
docker compose ps

# 查看服务日志
docker compose logs -f

# 查看指定服务日志
docker compose logs -f app
docker compose logs -f postgres
```

### 启动和停止

```bash
# 启动所有服务
docker compose up -d

# 停止所有服务
docker compose down

# 重启服务
docker compose restart

# 重启指定服务
docker compose restart app
```

### 进入容器

```bash
# 进入应用容器
docker compose exec app sh

# 进入数据库容器
docker compose exec postgres psql -U opendomain

# 查看数据库
docker compose exec postgres psql -U opendomain -d opendomain -c "SELECT * FROM users;"
```

### 更新部署

```bash
# 方式 1: 使用部署脚本
./docker-deploy.sh --build

# 方式 2: 手动更新
git pull
docker compose build --no-cache
docker compose down
docker compose up -d
```

### 数据备份

```bash
# 备份数据库
docker compose exec postgres pg_dump -U opendomain opendomain > backup_$(date +%Y%m%d).sql

# 恢复数据库
cat backup_20240101.sql | docker compose exec -T postgres psql -U opendomain opendomain
```

### 清理数据

```bash
# 停止并删除容器
docker compose down

# 停止并删除容器和数据卷
docker compose down -v

# 清理所有未使用的镜像
docker system prune -a
```

---

## 🌐 生产环境配置

### 使用 Nginx 反向代理

启用 Nginx 服务：

```bash
# 启动时包含 nginx
docker compose --profile with-nginx up -d

# 配置 Nginx
vim nginx/conf.d/opendomain.conf
docker compose restart nginx
```

### 配置 SSL 证书

**方式 1: 使用 Let's Encrypt (推荐)**

```bash
# 安装 Certbot
sudo apt install certbot

# 获取证书
sudo certbot certonly --standalone -d your-domain.com

# 复制证书到项目
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem nginx/ssl/
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem nginx/ssl/

# 重启 nginx
docker compose restart nginx
```

**方式 2: 手动配置证书**

将证书文件放到 `nginx/ssl/` 目录：
```
nginx/ssl/
├── fullchain.pem
└── privkey.pem
```

修改 `nginx/conf.d/opendomain.conf` 添加 SSL 配置：
```nginx
server {
    listen 443 ssl http2;
    server_name your-domain.com;
    
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;
    
    # ... 其他配置
}
```

---

## 🔍 故障排查

### 容器无法启动

```bash
# 查看容器日志
docker compose logs app

# 查看所有服务状态
docker compose ps -a

# 检查配置文件
docker compose config
```

### 数据库连接失败

```bash
# 检查数据库容器是否运行
docker compose ps postgres

# 测试数据库连接
docker compose exec postgres psql -U opendomain -d opendomain -c "SELECT 1;"

# 查看数据库日志
docker compose logs postgres
```

### 前端无法访问

```bash
# 检查应用容器端口映射
docker compose port app 8000

# 检查防火墙
sudo ufw status
sudo ufw allow 8000/tcp

# 查看应用日志
docker compose logs -f app
```

### 磁盘空间不足

```bash
# 查看 Docker 磁盘使用
docker system df

# 清理未使用的镜像
docker image prune -a

# 清理未使用的容器
docker container prune

# 清理未使用的卷
docker volume prune
```

---

## 📈 性能优化

### 1. 调整资源限制

修改 `docker-compose.yml`：

```yaml
services:
  app:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

### 2. 使用生产模式

确保 `.env` 中设置：
```bash
APP_ENV=production
LOG_LEVEL=info
```

### 3. 数据库优化

在 `docker-compose.yml` 中添加 PostgreSQL 参数：

```yaml
postgres:
  command: postgres -c max_connections=200 -c shared_buffers=256MB -c effective_cache_size=1GB
```

### 4. Redis 持久化

已启用 AOF 持久化：
```yaml
redis:
  command: redis-server --appendonly yes
```

---

## 🔒 安全建议

1. **修改默认密码**
   - 数据库密码
   - Redis 密码（如果启用）
   - JWT 密钥

2. **限制端口访问**
   ```bash
   # 只暴露必要的端口
   # 在 docker-compose.yml 中移除不必要的端口映射
   ```

3. **使用非 root 用户**
   - Dockerfile 已配置为使用 appuser

4. **定期更新镜像**
   ```bash
   docker compose pull
   docker compose up -d
   ```

5. **配置防火墙**
   ```bash
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw deny 5432/tcp  # 不要暴露数据库端口
   sudo ufw deny 6379/tcp  # 不要暴露 Redis 端口
   ```

---

## 🎯 高可用部署

### Docker Swarm 集群

```bash
# 初始化 Swarm
docker swarm init

# 部署 Stack
docker stack deploy -c docker-compose.yml opendomain

# 查看服务
docker stack services opendomain

# 扩容应用
docker service scale opendomain_app=3
```

### Kubernetes 部署

参考 `k8s/` 目录（需要单独创建）

---

## 📚 相关文档

- [Docker 官方文档](https://docs.docker.com/)
- [Docker Compose 文档](https://docs.docker.com/compose/)
- [OpenDomain 部署指南](DEPLOY.md)
- [快速开始](QUICKSTART.md)

---

## 🆘 获取帮助

- **查看日志**: `docker compose logs -f app`
- **健康检查**: `curl http://localhost:8000/health`
- **提交 Issue**: https://github.com/your-username/opendomain/issues
