BEGIN;

ALTER TABLE guilds
    ADD COLUMN comment varchar(255) NOT NULL DEFAULT '';

END;