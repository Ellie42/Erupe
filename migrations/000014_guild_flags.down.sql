BEGIN;
ALTER TABLE guild_characters
    RENAME COLUMN avoid_leadership TO is_sub_leader;

ALTER TABLE guilds
    DROP COLUMN icon,
    ALTER COLUMN main_motto TYPE varchar USING '',
    DROP COLUMN sub_motto;

ALTER TABLE guilds
    ALTER COLUMN main_motto SET DEFAULT '';

END;