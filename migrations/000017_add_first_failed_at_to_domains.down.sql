-- Remove index
DROP INDEX IF EXISTS idx_domains_first_failed_at;

-- Remove first_failed_at column from domains table
ALTER TABLE domains DROP COLUMN IF EXISTS first_failed_at;
