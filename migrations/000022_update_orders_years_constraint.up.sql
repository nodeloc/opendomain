-- Drop the old years check constraint
ALTER TABLE orders DROP CONSTRAINT IF EXISTS orders_years_check;

-- Add new constraint allowing 0-100 years (0 for lifetime, 1-10 for regular, 100 for lifetime storage)
ALTER TABLE orders ADD CONSTRAINT orders_years_check CHECK (years >= 0 AND years <= 100);
