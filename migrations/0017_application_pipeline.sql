-- Per-application interview pipeline: an ordered, user-defined list of stages
-- (e.g. Recruiter call → Manager → Take-home → HRBP → Offer) that each app
-- tracks independently of the global status. Stored as a JSONB array of
-- {name, done}. Mutable any time — stages are added/renamed/reordered/checked
-- off as the process unfolds.

ALTER TABLE applications ADD COLUMN pipeline JSONB NOT NULL DEFAULT '[]'::jsonb;
