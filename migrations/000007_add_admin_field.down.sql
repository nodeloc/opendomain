-- Remove default admin user
DELETE FROM users WHERE email = 'admin@opendomain.local';

-- Remove is_admin field
ALTER TABLE users DROP COLUMN IF EXISTS is_admin;
