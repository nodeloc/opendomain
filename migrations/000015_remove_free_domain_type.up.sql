-- Convert existing free_domain coupons to quota_increase
UPDATE coupons SET discount_type = 'quota_increase', quota_increase = 1 WHERE discount_type = 'free_domain';

-- Update CHECK constraint to remove free_domain
ALTER TABLE coupons DROP CONSTRAINT IF EXISTS coupons_discount_type_check;
ALTER TABLE coupons ADD CONSTRAINT coupons_discount_type_check CHECK (discount_type IN ('percentage', 'fixed', 'quota_increase'));
