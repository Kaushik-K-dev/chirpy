-- name: CreateChirp :one
INSERT INTO chirps (created_at, updated_at, body, user_id)
VALUES (
    NOW(),
    NOW(),
    $1,
    $2
    )
RETURNING id, created_at, updated_at, body, user_id;

-- name: GetChirpByID :one
SELECT id, created_at, updated_at, body, user_id
FROM chirps
WHERE id = $1;

-- name: ListChirps :many
SELECT id, created_at, updated_at, body, user_id
FROM chirps
ORDER BY created_at ASC;

-- name: UpdateChirp :one
UPDATE chirps
SET updated_at = NOW(), body = $2
WHERE id = $1
RETURNING id, created_at, updated_at, body, user_id;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;