-- A/B: prep-first cold start (see AB_TESTS.md).
-- New signups are assigned a variant at INSERT time in auth.UpsertUserAndSession.
-- NULL = pre-test users / control (the guided tour). New users get 'prepfirst'.
ALTER TABLE users ADD COLUMN IF NOT EXISTS onboarding_variant text;
