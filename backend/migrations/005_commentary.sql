-- ── Commentary event type support ────────────────────────────────────────────

-- Extend events table with commentary-specific columns
ALTER TABLE events ADD COLUMN IF NOT EXISTS event_type TEXT NOT NULL DEFAULT 'video'
  CHECK (event_type IN ('video', 'commentary'));
ALTER TABLE events ADD COLUMN IF NOT EXISTS livekit_room TEXT NOT NULL DEFAULT '';
ALTER TABLE events ADD COLUMN IF NOT EXISTS teaser_hook TEXT NOT NULL DEFAULT '';

-- ── Lobby participants ────────────────────────────────────────────────────────

CREATE TABLE IF NOT EXISTS lobby_participants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  nickname TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'listener' CHECK (role IN ('host', 'speaker', 'listener')),
  joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (event_id, user_id)
);

-- ── Chat messages (persisted for replay) ─────────────────────────────────────

CREATE TABLE IF NOT EXISTS lobby_messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  nickname TEXT NOT NULL,
  content TEXT NOT NULL,
  message_type TEXT NOT NULL DEFAULT 'text' CHECK (message_type IN ('text', 'meme')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_lobby_messages_event ON lobby_messages(event_id, created_at);
CREATE INDEX IF NOT EXISTS idx_lobby_participants_event ON lobby_participants(event_id);
