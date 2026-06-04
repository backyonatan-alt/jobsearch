-- Per-user cap on AI interview-prep generations (cost control for the beta).
-- New users get 10; the admin can grant more from /admin/people.
ALTER TABLE users ADD COLUMN prep_credits_used  INT NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN prep_credits_limit INT NOT NULL DEFAULT 10;

-- Existing users (the admin + current testers) predate the cap — leave them
-- effectively uncapped. New rows after this migration default to 10.
UPDATE users SET prep_credits_limit = 100000;
