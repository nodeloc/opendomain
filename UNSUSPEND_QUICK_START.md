# 快速使用指南

## 域名恢复工具

将被挂起（suspended）的域名批量恢复为活跃（active）状态。

### 快速开始

```bash
# 1. 查看所有挂起的域名
make unsuspend-list

# 2. 交互式恢复所有域名（会询问确认）
make unsuspend-all

# 3. 自动恢复所有域名（不询问）
make unsuspend-auto
```

### 高级用法

```bash
# 编译工具
make build-tools

# 查看帮助
./bin/unsuspend-domains-advanced -h

# 只恢复特定域名
./bin/unsuspend-domains-advanced -domain test.com

# 自动恢复特定域名
./bin/unsuspend-domains-advanced -domain test -y
```

### 详细文档

查看完整文档: [docs/UNSUSPEND_DOMAINS.md](docs/UNSUSPEND_DOMAINS.md)
