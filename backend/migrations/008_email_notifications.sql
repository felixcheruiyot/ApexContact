-- Email verification columns on users
ALTER TABLE users
  ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT FALSE,
  ADD COLUMN verification_token TEXT,
  ADD COLUMN verification_token_expires_at TIMESTAMPTZ;

-- Notification audit log
CREATE TABLE notification_log (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  notification_type TEXT NOT NULL,  -- 'email_verification', 'event_reminder', etc.
  reference_id TEXT,                -- event_id or other entity id
  channel TEXT NOT NULL DEFAULT 'email',
  sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX ON notification_log(user_id, notification_type, reference_id);
