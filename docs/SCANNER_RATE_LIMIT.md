# 扫描器速率限制优化

## 🎯 优化内容

### 1. API 速率限制
- **Google Safe Browsing**: 1 请求/秒，10,000 请求/天
- **VirusTotal**: 4 请求/分钟（15秒间隔），500 请求/天

### 2. 分批扫描
- 每批处理 50 个域名
- 批次之间等待 5 秒
- 整个流程完成后才开始下一轮

### 3. 配额持久化
- 新增 `api_quotas` 表存储配额使用情况
- 服务重启后配额计数不丢失
- 每次 API 调用后自动保存到数据库

### 4. 配额管理
- 自动检测每日配额重置
- 超出配额时自动跳过，不影响其他扫描
- 提供 API 查询配额使用状态

## 📊 新增表结构

```sql
CREATE TABLE api_quotas (
    id SERIAL PRIMARY KEY,
    api_name VARCHAR(100) NOT NULL UNIQUE,
    date VARCHAR(10) NOT NULL,
    used_count INTEGER NOT NULL DEFAULT 0,
    daily_limit INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
```

## 🚀 使用方法

### 运行数据库迁移

```bash
make migrate-up
```

### 查看配额状态

**API 端点**:
```
GET /api/admin/api-quota
```

**响应示例**:
```json
{
  "google_safe_browsing": {
    "used": 150,
    "limit": 10000,
    "remaining": 9850,
    "date": "2026-02-14"
  },
  "virustotal": {
    "used": 45,
    "limit": 500,
    "remaining": 455,
    "date": "2026-02-14"
  }
}
```

## 📝 扫描日志

扫描时会输出详细日志：

```
[INFO] Starting scan for 1234 domains
[INFO] Google Safe Browsing quota used: 150/10000
[INFO] VirusTotal quota used: 45/500
[INFO] Processing batch 1-50 of 1234 domains
[INFO] Batch complete. Waiting 5 seconds before next batch...
[INFO] Processing batch 51-100 of 1234 domains
...
[INFO] Scan complete. GSB used: 1384/10000, VT used: 545/500
```

## ⚠️ 配额超限处理

### Google Safe Browsing
- 超出配额时，状态标记为 `quota_exceeded`
- Safe Browsing 状态设为 `unknown`（不影响域名）
- 不会触发自动暂停

### VirusTotal
- 超出配额时，状态标记为 `quota_exceeded`
- VirusTotal 状态设为 `unknown`（不影响域名）
- 不会触发自动暂停

## 🔧 配置建议

### 大量域名场景
如果域名数量超过配额限制，建议：

1. **分时段扫描**: 每天只扫描部分域名
2. **优先级排序**: 优先扫描重要域名
3. **升级 API**: 考虑使用付费 API 提升配额

### 示例：分时段扫描
```go
// 每天只扫描 1/7 的域名
dayOfWeek := time.Now().Weekday()
domains := getAllDomains()
batchSize := len(domains) / 7
start := int(dayOfWeek) * batchSize
end := start + batchSize
todayDomains := domains[start:end]
```

## 📈 监控建议

定期检查配额使用情况：
- 每天早上查看配额重置
- 监控接近限制的情况
- 及时调整扫描策略

## 🐛 故障恢复

如果配额计数异常：

```sql
-- 查看当前配额
SELECT * FROM api_quotas WHERE date = CURRENT_DATE;

-- 重置配额（谨慎操作）
UPDATE api_quotas SET used_count = 0 WHERE api_name = 'google_safe_browsing';
UPDATE api_quotas SET used_count = 0 WHERE api_name = 'virustotal';
```
