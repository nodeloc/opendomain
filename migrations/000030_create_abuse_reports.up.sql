-- Abuse reports submitted from the public report form
CREATE TABLE IF NOT EXISTS abuse_reports (
    id BIGSERIAL PRIMARY KEY,
    subdomain VARCHAR(255) NOT NULL,
    category VARCHAR(20) NOT NULL CHECK (category IN ('phishing', 'malware', 'scam', 'spam', 'impersonation', 'other')),
    url VARCHAR(2048) NOT NULL DEFAULT '',
    details TEXT NOT NULL,
    reporter_email VARCHAR(255) NOT NULL DEFAULT '',
    reporter_ip VARCHAR(64) NOT NULL DEFAULT '',
    user_agent VARCHAR(512) NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'reviewing', 'resolved', 'dismissed')),
    admin_note TEXT NOT NULL DEFAULT '',
    domain_id BIGINT REFERENCES domains(id) ON DELETE SET NULL,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_abuse_reports_status_created ON abuse_reports(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_abuse_reports_created_at ON abuse_reports(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_abuse_reports_subdomain ON abuse_reports(subdomain);
CREATE INDEX IF NOT EXISTS idx_abuse_reports_domain_id ON abuse_reports(domain_id);
CREATE INDEX IF NOT EXISTS idx_abuse_reports_reporter_ip_created ON abuse_reports(reporter_ip, created_at DESC);
