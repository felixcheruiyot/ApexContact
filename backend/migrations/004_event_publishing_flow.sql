-- 004_event_publishing_flow.sql
-- Adds draft → pending_review → scheduled/declined workflow for events.

ALTER TABLE events DROP CONSTRAINT IF EXISTS events_status_check;
ALTER TABLE events ADD CONSTRAINT events_status_check
  CHECK (status IN ('draft','pending_review','scheduled','live','completed','cancelled','declined'));

ALTER TABLE events ADD COLUMN IF NOT EXISTS review_note TEXT NOT NULL DEFAULT '';
