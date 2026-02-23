-- Migration 007: Expand sport_type to support knowledge & entertainment categories.
-- Drops the old boxing/racing-only constraint and replaces it with the full category list.

ALTER TABLE events DROP CONSTRAINT IF EXISTS events_sport_type_check;

ALTER TABLE events ADD CONSTRAINT events_sport_type_check
  CHECK (sport_type IN (
    'sales','mentoring','education','business','legal','fitness','visa',
    'music','gaming','cooking','community','other'
  ));
