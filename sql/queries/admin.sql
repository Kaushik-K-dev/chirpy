-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;