-- name: GetChallengeMonth :one
SELECT *
FROM challenge_months
WHERE id = $1
LIMIT 1;

-- name: ListChallengeMonths :many
SELECT *
FROM challenge_months
ORDER BY year DESC, month DESC;

-- name: CreateChallengeMonth :one
INSERT INTO challenge_months (
  month,
  year
) VALUES (
  $1,
  $2
)
RETURNING *;

-- name: UpdateChallengeMonth :exec
UPDATE challenge_months
SET
  month = $2,
  year = $3,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteChallengeMonth :exec
DELETE FROM challenge_months
WHERE id = $1;

-- name: GetMostRecentChallengeMonth :one
SELECT *
FROM challenge_months
ORDER BY year DESC, month DESC
LIMIT 1;

-- name: GetChallengeMonthByYearMonth :one
SELECT *
FROM challenge_months
WHERE year = $1 AND month = $2
LIMIT 1;

-- name: GetDayByMonthIDNumber :one
SELECT *
FROM days
WHERE challenge_month_id = $1
  AND day_number = $2
LIMIT 1;
