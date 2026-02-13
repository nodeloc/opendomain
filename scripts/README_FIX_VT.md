# VirusTotal Malicious 状态修复脚本

## 问题描述
当 VirusTotal API 超出配额时，之前的代码错误地将所有域名标记为 `malicious` 并自动 `suspend`。

## 修复方案

### 方式 1: SQL 脚本 (推荐 - 快速)

```bash
# 连接数据库并执行修复脚本
psql $DATABASE_URL -f scripts/fix_virustotal_malicious.sql
```

或者使用 Docker:
```bash
docker exec -i opendomain-db psql -U yourusername -d opendomain < scripts/fix_virustotal_malicious.sql
```

### 方式 2: Go 脚本 (备选)

```bash
cd scripts
go run fix_virustotal_malicious.go
```

## SQL 脚本功能

1. ✅ 识别所有被错误标记为 `malicious` 的域名
   - 检查最新的 VirusTotal 扫描记录
   - 只修复扫描状态为 `failed` 或 `quota_exceeded` 的记录

2. ✅ 更新 `domain_scan_summaries` 表
   - 将 `virus_total_status` 从 `malicious` 改为 `unknown`
   - 重新计算 `overall_health` 状态

3. ✅ 恢复域名状态
   - 将符合条件的域名从 `suspended` 改回 `active`
   - 清除 `first_failed_at` 标记
   - 仅恢复实际健康的域名（DNS 或 HTTP 正常）

4. ✅ 显示详细的修复统计

## 安全特性

- 使用事务 (BEGIN/COMMIT)，确保原子性
- 如果需要先预览，可以将最后的 `COMMIT` 改为 `ROLLBACK`
- 临时表不会影响现有数据
- 显示修复前后的统计对比

## 预期输出示例

```
准备修复的域名: 245
已更新 domain_scan_summaries: 245
已重新计算 overall_health: 245
已恢复域名状态: 198

VirusTotal 状态分布:
- unknown: 245

整体健康状态分布:
- healthy: 180
- degraded: 45
- down: 20

域名状态分布:
- active: 198
- suspended: 47
```

## 注意事项

1. 执行前建议先备份数据库
2. 脚本会自动识别真正恶意的域名（实际检测到威胁的）并保留其状态
3. 只恢复那些实际上健康的域名，避免错误恢复真正有问题的域名
4. 脚本是幂等的，可以安全地多次执行

## 验证修复

执行后可以运行以下查询验证：

```sql
-- 检查还有多少 malicious 状态（应该只剩下真正恶意的）
SELECT COUNT(*) FROM domain_scan_summaries
WHERE virus_total_status = 'malicious';

-- 检查被恢复的域名
SELECT full_domain, status, virus_total_status, overall_health
FROM domains d
JOIN domain_scan_summaries dss ON d.id = dss.domain_id
WHERE d.status = 'active'
  AND dss.virus_total_status = 'unknown'
LIMIT 10;
```
