-- 移除索引
DROP INDEX IF EXISTS idx_coupon_usage_coupon_user;

-- 恢复唯一约束
ALTER TABLE coupon_usage ADD CONSTRAINT coupon_usage_coupon_id_user_id_key UNIQUE (coupon_id, user_id);

-- 移除 is_reusable 字段
ALTER TABLE coupons DROP COLUMN IF EXISTS is_reusable;
