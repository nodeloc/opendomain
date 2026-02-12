-- Restore CHECK constraint with free_domain
ALTER TABLE coupons DROP CONSTRAINT IF EXISTS coupons_discount_type_check;
ALTER TABLE coupons ADD CONSTRAINT coupons_discount_type_check CHECK (discount_type IN ('percentage', 'fixed', 'free_domain', 'quota_increase'));
