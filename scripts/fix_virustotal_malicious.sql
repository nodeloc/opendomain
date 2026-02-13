-- ========================================
-- 修复因 VirusTotal API 失败而被错误标记为 Malicious 和 Suspended 的域名
-- ========================================

BEGIN;

-- 创建临时表存储需要修复的域名
CREATE TEMP TABLE temp_fix_domains AS
SELECT DISTINCT
    dss.domain_id,
    d.full_domain,
    d.status as domain_status,
    dss.virustotal_status,
    ds.status as vt_scan_status,
    ds.error_message as vt_error,
    ds.scan_details
FROM domain_scan_summaries dss
JOIN domains d ON d.id = dss.domain_id
LEFT JOIN LATERAL (
    SELECT status, error_message, scan_details, scanned_at
    FROM domain_scans
    WHERE domain_id = dss.domain_id
      AND scan_type = 'virustotal'
    ORDER BY scanned_at DESC
    LIMIT 1
) ds ON true
WHERE dss.virustotal_status = 'malicious'
  AND (ds.status IN ('failed', 'quota_exceeded') OR ds.status IS NULL);

-- 显示将要修复的域名
SELECT
    '准备修复的域名:' as info,
    COUNT(*) as total_count
FROM temp_fix_domains;

SELECT
    full_domain,
    domain_status,
    vt_scan_status,
    vt_error
FROM temp_fix_domains
ORDER BY full_domain
LIMIT 20;

-- 1. 更新 domain_scan_summaries: 将错误的 malicious 改为 unknown
UPDATE domain_scan_summaries dss
SET
    virustotal_status = 'unknown',
    updated_at = NOW()
FROM temp_fix_domains tfd
WHERE dss.domain_id = tfd.domain_id;

SELECT
    '已更新 domain_scan_summaries' as step,
    COUNT(*) as affected_rows
FROM temp_fix_domains;

-- 2. 重新计算 overall_health (排除 VirusTotal malicious 影响)
UPDATE domain_scan_summaries dss
SET
    overall_health = CASE
        -- 如果 SafeBrowsing 不安全，仍然是 degraded
        WHEN dss.safe_browsing_status = 'unsafe' THEN 'degraded'
        -- DNS 和 HTTP 都正常，且没有安全问题
        WHEN dss.dns_status = 'resolved' AND dss.http_status = 'online' THEN 'healthy'
        -- 至少一个正常
        WHEN dss.dns_status = 'resolved' OR dss.http_status = 'online' THEN 'degraded'
        -- 都失败
        ELSE 'down'
    END,
    updated_at = NOW()
FROM temp_fix_domains tfd
WHERE dss.domain_id = tfd.domain_id;

SELECT
    '已重新计算 overall_health' as step,
    COUNT(*) as affected_rows
FROM temp_fix_domains;

-- 3. 恢复被错误 suspend 的域名
-- 只恢复那些实际上是健康的域名（DNS 或 HTTP 正常）
UPDATE domains d
SET
    status = 'active',
    first_failed_at = NULL
FROM temp_fix_domains tfd
JOIN domain_scan_summaries dss ON dss.domain_id = tfd.domain_id
WHERE d.id = tfd.domain_id
  AND d.status = 'suspended'
  AND (dss.dns_status = 'resolved' OR dss.http_status = 'online')
  AND dss.safe_browsing_status != 'unsafe';

SELECT
    '已恢复域名状态 (suspended -> active)' as step,
    COUNT(*) as affected_rows
FROM temp_fix_domains tfd
JOIN domains d ON d.id = tfd.domain_id
WHERE d.status = 'active';

-- 显示修复后的统计
SELECT
    '===== 修复完成统计 =====' as info;

SELECT
    'VirusTotal 状态分布' as category,
    virustotal_status,
    COUNT(*) as count
FROM domain_scan_summaries
WHERE domain_id IN (SELECT domain_id FROM temp_fix_domains)
GROUP BY virustotal_status;

SELECT
    '整体健康状态分布' as category,
    overall_health,
    COUNT(*) as count
FROM domain_scan_summaries
WHERE domain_id IN (SELECT domain_id FROM temp_fix_domains)
GROUP BY overall_health;

SELECT
    '域名状态分布' as category,
    status,
    COUNT(*) as count
FROM domains
WHERE id IN (SELECT domain_id FROM temp_fix_domains)
GROUP BY status;

-- 如果一切正常，提交事务
-- 如果需要回滚，执行 ROLLBACK; 而不是 COMMIT;
COMMIT;

-- 如果想要先预览而不实际执行，将上面的 COMMIT 改为：
-- ROLLBACK;
