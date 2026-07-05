-- Decouple a "round" from a calendar event: one-tap rounds have no date.
-- scheduled = true means it came from a real calendar interview (has a real
-- starts_at); false means a quick-logged round (starts_at defaulted to now()).
-- Existing rows were all calendar events → default true.
ALTER TABLE interviews ADD COLUMN scheduled BOOLEAN NOT NULL DEFAULT true;
