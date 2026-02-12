-- 添加优惠券可重复使用标志
ALTER TABLE coupons ADD COLUMN is_reusable BOOLEAN DEFAULT FALSE;

-- 移除 coupon_usage 表的唯一约束，允许同一用户多次使用可重复使用的优惠券
ALTER TABLE coupon_usage DROP CONSTRAINT IF EXISTS coupon_usage_coupon_id_user_id_key;

-- 添加普通索引（用于查询性能）
CREATE INDEX IF NOT EXISTS idx_coupon_usage_coupon_user ON coupon_usage(coupon_id, user_id);
