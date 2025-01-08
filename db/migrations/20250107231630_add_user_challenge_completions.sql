-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_challenge_completions (
  id              BIGSERIAL PRIMARY KEY,
  user_id         BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
  challenge_id    BIGINT NOT NULL REFERENCES challenges (id) ON DELETE CASCADE,
  day_id          BIGINT NOT NULL REFERENCES days (id) ON DELETE CASCADE,
  completed_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (user_id, challenge_id, day_id) -- Ensure users can't complete the same challenge multiple times on the same day
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_challenge_completions;
-- +goose StatementEnd
