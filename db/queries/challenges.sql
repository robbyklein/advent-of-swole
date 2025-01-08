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
  category,
  muscle_groups,
  difficulty,
  calories_burned_estimate
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: UpdateChallenge :exec
UPDATE challenges
SET
  description = $2,
  category = $3,
  muscle_groups = $4,
  difficulty = $5,
  calories_burned_estimate = $6,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteChallenge :exec
DELETE FROM challenges
WHERE id = $1;
