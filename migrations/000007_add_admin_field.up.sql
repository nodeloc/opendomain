-- Add is_admin field to users table
ALTER TABLE users ADD COLUMN is_admin BOOLEAN DEFAULT FALSE;

-- Create index
CREATE INDEX idx_users_is_admin ON users(is_admin);

-- Create default admin user (password: admin123)
-- Password hash for 'admin123' using bcrypt
INSERT INTO users (
  username,
  email,
  password_hash,
  is_admin,
  email_verified,
  domain_quota,
  invite_code,
  user_level
) VALUES (
  'admin',
  'admin@opendomain.local',
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', -- admin123
  TRUE,
  TRUE,
  999,
  'ADMIN_' || substr(md5(random()::text), 1, 10),
  'whitelist'
);
