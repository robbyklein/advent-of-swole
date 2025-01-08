-- name: CompleteChallenge :exec
INSERT INTO user_challenge_completions (
  user_id,
  challenge_id,
  day_id
) VALUES (
  $1,
  $2,
  $3
)
ON CONFLICT DO NOTHING;

-- name: GetCompletedChallengesForUser :many
SELECT challenge_id
FROM user_challenge_completions
WHERE user_id = $1 AND day_id = $2;