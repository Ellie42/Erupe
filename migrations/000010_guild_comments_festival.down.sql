BEGIN;

ALTER TABLE guilds
    DROP COLUMN comment,
    DROP COLUMN festival_colour;

DROP TYPE festival_colour;

END;