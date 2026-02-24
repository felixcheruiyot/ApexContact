-- Migration 009: Rename commentary event_type → audio_video, add audio type
-- ─────────────────────────────────────────────────────────────────────────────

-- 1. Drop the existing CHECK constraint (auto-named by Postgres)
ALTER TABLE events DROP CONSTRAINT IF EXISTS events_event_type_check;

-- 2. Add new CHECK constraint with all valid types
ALTER TABLE events
  ADD CONSTRAINT events_event_type_check
  CHECK (event_type IN ('video', 'commentary', 'audio_video', 'audio'));
-- Note: 'commentary' kept temporarily so any in-flight rows are still valid
-- during the UPDATE below, then the constraint below removes it.

-- 3. Migrate existing 'commentary' rows to 'audio_video'
UPDATE events SET event_type = 'audio_video' WHERE event_type = 'commentary';

-- 4. Now tighten the constraint to remove the legacy 'commentary' value
ALTER TABLE events DROP CONSTRAINT events_event_type_check;
ALTER TABLE events
  ADD CONSTRAINT events_event_type_check
  CHECK (event_type IN ('video', 'audio_video', 'audio'));
