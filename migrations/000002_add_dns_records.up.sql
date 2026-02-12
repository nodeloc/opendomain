-- DNS Records table
CREATE TABLE IF NOT EXISTS dns_records (
  id BIGSERIAL PRIMARY KEY,
  domain_id BIGINT NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  type VARCHAR(20) NOT NULL CHECK (type IN ('A', 'AAAA', 'CNAME', 'MX', 'TXT', 'NS', 'SRV', 'CAA')),
  content TEXT NOT NULL,
  ttl INT DEFAULT 3600,
  priority INT,
  is_active BOOLEAN DEFAULT TRUE,
  synced_to_powerdns BOOLEAN DEFAULT FALSE,
  sync_error TEXT,
  last_synced_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_dns_records_domain_id ON dns_records(domain_id);
CREATE INDEX idx_dns_records_type ON dns_records(type);
CREATE INDEX idx_dns_records_is_active ON dns_records(is_active);
CREATE INDEX idx_dns_records_synced ON dns_records(synced_to_powerdns);
