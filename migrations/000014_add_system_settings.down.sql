-- 恢复旧等级
UPDATE users SET user_level = 'verified' WHERE user_level = 'basic';
UPDATE users SET user_level = 'vip' WHERE user_level = 'regular';
UPDATE users SET user_level = 'whitelist' WHERE user_level = 'leader';

-- 恢复旧 CHECK 约束
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_user_level_check;
ALTER TABLE users ADD CONSTRAINT users_user_level_check
  CHECK (user_level IN ('normal', 'verified', 'vip', 'whitelist'));

-- 删除 system_settings 表
DROP TABLE IF EXISTS system_settings;
