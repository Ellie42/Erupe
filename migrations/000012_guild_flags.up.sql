BEGIN;
ALTER TABLE guild_characters
    RENAME COLUMN is_sub_leader TO avoid_leadership;
END;