-- 删除 API 配额追踪表
DROP INDEX IF EXISTS idx_api_quotas_date;
DROP INDEX IF EXISTS idx_api_quotas_api_name;
DROP TABLE IF EXISTS api_quotas;
