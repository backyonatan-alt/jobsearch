-- Public "request access" list. Anyone hitting / can drop their email +
-- a one-line "why I want in". Admin reviews on /admin/people and
-- one-click promotes to invited_emails (which lets that Gmail sign in).
--
-- invited_at is set when the admin promotes the row; it's the join between
-- this table and invited_emails so we never lose the original context note.

CREATE TABLE beta_interest (
    email      TEXT PRIMARY KEY,
    note       TEXT,
    source     TEXT,                              -- optional referrer / channel tag
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    invited_at TIMESTAMPTZ                        -- null = pending, set = promoted
);

CREATE INDEX beta_interest_pending_idx
    ON beta_interest (created_at DESC)
    WHERE invited_at IS NULL;
