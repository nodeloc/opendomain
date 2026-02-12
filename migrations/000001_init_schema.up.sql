-- Users table
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(100) NOT NULL UNIQUE,
  email_verified BOOLEAN DEFAULT FALSE,
  phone VARCHAR(20),
  phone_verified BOOLEAN DEFAULT FALSE,
  password_hash VARCHAR(255) NOT NULL,
  avatar VARCHAR(255),
  real_name VARCHAR(50),
  is_verified BOOLEAN DEFAULT FALSE,
  user_level VARCHAR(20) DEFAULT 'normal' CHECK (user_level IN ('normal', 'verified', 'vip', 'whitelist')),
  domain_quota INT DEFAULT 2,
  invite_code VARCHAR(20) NOT NULL UNIQUE,
  invited_by BIGINT,
  status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'frozen', 'banned')),
  last_login_at TIMESTAMP,
  last_login_ip VARCHAR(45),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_invite_code ON users(invite_code);
CREATE INDEX idx_users_invited_by ON users(invited_by);
CREATE INDEX idx_users_status ON users(status);

-- Root domains table
CREATE TABLE IF NOT EXISTS root_domains (
  id SERIAL PRIMARY KEY,
  domain VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  nameservers TEXT NOT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  is_hot BOOLEAN DEFAULT FALSE,
  is_new BOOLEAN DEFAULT FALSE,
  priority INT DEFAULT 0,
  min_length INT DEFAULT 3,
  max_length INT DEFAULT 63,
  registration_count INT DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_root_domains_is_active ON root_domains(is_active);
CREATE INDEX idx_root_domains_priority ON root_domains(priority);

-- Domains table
CREATE TABLE IF NOT EXISTS domains (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  root_domain_id INT NOT NULL REFERENCES root_domains(id) ON DELETE RESTRICT,
  subdomain VARCHAR(63) NOT NULL,
  full_domain VARCHAR(255) NOT NULL UNIQUE,
  status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'suspended', 'deleted')),
  registered_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP NOT NULL,
  auto_renew BOOLEAN DEFAULT FALSE,
  reminder_sent_30d BOOLEAN DEFAULT FALSE,
  reminder_sent_7d BOOLEAN DEFAULT FALSE,
  dns_synced BOOLEAN DEFAULT FALSE,
  dns_sync_error TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_domains_user_id ON domains(user_id);
CREATE INDEX idx_domains_root_domain_id ON domains(root_domain_id);
CREATE INDEX idx_domains_status ON domains(status);
CREATE INDEX idx_domains_expires_at ON domains(expires_at);
CREATE INDEX idx_domains_full_domain ON domains(full_domain);
CREATE INDEX idx_domains_auto_renew ON domains(auto_renew);
CREATE INDEX idx_domains_dns_synced ON domains(dns_synced);
