-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS settings (
  id              BIGSERIAL PRIMARY KEY,
  user_id         BIGINT NOT NULL UNIQUE REFERENCES users (id) ON DELETE CASCADE,
  timezone        TEXT NOT NULL,
  display_name    TEXT NOT NULL,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for quick lookup by user_id
CREATE INDEX IF NOT EXISTS idx_settings_user_id
  ON settings (user_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_settings_user_id;
DROP TABLE IF EXISTS settings;
-- +goose StatementEnd
