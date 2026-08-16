-- Pre-round reminder emails. sent_emails is the idempotency ledger: one row
-- per (interview, kind) claimed before sending so a restart never double-sends.
ALTER TABLE users ADD COLUMN email_reminders BOOLEAN NOT NULL DEFAULT true;

-- IANA timezone captured from the browser at save time, so the reminder can
-- render wall-clock time the way the user sees it.
ALTER TABLE interviews ADD COLUMN tz TEXT;

CREATE TABLE sent_emails (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    interview_id BIGINT NOT NULL REFERENCES interviews(id) ON DELETE CASCADE,
    kind         TEXT NOT NULL,
    sent_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (interview_id, kind)
);
