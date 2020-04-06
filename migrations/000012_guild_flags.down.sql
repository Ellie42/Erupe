BEGIN;
ALTER TABLE guild_characters
    RENAME COLUMN avoid_leadership TO is_sub_leader;
END;