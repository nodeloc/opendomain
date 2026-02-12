-- Drop payment configuration table
DROP TABLE IF EXISTS payment_configs;

-- Drop payments table
DROP TABLE IF EXISTS payments;

-- Drop orders table
DROP TABLE IF EXISTS orders;

-- Remove coupon_usage domain_id index
DROP INDEX IF EXISTS idx_coupon_usage_domain_id;

-- Remove pricing fields from root_domains
ALTER TABLE root_domains DROP COLUMN IF EXISTS is_free;
ALTER TABLE root_domains DROP COLUMN IF EXISTS lifetime_price;
ALTER TABLE root_domains DROP COLUMN IF EXISTS price_per_year;

DROP INDEX IF EXISTS idx_root_domains_is_free;
