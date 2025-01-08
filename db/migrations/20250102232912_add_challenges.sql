-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS challenge_months (
  id         BIGSERIAL PRIMARY KEY,
  month      INTEGER NOT NULL CHECK (month BETWEEN 1 AND 12),
  year       INTEGER NOT NULL CHECK (year >= 2000),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS days (
  id                  BIGSERIAL PRIMARY KEY,
  challenge_month_id  BIGINT NOT NULL REFERENCES challenge_months (id) ON DELETE CASCADE,
  day_number          INTEGER NOT NULL CHECK (day_number BETWEEN 1 AND 25),
  created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (challenge_month_id, day_number) -- Ensures no duplicate days in a single month
);

-- For performance
CREATE INDEX IF NOT EXISTS idx_days_challenge_month_id
  ON days (challenge_month_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_days_challenge_month_id;
DROP TABLE IF EXISTS days;
DROP TABLE IF EXISTS challenge_months;
-- +goose StatementEnd