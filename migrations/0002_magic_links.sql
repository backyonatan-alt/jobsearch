CREATE TABLE magic_links (
    token       TEXT PRIMARY KEY,
    email       TEXT NOT NULL,
    expires_at  TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX magic_links_email_idx ON magic_links (lower(email));
CREATE INDEX magic_links_expires_idx ON magic_links (expires_at) WHERE consumed_at IS NULL;
