# OpenDomain 部署指南

## 前置要求

### 服务器环境
- **操作系统**: Linux (Ubuntu 20.04+ 或 CentOS 7+)
- **CPU**: 2核以上
- **内存**: 2GB 以上
- **磁盘**: 20GB 以上

### 软件依赖
- **Go**: 1.21+ ([安装指南](https://golang.org/doc/install))
- **Node.js**: 18+ ([安装指南](https://nodejs.org/))
- **PostgreSQL**: 13+ ([安装指南](https://www.postgresql.org/download/))
- **Redis**: 6+ ([安装指南](https://redis.io/download))
- **PowerDNS**: 4.5+ ([安装指南](https://doc.powerdns.com/))
- **migrate**: 用于数据库迁移 ([安装指南](https://github.com/golang-migrate/migrate))

---

## 快速部署

### 1. 克隆代码

```bash
# 创建项目目录
sudo mkdir -p /var/www/opendomain
sudo chown $USER:$USER /var/www/opendomain

# 克隆代码
cd /var/www/opendomain
git clone https://github.com/your-username/opendomain.git .
```

### 2. 配置环境变量

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置文件
vim .env
```

**重要配置项**:
```bash
# 应用配置
APP_ENV=production
PORT=8000
FRONTEND_URL=https://your-domain.com

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=opendomain
DB_PASSWORD=your-secure-password
DB_NAME=opendomain

# JWT 密钥（请修改为随机字符串）
JWT_SECRET=your-super-secret-jwt-key-change-this

# PowerDNS API
POWERDNS_API_URL=http://localhost:8081
POWERDNS_API_KEY=your-powerdns-api-key

# 邮件配置（可选）
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USER=your-email@gmail.com
MAIL_PASSWORD=your-app-password

# Redis（可选）
REDIS_HOST=localhost
REDIS_PORT=6379
```

### 3. 创建数据库

```bash
# 登录 PostgreSQL
sudo -u postgres psql

# 创建数据库和用户
CREATE DATABASE opendomain;
CREATE USER opendomain WITH PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE opendomain TO opendomain;
\q
```

### 4. 安装 migrate 工具

```bash
# Linux x64
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
sudo chmod +x /usr/local/bin/migrate
```

### 5. 运行部署脚本

```bash
chmod +x deploy.sh
./deploy.sh
```

部署脚本会自动完成：
1. 拉取最新代码
2. 安装 Go 依赖
3. 编译后端程序
4. 运行数据库迁移
5. 安装前端依赖
6. 构建前端应用
7. 重启服务（如果已运行）

---

## Systemd 服务配置（推荐）

### 1. 创建服务文件

```bash
sudo cp opendomain.service /etc/systemd/system/
sudo vim /etc/systemd/system/opendomain.service
```

修改以下配置项：
- `User`: 运行服务的用户（建议使用 www-data 或创建专用用户）
- `Group`: 运行服务的用户组
- `WorkingDirectory`: 项目目录路径
- `EnvironmentFile`: .env 文件路径
- `ExecStart`: opendomain 可执行文件路径

### 2. 创建日志目录

```bash
mkdir -p logs
```

### 3. 启用并启动服务

```bash
# 重新加载 systemd 配置
sudo systemctl daemon-reload

# 启用开机自启
sudo systemctl enable opendomain

# 启动服务
sudo systemctl start opendomain

# 查看服务状态
sudo systemctl status opendomain

# 查看日志
sudo journalctl -u opendomain -f
```

---

## Nginx 反向代理配置

### 1. 安装 Nginx

```bash
sudo apt install nginx  # Ubuntu/Debian
sudo yum install nginx  # CentOS
```

### 2. 创建配置文件

```bash
sudo vim /etc/nginx/sites-available/opendomain
```

**配置示例**:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 证书配置
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # 前端静态文件
    location / {
        root /var/www/opendomain/web/dist;
        try_files $uri $uri/ /index.html;
        
        # 缓存静态资源
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # API 反向代理
    location /api/ {
        proxy_pass http://localhost:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 文件上传大小限制
    client_max_body_size 10M;
}
```

### 3. 启用配置

```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/opendomain /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
```

### 4. 配置 SSL（Let's Encrypt）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期测试
sudo certbot renew --dry-run
```

---

## 手动启动（不使用 Systemd）

如果不想使用 systemd，可以手动启动：

```bash
# 后台运行
nohup ./opendomain > logs/app.log 2>&1 &

# 查看日志
tail -f logs/app.log

# 停止服务
pkill -f opendomain
```

---

## 更新部署

后续更新只需运行：

```bash
cd /var/www/opendomain
./deploy.sh
```

---

## 常见问题

### 1. 数据库连接失败

检查 PostgreSQL 是否运行：
```bash
sudo systemctl status postgresql
```

检查数据库配置：
```bash
psql -h localhost -U opendomain -d opendomain
```

### 2. 端口被占用

查看端口占用：
```bash
sudo netstat -tulpn | grep :8000
```

修改 `.env` 中的 `PORT` 配置。

### 3. 权限问题

确保运行用户对项目目录有读写权限：
```bash
sudo chown -R www-data:www-data /var/www/opendomain
sudo chmod -R 755 /var/www/opendomain
```

### 4. 前端构建失败

清除 node_modules 重新安装：
```bash
cd web
rm -rf node_modules package-lock.json
npm install
npm run build
```

### 5. 数据库迁移失败

手动运行迁移：
```bash
export $(cat .env | grep -v '^#' | xargs)
migrate -path ./migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" up
```

---

## 监控和维护

### 查看系统日志

```bash
# systemd 日志
sudo journalctl -u opendomain -f

# 应用日志
tail -f logs/app.log

# Nginx 日志
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

### 数据库备份

```bash
# 备份数据库
pg_dump -U opendomain opendomain > backup_$(date +%Y%m%d).sql

# 恢复数据库
psql -U opendomain opendomain < backup_20240101.sql
```

### 性能监控

```bash
# 查看进程资源占用
htop

# 查看连接数
sudo netstat -ant | grep :8000 | wc -l

# 查看数据库连接
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity;"
```

---

## 安全建议

1. **防火墙配置**
   ```bash
   sudo ufw allow 22/tcp      # SSH
   sudo ufw allow 80/tcp      # HTTP
   sudo ufw allow 443/tcp     # HTTPS
   sudo ufw enable
   ```

2. **定期更新**
   ```bash
   sudo apt update && sudo apt upgrade
   ```

3. **最小权限原则**
   - 不要使用 root 用户运行应用
   - 限制文件权限

4. **备份策略**
   - 定期备份数据库
   - 定期备份配置文件
   - 定期备份用户上传的文件

5. **监控告警**
   - 配置监控系统（如 Prometheus + Grafana）
   - 配置日志分析（如 ELK Stack）
   - 配置告警通知

---

## 技术支持

- **文档**: https://github.com/your-username/opendomain
- **问题反馈**: https://github.com/your-username/opendomain/issues
