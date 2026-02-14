-- 恢复被错误暂停的域名脚本
-- 这个脚本会找出因 Safe Browsing API 失败而被误判暂停的域名

-- 1. 查看所有被暂停的域名及其扫描结果
SELECT 
    d.id,
    d.full_domain,
    d.status,
    d.first_failed_at,
    dss.safe_browsing_status,
    dss.virustotal_status,
    dss.overall_health
FROM domains d
LEFT JOIN domain_scan_summaries dss ON d.id = dss.domain_id
WHERE d.status = 'suspended'
ORDER BY d.id;

-- 2. 查看这些域名最近的 Safe Browsing 扫描详情
SELECT 
    ds.domain_id,
    d.full_domain,
    ds.scan_type,
    ds.status,
    ds.error_message,
    ds.scan_details,
    ds.scanned_at
FROM domain_scans ds
JOIN domains d ON ds.domain_id = d.id
WHERE d.status = 'suspended' 
  AND ds.scan_type = 'safebrowsing'
  AND ds.scanned_at > NOW() - INTERVAL '7 days'
ORDER BY ds.domain_id, ds.scanned_at DESC;

-- 3. 恢复那些因 Safe Browsing API 调用失败而被误判的域名
-- （status = 'failed' 且 error_message 包含错误信息，而不是真正的威胁检测）
-- 注意：请先检查上面的查询结果，确认哪些域名应该恢复

-- 如果确定要恢复，取消下面这些语句的注释：

-- 恢复所有因 API 失败（非威胁检测）被暂停的域名
-- UPDATE domains
-- SET status = 'active', first_failed_at = NULL
-- WHERE id IN (
--     SELECT DISTINCT d.id
--     FROM domains d
--     JOIN domain_scans ds ON d.id = ds.domain_id
--     WHERE d.status = 'suspended'
--       AND ds.scan_type = 'safebrowsing'
--       AND ds.status = 'failed'
--       AND ds.error_message IS NOT NULL
--       AND ds.error_message != ''
--       AND ds.scanned_at > NOW() - INTERVAL '1 day'
-- );

-- 或者手动恢复特定域名（替换 YOUR_DOMAIN_ID）：
-- UPDATE domains SET status = 'active', first_failed_at = NULL WHERE id = YOUR_DOMAIN_ID;

-- 4. 验证恢复结果
-- SELECT id, full_domain, status, first_failed_at 
-- FROM domains 
-- WHERE full_domain IN ('dsm.loc.cc', '...其他域名...');
