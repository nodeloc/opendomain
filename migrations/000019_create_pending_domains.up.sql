-- 创建待激活域名表
CREATE TABLE IF NOT EXISTS pending_domains (
    id BIGSERIAL PRIMARY KEY,
    root_domain_id BIGINT NOT NULL REFERENCES root_domains(id),
    subdomain VARCHAR(63) NOT NULL,
    full_domain VARCHAR(255) NOT NULL UNIQUE,
    fossbilling_order_id INTEGER NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    registered_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    first_failed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_pending_domains_root_domain_id ON pending_domains(root_domain_id);
CREATE INDEX idx_pending_domains_first_failed_at ON pending_domains(first_failed_at);
CREATE INDEX idx_pending_domains_deleted_at ON pending_domains(deleted_at);

-- 添加注释
COMMENT ON TABLE pending_domains IS '从FOSSBilling预同步的待激活域名';
COMMENT ON COLUMN pending_domains.status IS '状态：pending-待激活, healthy-健康, unhealthy-不健康';
COMMENT ON COLUMN pending_domains.fossbilling_order_id IS 'FOSSBilling订单ID';
COMMENT ON COLUMN pending_domains.first_failed_at IS '首次健康检查失败时间';
