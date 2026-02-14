# 恢复被错误暂停的域名

此工具用于恢复因 Safe Browsing API 调用失败而被错误暂停的域名。

## 编译

```bash
go build -o recover-suspended cmd/recover-suspended/main.go
```

## 使用方法

### 方式1：使用环境变量

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=opendomain

./recover-suspended
```

### 方式2：使用 .env 文件

如果你的项目使用 .env 文件，可以直接运行：

```bash
# 确保 .env 文件中有数据库配置
./recover-suspended
```

### 方式3：一行命令

```bash
DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=your_password DB_NAME=opendomain ./recover-suspended
```

## 恢复逻辑

该工具会：

1. ✅ **自动恢复**以下情况的域名：
   - Safe Browsing 状态为 `safe` 或 `unknown`，且 VirusTotal 不是 `malicious`
   - 域名整体健康状态为 `healthy`

2. ⚠️ **保持暂停**以下情况的域名：
   - Safe Browsing 状态为 `unsafe` 且 VirusTotal 为 `malicious`（真正的恶意域名）

3. ❌ **跳过**没有扫描记录的域名

## 示例输出

```
=== 恢复被错误暂停的域名 ===
正在查找被暂停的域名...
找到 5 个被暂停的域名

✅ dsm.loc.cc - 已恢复 (Safe Browsing: unknown, VirusTotal: unknown (非恶意))
✅ test.loc.cc - 已恢复 (域名健康状态正常)
⚠️  malicious.loc.cc - 保持暂停状态 (Safe Browsing: unsafe, VirusTotal: malicious)
✅ example.loc.cc - 已恢复 (Safe Browsing: safe, VirusTotal: clean (非恶意))
❌ orphan.loc.cc - 未找到扫描记录，跳过

=== 完成 ===
恢复: 3 个
跳过: 2 个
总计: 5 个
```

## 注意事项

- 该工具只会恢复状态为 `suspended` 的域名
- 恢复时会将域名状态改为 `active` 并清除 `first_failed_at` 字段
- 真正的恶意域名（Safe Browsing: unsafe + VirusTotal: malicious）不会被恢复
- 建议在运行前备份数据库
