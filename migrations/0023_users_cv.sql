-- Candidate CV, one per user, plain text (pasted or Haiku-extracted from an
-- upload). Feeds the per-round "make sure they hear" playbook section.
ALTER TABLE users ADD COLUMN cv_text TEXT;
