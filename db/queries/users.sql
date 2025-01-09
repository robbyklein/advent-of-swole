-- name: GetUser :one
SELECT id, oauth_provider, oauth_provider_id, email, timezone, display_name, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByProviderId :one
SELECT id, oauth_provider, oauth_provider_id, email, timezone, display_name, created_at, updated_at
FROM users
WHERE oauth_provider = $1 AND oauth_provider_id = $2
LIMIT 1;

-- name: ListUsers :many
SELECT id, oauth_provider, oauth_provider_id, email, timezone, display_name, created_at, updated_at
FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  oauth_provider,
  oauth_provider_id,
  email,
  timezone,
  display_name
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING id, oauth_provider, oauth_provider_id, email, timezone, display_name, created_at, updated_at;

-- name: UpdateUser :exec
UPDATE users
SET
  timezone = $2,
  display_name = $3,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
