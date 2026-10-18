-- One-time tokens for email verification and password reset links. Only the
-- SHA-256 hex digest of the raw token is stored.
CREATE TABLE IF NOT EXISTS email_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    purpose VARCHAR(20) NOT NULL CHECK (purpose IN ('verify_email', 'reset_password')),
    token_hash CHAR(64) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_email_tokens_user_purpose ON email_tokens(user_id, purpose);

-- Locale of the reporter, used to pick the language of outcome emails
ALTER TABLE abuse_reports ADD COLUMN IF NOT EXISTS locale VARCHAR(10) NOT NULL DEFAULT '';

-- Claiming a domain now requires a verified email. Accounts created before
-- verification emails existed never had a way to verify, so they are
-- grandfathered in rather than locked out.
UPDATE users SET email_verified = TRUE WHERE email_verified = FALSE;
