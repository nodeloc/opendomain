-- Add first_failed_at column to domains table for tracking health check failures
ALTER TABLE domains ADD COLUMN IF NOT EXISTS first_failed_at TIMESTAMP NULL DEFAULT NULL;

-- Create index for efficient querying
CREATE INDEX IF NOT EXISTS idx_domains_first_failed_at ON domains(first_failed_at);
