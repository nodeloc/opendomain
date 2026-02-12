-- Add use_default_nameservers field to root_domains table
ALTER TABLE root_domains
ADD COLUMN use_default_nameservers BOOLEAN DEFAULT TRUE;
