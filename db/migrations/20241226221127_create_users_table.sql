-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id                       BIGSERIAL PRIMARY KEY,
  email                    CITEXT NOT NULL UNIQUE,
  password_hash            TEXT NOT NULL,
  created_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  confirmed_at             TIMESTAMPTZ,
  confirmation_token       TEXT,
  confirmation_sent_at     TIMESTAMPTZ,
  reset_password_token     TEXT,
  reset_password_sent_at   TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_confirmation_token
  ON users (confirmation_token);

CREATE INDEX IF NOT EXISTS idx_users_reset_password_token
  ON users (reset_password_token);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_reset_password_token;
DROP INDEX IF EXISTS idx_users_confirmation_token;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd