BEGIN;
ALTER TABLE guild_characters
    RENAME COLUMN is_sub_leader TO avoid_leadership;

ALTER TABLE guilds
    ALTER COLUMN main_motto SET DEFAULT 0;

ALTER TABLE guilds
    ADD COLUMN icon      bytea,
    ADD COLUMN sub_motto int DEFAULT 0,
    ALTER COLUMN main_motto TYPE int USING 0;
END;