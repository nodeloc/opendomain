-- Invitations table
CREATE TABLE IF NOT EXISTS invitations (
  id BIGSERIAL PRIMARY KEY,
  inviter_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  invitee_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  invite_code VARCHAR(20) NOT NULL,
  reward_given BOOLEAN DEFAULT FALSE,
  reward_type VARCHAR(50),
  reward_value TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(invitee_id)
);

-- Add invited_by field to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS invited_by BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS total_invites INT DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS successful_invites INT DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_invitations_inviter ON invitations(inviter_id);
CREATE INDEX IF NOT EXISTS idx_invitations_invitee ON invitations(invitee_id);
CREATE INDEX IF NOT EXISTS idx_users_invited_by ON users(invited_by);
