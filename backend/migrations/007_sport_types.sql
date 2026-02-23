-- Migration 007: Expand sport_type to support knowledge & entertainment categories.
-- Drops the old boxing/racing-only constraint, remaps legacy values, and adds the new constraint.

ALTER TABLE events DROP CONSTRAINT IF EXISTS events_sport_type_check;

-- Remap legacy sport types to 'other'
UPDATE events SET sport_type = 'other' WHERE sport_type IN ('boxing', 'racing');

ALTER TABLE events ADD CONSTRAINT events_sport_type_check
  CHECK (sport_type IN (
    'sales','mentoring','education','business','legal','fitness','visa',
    'music','gaming','cooking','community','other'
  ));
