-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS challenges (
  id                     BIGSERIAL PRIMARY KEY,
  description            TEXT NOT NULL UNIQUE,
  category               TEXT NOT NULL,
  muscle_groups          TEXT[] NOT NULL,
  difficulty             INTEGER NOT NULL CHECK (difficulty BETWEEN 1 AND 5),
  calories_burned_estimate INTEGER NOT NULL,
  created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS day_challenges (
  day_id       BIGINT NOT NULL REFERENCES days (id) ON DELETE CASCADE,
  challenge_id BIGINT NOT NULL REFERENCES challenges (id) ON DELETE CASCADE,
  PRIMARY KEY (day_id, challenge_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS day_challenges;
DROP TABLE IF EXISTS challenges;
-- +goose StatementEnd