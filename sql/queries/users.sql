-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email)
VALUES (
    NOW(),
    NOW(),
    $1
)
RETURNING id, created_at, updated_at, email;

-- name: GetUserByID :one
SELECT id, created_at, updated_at, email
FROM users
WHERE id = $1;


-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, created_at, updated_at, email
FROM users
ORDER BY created_at ASC;

-- name: UpdateUserEmail :one
UPDATE users
SET updated_at = NOW(), email = $2
WHERE id = $1
RETURNING id, created_at, updated_at, email;