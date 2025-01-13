-- name: GetChallenge :one
SELECT *
FROM challenges
WHERE id = $1
LIMIT 1;

-- name: GetChallengeByDescription :one
SELECT *
FROM challenges
WHERE description = $1
LIMIT 1;

-- name: ListChallenges :many
SELECT *
FROM challenges
ORDER BY id;

-- name: CreateChallenge :one
INSERT INTO challenges (
  description,
  description_metric,
  category,
  muscle_groups,
  difficulty,
  calories_burned_estimate
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: UpdateChallenge :exec
UPDATE challenges
SET
  description = $2,
  description_metric = $3,
  category = $4,
  muscle_groups = $5,
  difficulty = $6,
  calories_burned_estimate = $7,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteChallenge :exec
DELETE FROM challenges
WHERE id = $1;
