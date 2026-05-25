-- Interview events linked to applications. Source is 'ics' for calendar
-- imports, 'manual' for things the user types in. We keep (application_id, uid)
-- unique so re-importing an updated invite updates in place instead of
-- duplicating — the calendar UID is stable across REQUEST/UPDATE sequences.

CREATE TABLE interviews (
    id              BIGSERIAL PRIMARY KEY,
    application_id  BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source          TEXT NOT NULL DEFAULT 'manual',
    uid             TEXT,
    summary         TEXT NOT NULL DEFAULT '',
    location        TEXT,
    description     TEXT,
    starts_at       TIMESTAMPTZ NOT NULL,
    ends_at         TIMESTAMPTZ,
    all_day         BOOLEAN NOT NULL DEFAULT false,
    organizer       JSONB,
    attendees       JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX interviews_application_id_idx ON interviews (application_id, starts_at);
CREATE INDEX interviews_user_starts_idx ON interviews (user_id, starts_at);
CREATE UNIQUE INDEX interviews_app_uid_uidx ON interviews (application_id, uid) WHERE uid IS NOT NULL;
