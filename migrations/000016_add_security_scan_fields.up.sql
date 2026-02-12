ALTER TABLE domain_scan_summaries
  ADD COLUMN IF NOT EXISTS safe_browsing_status VARCHAR(20) DEFAULT 'unknown',
  ADD COLUMN IF NOT EXISTS virustotal_status VARCHAR(20) DEFAULT 'unknown';
