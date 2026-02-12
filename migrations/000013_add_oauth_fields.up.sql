-- Add OAuth support to users table
ALTER TABLE users ADD COLUMN provider VARCHAR(20) DEFAULT 'local';
ALTER TABLE users ADD COLUMN oauth_id VARCHAR(255);

-- Allow password_hash to be NULL for OAuth users
ALTER TABLE users ALTER COLUMN password_hash DROP NOT NULL;

-- Unique index: one OAuth account per provider
CREATE UNIQUE INDEX idx_users_oauth ON users(provider, oauth_id) WHERE oauth_id IS NOT NULL;
