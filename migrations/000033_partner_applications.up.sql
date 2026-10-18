-- Partnership applications submitted from the public sponsor page
CREATE TABLE IF NOT EXISTS partner_applications (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(254) NOT NULL,
    organization VARCHAR(200) NOT NULL DEFAULT '',
    website VARCHAR(500) NOT NULL DEFAULT '',
    partnership_type VARCHAR(20) NOT NULL CHECK (partnership_type IN ('sponsorship', 'resources', 'business', 'community', 'other')),
    message TEXT NOT NULL,
    contact VARCHAR(200) NOT NULL DEFAULT '',
    status VARCHAR(20) NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'contacted', 'accepted', 'declined')),
    admin_note TEXT NOT NULL DEFAULT '',
    ip VARCHAR(64) NOT NULL DEFAULT '',
    user_agent VARCHAR(512) NOT NULL DEFAULT '',
    locale VARCHAR(10) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_partner_applications_status_created ON partner_applications(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_partner_applications_created_at ON partner_applications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_partner_applications_ip_created ON partner_applications(ip, created_at DESC);
