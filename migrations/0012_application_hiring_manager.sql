-- Hiring manager fields on applications. Captured from the JD or by hand;
-- shown on the Brief tab so you know who the role rolls up to and can open
-- their LinkedIn profile before any conversation. Both nullable — most apps
-- will start without this info and you fill it in as you learn.

ALTER TABLE applications ADD COLUMN hiring_manager_name TEXT;
ALTER TABLE applications ADD COLUMN hiring_manager_linkedin TEXT;
