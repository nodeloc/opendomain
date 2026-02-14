-- Drop the updated constraint
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_years_check;

-- Restore the original constraint
ALTER TABLE orders ADD CONSTRAINT orders_years_check CHECK (years >= 1 AND years <= 10);
