CREATE TABLE payout_accounts (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  account_type   TEXT NOT NULL CHECK (account_type IN ('mpesa','bank')),
  account_number TEXT NOT NULL,
  account_name   TEXT NOT NULL DEFAULT '',
  bank_name      TEXT NOT NULL DEFAULT '',
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id)
);

CREATE TABLE withdrawals (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id           UUID NOT NULL REFERENCES users(id),
  payout_account_id UUID NOT NULL REFERENCES payout_accounts(id),
  amount            NUMERIC(12,2) NOT NULL,
  currency          TEXT NOT NULL DEFAULT 'KES',
  status            TEXT NOT NULL DEFAULT 'pending_otp'
                    CHECK (status IN ('pending_otp','processing','completed','failed','cancelled')),
  intasend_ref      TEXT NOT NULL DEFAULT '',
  failure_reason    TEXT NOT NULL DEFAULT '',
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX ON withdrawals(user_id);
