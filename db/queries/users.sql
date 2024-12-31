-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByProviderId :one
SELECT *
FROM users
WHERE oauth_provider = $1 AND oauth_provider_id = $2
LIMIT 1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  oauth_provider,
  oauth_provider_id
) VALUES (
  $1,
  $2
)
RETURNING *;

-- name: UpdateUser :exec
UPDATE users
SET
  oauth_provider = $2,
  oauth_provider_id = $3,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
