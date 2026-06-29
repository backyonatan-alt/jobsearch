-- Per-round interview prep. A dossier may now be scoped to a specific interview
-- (round); interview_id NULL is the application-level ("General") prep, which is
-- also what every pre-existing dossier becomes (no data loss). Round-specific
-- prep wins over General when both exist.
ALTER TABLE dossiers ADD COLUMN interview_id BIGINT REFERENCES interviews(id) ON DELETE CASCADE;

-- Was: one dossier per application. Now: at most one per interview, plus at most
-- one application-level (General) slot. Drop the old single-per-app UNIQUE by
-- whatever name Postgres gave it — leaving it would block per-round rows.
DO $$
DECLARE c text;
BEGIN
  FOR c IN
    SELECT con.conname
    FROM pg_constraint con
    JOIN pg_class rel ON rel.oid = con.conrelid
    WHERE rel.relname = 'dossiers' AND con.contype = 'u'
  LOOP
    EXECUTE format('ALTER TABLE dossiers DROP CONSTRAINT %I', c);
  END LOOP;
END $$;
CREATE UNIQUE INDEX dossiers_interview_uidx       ON dossiers (interview_id)   WHERE interview_id IS NOT NULL;
CREATE UNIQUE INDEX dossiers_app_general_uidx     ON dossiers (application_id) WHERE interview_id IS NULL;
CREATE INDEX        dossiers_interview_id_idx     ON dossiers (interview_id);
