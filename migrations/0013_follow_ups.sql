-- Follow-ups are self-logged outreach: the user records that they reached out
-- to a company themselves. Pursuit never sends anything — this is just a log
-- entry. The MAX(occurred_at) per application is what "resets the clock" in the
-- application list (exposed as last_follow_up_at).

CREATE TABLE follow_ups (
    id              BIGSERIAL PRIMARY KEY,
    application_id  BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    note            TEXT NOT NULL DEFAULT '',
    channel         TEXT NOT NULL DEFAULT '',
    occurred_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX follow_ups_application_id_idx ON follow_ups (application_id);
