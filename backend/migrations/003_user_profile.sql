-- Add extended profile fields to users
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS age     INT,
  ADD COLUMN IF NOT EXISTS gender  TEXT NOT NULL DEFAULT '',
  ADD COLUMN IF NOT EXISTS country TEXT NOT NULL DEFAULT '';
