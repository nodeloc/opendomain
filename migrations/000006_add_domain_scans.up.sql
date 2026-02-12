-- Domain scans table
CREATE TABLE IF NOT EXISTS domain_scans (
  id BIGSERIAL PRIMARY KEY,
  domain_id BIGINT NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
  scan_type VARCHAR(50) NOT NULL, -- http, dns, ssl
  status VARCHAR(20) NOT NULL, -- success, failed, timeout
  response_time INT, -- milliseconds
  http_status_code INT,
  ssl_valid BOOLEAN,
  ssl_expiry_date TIMESTAMP,
  error_message TEXT,
  scan_details JSONB,
  scanned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Domain scan summary (latest scan results for each domain)
CREATE TABLE IF NOT EXISTS domain_scan_summaries (
  domain_id BIGINT PRIMARY KEY REFERENCES domains(id) ON DELETE CASCADE,
  last_scanned_at TIMESTAMP,
  http_status VARCHAR(20), -- online, offline, error
  dns_status VARCHAR(20), -- resolved, failed
  ssl_status VARCHAR(20), -- valid, invalid, none
  overall_health VARCHAR(20), -- healthy, degraded, down
  total_scans INT DEFAULT 0,
  successful_scans INT DEFAULT 0,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_domain_scans_domain_id ON domain_scans(domain_id, scanned_at DESC);
CREATE INDEX idx_domain_scans_status ON domain_scans(status);
CREATE INDEX idx_domain_scans_type ON domain_scans(scan_type);
CREATE INDEX idx_domain_scan_summaries_health ON domain_scan_summaries(overall_health);
