ALTER TABLE domain_scan_summaries
  DROP COLUMN IF EXISTS safe_browsing_status,
  DROP COLUMN IF EXISTS virustotal_status;
