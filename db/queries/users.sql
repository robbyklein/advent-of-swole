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
  oauth_provider = $2,
  oauth_provider_id = $3,
  email = $4,
  timezone = $5,
  display_name = $6,
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
