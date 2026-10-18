-- Sponsors shown on the public sponsor page
CREATE TABLE IF NOT EXISTS sponsors (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    website_url VARCHAR(500) NOT NULL DEFAULT '',
    logo_url VARCHAR(500) NOT NULL DEFAULT '',
    logo_dark_url VARCHAR(500) NOT NULL DEFAULT '',
    description_zh VARCHAR(300) NOT NULL DEFAULT '',
    description_en VARCHAR(300) NOT NULL DEFAULT '',
    tier VARCHAR(20) NOT NULL DEFAULT 'sponsor' CHECK (tier IN ('partner', 'sponsor', 'backer')),
    display_order INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sponsors_active_order ON sponsors(is_active, display_order);

INSERT INTO system_settings (setting_key, setting_value, description, created_at, updated_at)
VALUES
  ('sponsor_usdt_trc20_address', '', 'USDT (TRC20) address shown on the sponsor page; empty hides it', NOW(), NOW()),
  ('sponsor_github_url', '', 'GitHub Sponsors URL shown on the sponsor page; empty hides it', NOW(), NOW())
ON CONFLICT (setting_key) DO NOTHING;
