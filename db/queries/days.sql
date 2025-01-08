-- name: GetDay :one
SELECT *
FROM days
WHERE id = $1
LIMIT 1;

-- name: ListDaysForMonth :many
SELECT *
FROM days
WHERE challenge_month_id = $1
ORDER BY day_number;

-- name: CreateDay :one
INSERT INTO days (
  challenge_month_id,
  day_number
) VALUES (
  $1,
  $2
)
RETURNING *;

-- name: UpdateDay :exec
UPDATE days
SET
  day_number = $2,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteDay :exec
DELETE FROM days
WHERE id = $1;
