-- JD body text + recruiter/point-of-contact on applications.
--
-- jd_text: the full job-description body. We only stored jd_url before, but
-- postings get taken down mid-process (the URL rots) so the description is
-- lost. Persist the pasted/parsed body so it survives.
--
-- recruiter_*: the person running the process (schedules interviews, relays
-- answers). Distinct from hiring_manager_* (the role's manager).

ALTER TABLE applications ADD COLUMN jd_text            TEXT;
ALTER TABLE applications ADD COLUMN recruiter_name     TEXT;
ALTER TABLE applications ADD COLUMN recruiter_email    TEXT;
ALTER TABLE applications ADD COLUMN recruiter_linkedin TEXT;
