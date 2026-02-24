-- Migration 011: Add is_public flag to events for the Discover page
ALTER TABLE events
  ADD COLUMN IF NOT EXISTS is_public BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_events_public ON events(is_public) WHERE is_public = true;
