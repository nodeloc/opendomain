# 快速开始指南

## 🚀 服务器部署（30分钟完成）

### 步骤 1: 环境准备

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装 Go 1.21+
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Node.js 18+
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# 安装 PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# 安装 Redis
sudo apt install -y redis-server

# 安装 Nginx
sudo apt install -y nginx

# 安装 migrate
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### 步骤 2: 克隆项目

```bash
sudo mkdir -p /var/www/opendomain
sudo chown $USER:$USER /var/www/opendomain
cd /var/www/opendomain
git clone https://github.com/your-username/opendomain.git .
```

### 步骤 3: 配置数据库

```bash
# 创建数据库
sudo -u postgres psql << SQL
CREATE DATABASE opendomain;
CREATE USER opendomain WITH PASSWORD 'your-password-here';
GRANT ALL PRIVILEGES ON DATABASE opendomain TO opendomain;
\q
SQL
```

### 步骤 4: 配置环境变量

```bash
# 复制配置文件
cp .env.example .env

# 编辑配置（修改数据库密码、JWT密钥等）
vim .env
```

**必须修改的配置**:
- `DB_PASSWORD`: 数据库密码
- `JWT_SECRET`: JWT密钥（随机字符串）
- `FRONTEND_URL`: 你的域名
- `POWERDNS_API_KEY`: PowerDNS API密钥

### 步骤 5: 一键部署

```bash
chmod +x deploy.sh
./deploy.sh
```

### 步骤 6: 配置 Systemd 服务

```bash
# 修改服务文件中的路径
sudo cp opendomain.service /etc/systemd/system/
sudo vim /etc/systemd/system/opendomain.service

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable opendomain
sudo systemctl start opendomain
sudo systemctl status opendomain
```

### 步骤 7: 配置 Nginx

```bash
# 复制配置
sudo cp nginx.conf.example /etc/nginx/sites-available/opendomain
sudo vim /etc/nginx/sites-available/opendomain  # 修改域名

# 启用配置
sudo ln -s /etc/nginx/sites-available/opendomain /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 步骤 8: 配置 SSL（可选但推荐）

```bash
# 安装 Certbot
sudo apt install -y certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# 测试自动续期
sudo certbot renew --dry-run
```

---

## ✅ 验证部署

打开浏览器访问：
- 前端: https://your-domain.com
- 健康检查: https://your-domain.com/api/health

检查服务状态：
```bash
# 查看应用日志
sudo journalctl -u opendomain -f

# 查看 Nginx 日志
sudo tail -f /var/log/nginx/opendomain_access.log

# 查看进程
ps aux | grep opendomain
```

---

## 🔄 更新部署

后续更新只需要：

```bash
cd /var/www/opendomain
./deploy.sh
```

---

## 📊 常用命令

```bash
# 重启服务
sudo systemctl restart opendomain

# 查看日志
sudo journalctl -u opendomain -n 100 -f

# 查看状态
sudo systemctl status opendomain

# 停止服务
sudo systemctl stop opendomain

# 备份数据库
pg_dump -U opendomain opendomain > backup_$(date +%Y%m%d).sql

# 查看资源占用
htop
```

---

## 🐛 常见问题

### 1. 端口 8000 被占用
```bash
sudo netstat -tulpn | grep :8000
# 修改 .env 中的 PORT
```

### 2. 数据库连接失败
```bash
sudo systemctl status postgresql
psql -h localhost -U opendomain -d opendomain
```

### 3. 前端页面空白
```bash
# 检查 Nginx 配置
sudo nginx -t
# 检查前端构建
ls -la web/dist/
```

### 4. 502 Bad Gateway
```bash
# 检查后端是否运行
sudo systemctl status opendomain
# 查看日志
sudo journalctl -u opendomain -n 50
```

---

## 📚 更多文档

- 完整部署指南: [DEPLOY.md](DEPLOY.md)
- 环境变量配置: [.env.example](.env.example)
- Nginx 配置: [nginx.conf.example](nginx.conf.example)

---

## 🆘 技术支持

如遇到问题，请：
1. 查看日志: `sudo journalctl -u opendomain -f`
2. 检查 [常见问题](DEPLOY.md#常见问题)
3. 提交 Issue: https://github.com/your-username/opendomain/issues
