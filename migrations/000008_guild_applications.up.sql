BEGIN;

ALTER TABLE guild_characters ADD COLUMN is_applicant bool DEFAULT false;

END;