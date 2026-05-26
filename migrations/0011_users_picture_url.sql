-- Store the Google profile picture URL captured from the OAuth ID token so
-- the avatar in the sidebar can show the real face instead of initials.
-- Nullable: existing rows pre-date the column and we only learn the URL on
-- a fresh sign-in.

ALTER TABLE users ADD COLUMN picture_url TEXT;
