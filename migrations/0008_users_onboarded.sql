-- Drives the first-run experience. New users sign in with onboarded_at NULL
-- and see the overlay. Existing users (already have applications) get
-- onboarded_at set so the overlay doesn't bug them.

ALTER TABLE users ADD COLUMN onboarded_at TIMESTAMPTZ;

UPDATE users
SET onboarded_at = now()
WHERE id IN (SELECT DISTINCT user_id FROM applications);
