-- name: GetSettings :one
SELECT *
FROM settings
WHERE user_id = $1
LIMIT 1;

-- name: CreateSettings :one
INSERT INTO settings (
  user_id,
  timezone,
  display_name
) VALUES (
  $1,
  $2,
  $3
)
RETURNING *;

-- name: UpdateSettings :exec
UPDATE settings
SET
  timezone = $2,
  display_name = $3,
  updated_at = NOW()
WHERE user_id = $1;

-- name: DeleteSettings :exec
DELETE FROM settings
WHERE user_id = $1;
