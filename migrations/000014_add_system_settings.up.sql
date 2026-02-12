-- 创建 system_settings 表
CREATE TABLE IF NOT EXISTS system_settings (
  id BIGSERIAL PRIMARY KEY,
  setting_key VARCHAR(100) NOT NULL UNIQUE,
  setting_value TEXT NOT NULL,
  description VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认 quota 配置
INSERT INTO system_settings (setting_key, setting_value, description) VALUES
  ('quota_normal', '2', 'Domain quota for normal level'),
  ('quota_basic', '3', 'Domain quota for basic level'),
  ('quota_member', '5', 'Domain quota for member level'),
  ('quota_regular', '10', 'Domain quota for regular level'),
  ('quota_leader', '20', 'Domain quota for leader level');

-- 删除旧的 user_level CHECK 约束
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_user_level_check;

-- 将旧等级映射到新等级（在添加新约束之前）
UPDATE users SET user_level = 'basic' WHERE user_level = 'verified';
UPDATE users SET user_level = 'regular' WHERE user_level = 'vip';
UPDATE users SET user_level = 'leader' WHERE user_level = 'whitelist';

-- 添加新的 user_level CHECK 约束
ALTER TABLE users ADD CONSTRAINT users_user_level_check
  CHECK (user_level IN ('normal', 'basic', 'member', 'regular', 'leader'));
