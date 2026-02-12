-- Remove use_default_nameservers field from root_domains table
ALTER TABLE root_domains
DROP COLUMN IF EXISTS use_default_nameservers;
