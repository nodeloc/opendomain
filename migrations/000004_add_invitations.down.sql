DROP TABLE IF EXISTS invitations;

ALTER TABLE users DROP COLUMN IF EXISTS invited_by;
ALTER TABLE users DROP COLUMN IF EXISTS total_invites;
ALTER TABLE users DROP COLUMN IF EXISTS successful_invites;
