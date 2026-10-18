-- users.email_verified is intentionally left as is: the up migration's
-- grandfathering cannot be told apart from real verifications.
ALTER TABLE abuse_reports DROP COLUMN IF EXISTS locale;
DROP TABLE IF EXISTS email_tokens;
