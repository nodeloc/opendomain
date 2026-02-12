-- Add per-domain nameservers (previously shared at root_domain level)
ALTER TABLE domains ADD COLUMN nameservers TEXT;
ALTER TABLE domains ADD COLUMN use_default_nameservers BOOLEAN DEFAULT TRUE;

-- Migrate existing data: copy nameservers from root_domains to each domain
UPDATE domains d SET
  nameservers = r.nameservers,
  use_default_nameservers = r.use_default_nameservers
FROM root_domains r WHERE d.root_domain_id = r.id;
