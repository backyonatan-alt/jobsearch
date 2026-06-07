-- First-party product-analytics event log. Source of truth for in-app
-- behaviour (activation, AI-moment adoption, funnel) — GA4 stays for the
-- public homepage only. Small-N beta: we read individual user journeys here
-- (SELECT ... WHERE user_id = ? ORDER BY created_at), not aggregate charts.
CREATE TABLE events (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT REFERENCES users(id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    props      JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Per-user timeline reads (the primary beta analysis surface).
CREATE INDEX events_user_created_idx ON events (user_id, created_at DESC);
-- Group-by-name rollups (feature adoption).
CREATE INDEX events_name_created_idx ON events (name, created_at DESC);
