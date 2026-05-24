-- One current dossier per application. Regenerating UPSERTs in place; we don't
-- keep history (yet — easy to relax to a per-version table later if useful).

CREATE TABLE dossiers (
    id               BIGSERIAL PRIMARY KEY,
    application_id   BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    interviewer_name TEXT,
    content          JSONB NOT NULL,
    model_used       TEXT NOT NULL,
    generated_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (application_id)
);

CREATE INDEX dossiers_application_id_idx ON dossiers (application_id);
