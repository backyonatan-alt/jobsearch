-- Application status: a free-form text column rather than an enum.
-- Enums are painful to evolve in Postgres; we validate at the API layer instead.
-- Conventional values: wishlist, applied, screen, interview, offer, rejected, withdrawn.

CREATE TABLE applications (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company     TEXT NOT NULL,
    role        TEXT NOT NULL,
    status      TEXT NOT NULL DEFAULT 'applied',
    source      TEXT,
    jd_url      TEXT,
    location    TEXT,
    salary_note TEXT,
    cv_variant  TEXT,
    notes       TEXT,
    applied_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX applications_user_id_idx ON applications (user_id);
CREATE INDEX applications_user_status_idx ON applications (user_id, status);
CREATE INDEX applications_user_applied_at_idx ON applications (user_id, applied_at DESC);
