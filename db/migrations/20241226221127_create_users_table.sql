-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id                BIGSERIAL PRIMARY KEY,
  oauth_provider    TEXT NOT NULL,
  oauth_provider_id TEXT NOT NULL UNIQUE,
  email             TEXT NOT NULL,
  timezone          TEXT NOT NULL,
  display_name      TEXT NOT NULL,
  measurement_system      TEXT NOT NULL DEFAULT 'imperial',
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for quick lookup by OAuth provider and ID
CREATE INDEX IF NOT EXISTS idx_users_oauth_provider_id
  ON users (oauth_provider, oauth_provider_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_oauth_provider_id;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd