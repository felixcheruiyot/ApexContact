-- Live Streamify initial schema
-- Run with: psql $DATABASE_URL -f migrations/001_init.sql

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ─── Users ────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS users (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email          TEXT NOT NULL UNIQUE,
    password_hash  TEXT NOT NULL,
    full_name      TEXT NOT NULL DEFAULT '',
    phone          TEXT NOT NULL DEFAULT '',
    role           TEXT NOT NULL DEFAULT 'viewer' CHECK (role IN ('viewer','promoter','broadcaster','admin')),
    is_locked      BOOLEAN NOT NULL DEFAULT false,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_role  ON users (role);

-- ─── Events ───────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS events (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    promoter_id    UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    title          TEXT NOT NULL,
    description    TEXT NOT NULL DEFAULT '',
    sport_type     TEXT NOT NULL CHECK (sport_type IN ('boxing','racing')),
    scheduled_at   TIMESTAMPTZ NOT NULL,
    status         TEXT NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled','live','completed','cancelled')),
    price          NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency       TEXT NOT NULL DEFAULT 'KES',
    thumbnail_url  TEXT NOT NULL DEFAULT '',
    stream_key     TEXT NOT NULL UNIQUE,
    hls_path       TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_events_status      ON events (status);
CREATE INDEX idx_events_sport_type  ON events (sport_type);
CREATE INDEX idx_events_scheduled   ON events (scheduled_at);
CREATE INDEX idx_events_promoter    ON events (promoter_id);

-- ─── Payments ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS payments (
    id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id        UUID NOT NULL REFERENCES users(id),
    event_id       UUID NOT NULL REFERENCES events(id),
    amount         NUMERIC(12,2) NOT NULL,
    currency       TEXT NOT NULL DEFAULT 'KES',
    status         TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','success','failed','cancelled')),
    intasend_ref   TEXT NOT NULL DEFAULT '',
    phone_number   TEXT NOT NULL DEFAULT '',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payments_user    ON payments (user_id);
CREATE INDEX idx_payments_event   ON payments (event_id);
CREATE INDEX idx_payments_status  ON payments (status);

-- ─── Subscriptions (stream tokens) ───────────────────────────────────────────
CREATE TABLE IF NOT EXISTS subscriptions (
    id                  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id             UUID NOT NULL REFERENCES users(id),
    event_id            UUID NOT NULL REFERENCES events(id),
    payment_id          UUID NOT NULL REFERENCES payments(id),
    stream_token        TEXT NOT NULL UNIQUE,
    device_fingerprint  TEXT NOT NULL DEFAULT '',
    ip_lock             TEXT NOT NULL DEFAULT '',
    active_session_id   TEXT NOT NULL DEFAULT '',
    expires_at          TIMESTAMPTZ NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, event_id)
);

CREATE INDEX idx_subs_token    ON subscriptions (stream_token);
CREATE INDEX idx_subs_user     ON subscriptions (user_id);
CREATE INDEX idx_subs_event    ON subscriptions (event_id);

-- ─── Stream Analytics ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS stream_analytics (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id        UUID NOT NULL REFERENCES events(id),
    viewer_count    INT NOT NULL DEFAULT 0,
    peak_viewers    INT NOT NULL DEFAULT 0,
    total_revenue   NUMERIC(12,2) NOT NULL DEFAULT 0,
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analytics_event ON stream_analytics (event_id);

-- ─── Fraud Flags ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS fraud_flags (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID NOT NULL REFERENCES users(id),
    subscription_id  UUID NOT NULL REFERENCES subscriptions(id),
    reason           TEXT NOT NULL,
    detected_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved         BOOLEAN NOT NULL DEFAULT false,
    resolved_at      TIMESTAMPTZ
);

CREATE INDEX idx_fraud_user      ON fraud_flags (user_id);
CREATE INDEX idx_fraud_resolved  ON fraud_flags (resolved);
