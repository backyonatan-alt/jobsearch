-- Admin flag + invited-emails table for moving the closed-beta allow-list
-- out of the .env file and into a real, mutable surface.

ALTER TABLE users ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE invited_emails (
    email              TEXT PRIMARY KEY,
    invited_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    invited_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    note               TEXT
);

CREATE INDEX invited_emails_invited_by_idx ON invited_emails (invited_by_user_id);
