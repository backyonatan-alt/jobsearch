-- Post-interview debrief (Phase 3a). One per interview round. Feeds the next
-- round's interviewer brief (feed-forward) and captures the "was the prep right?"
-- trust signal.
CREATE TABLE debriefs (
    id             BIGSERIAL PRIMARY KEY,
    interview_id   BIGINT NOT NULL REFERENCES interviews(id) ON DELETE CASCADE,
    application_id BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feel           TEXT NOT NULL DEFAULT '',   -- strong | mixed | rough
    prep_accuracy  TEXT NOT NULL DEFAULT '',   -- spot_on | partly | off
    topics         TEXT NOT NULL DEFAULT '',   -- what actually came up
    notes          TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (interview_id)
);
CREATE INDEX debriefs_application_id_idx ON debriefs (application_id);
