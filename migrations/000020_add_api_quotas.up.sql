-- 创建 API 配额追踪表
CREATE TABLE IF NOT EXISTS api_quotas (
    id SERIAL PRIMARY KEY,
    api_name VARCHAR(100) NOT NULL,
    date VARCHAR(10) NOT NULL,
    used_count INTEGER NOT NULL DEFAULT 0,
    daily_limit INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 为 api_name 创建唯一索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_api_quotas_api_name ON api_quotas(api_name);

-- 为日期查询创建索引
CREATE INDEX IF NOT EXISTS idx_api_quotas_date ON api_quotas(date);

-- 添加注释
COMMENT ON TABLE api_quotas IS 'API 配额使用追踪';
COMMENT ON COLUMN api_quotas.api_name IS 'API 名称（google_safe_browsing 或 virustotal）';
COMMENT ON COLUMN api_quotas.date IS '日期（YYYY-MM-DD）';
COMMENT ON COLUMN api_quotas.used_count IS '当日已使用次数';
COMMENT ON COLUMN api_quotas.daily_limit IS '每日配额限制';
